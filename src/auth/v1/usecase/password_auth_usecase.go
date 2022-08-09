package usecase

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	corporateModel "github.com/Bhinneka/user-service/src/corporate/v2/model"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

// RequestTokenPasswordType function for validate request token with grantType is password
func (au *AuthUseCaseImpl) RequestTokenPasswordType(ctxReq context.Context, data *model.RequestToken) (*memberModel.Member, int, error) {
	ctx := "AuthUseCase-RequestTokenPasswordType"
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{"email": data.Email}
	defer tr.Finish(tags)

	// validate request
	mErr := jsonschema.ValidateTemp(textGetTokenPass, data)
	if mErr != nil {
		return nil, http.StatusBadRequest, mErr
	}

	memberResult := <-au.MemberQueryRead.FindByEmail(ctxReq, data.Email)
	if memberResult.Error != nil {
		if memberResult.Error == sql.ErrNoRows {
			memberResult.Error = errors.New(model.ErrorInvalidUsernameOrPasswordBahasa)
		}
		return nil, http.StatusUnauthorized, memberResult.Error
	}

	member, ok := memberResult.Result.(memberModel.Member)
	if !ok {
		return nil, http.StatusUnauthorized, errors.New(msgResultNotMember)
	}

	// validate member login password
	err := au.validateMemberLoginPassword(member)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// validate password
	passwordSaltDB := strings.TrimSpace(member.Salt)
	passwordHasher := memberModel.NewPBKDF2Hasher(memberModel.SaltSize, memberModel.SaltSize, memberModel.IterationsCount, sha1.New)
	passwordHasher.ParseSalt(passwordSaltDB)
	base64Data := base64.StdEncoding.EncodeToString(passwordHasher.Hash([]byte(data.Password)))

	if member.Password != base64Data {
		// check attempt and save it
		attemptResult := <-au.checkAttempt(ctxReq, data.Email)
		if attemptResult.Error != nil && attemptResult.Error.Error() == helper.ErrorLoginAttempt {
			err := errors.New(model.ErrorAccountBlockedBahasa)
			return nil, http.StatusBadRequest, err
		}
		tags["encoded_password"] = member.Password
		tags["email"] = data.Email
		tags["base64"] = base64Data

		err := errors.New(model.ErrorInvalidUsernameOrPasswordBahasa)
		return nil, http.StatusBadRequest, err
	}

	key := fmt.Sprintf(keyAttempt, member.Email)

	<-au.LoginAttemptRepo.Delete(ctxReq, key)
	passwordHasher = nil
	return &member, 200, nil
}

// validateMemberLoginPassword function for validate request token with grantType is password
func (au *AuthUseCaseImpl) validateMemberLoginPassword(member memberModel.Member) error {
	// check member status
	// only active member can login
	if member.StatusString == memberModel.InactiveString {
		err := errors.New(model.ErrorAccountInActiveBahasa)
		return err
	}
	if member.StatusString == memberModel.NewString {
		err := errors.New(model.ErrorNewAccountBahasa)
		return err
	}

	if member.StatusString == memberModel.BlockedString {
		err := errors.New(model.ErrorAccountBlockedBahasa)
		return err
	}

	// check whether member has password or not
	if len(member.Password) == 0 && len(member.Salt) == 0 {
		if len(member.SocialMedia.AzureID) > 0 {
			err := errors.New(model.ErrorEmailUnregisteredBahasa)
			return err
		}

		if len(member.SocialMedia.FacebookID) > 0 {
			err := errors.New(model.ErrorEmailUnregisteredBahasa)
			return err
		}

		if len(member.SocialMedia.GoogleID) > 0 {
			err := errors.New(model.ErrorEmailUnregisteredBahasa)
			return err
		}
		err := errors.New(model.ErrorInvalidUsernameOrPasswordBahasa)
		return err
	}
	return nil
}

// RequestTokenPasswordTypeCorporate function for validate request token with grantType is password for corporate
func (au *AuthUseCaseImpl) RequestTokenPasswordTypeCorporate(ctxReq context.Context, data *model.RequestToken) (sharedModel.B2BContactData, int, error) {
	ctx := "AuthUseCase-RequestTokenPasswordTypeCorporate"

	trace := tracer.StartTrace(ctxReq, ctx)
	defer trace.Finish(nil)

	member := sharedModel.B2BContactData{}
	// validate request
	mErr := jsonschema.ValidateTemp(textGetTokenPass, data)
	if mErr != nil {
		return member, http.StatusBadRequest, mErr
	}

	memberResult := <-au.CorporateContactQueryRead.FindContactCorporateByEmail(ctxReq, data.Email)
	if memberResult.Error != nil {
		if memberResult.Error == sql.ErrNoRows {
			memberResult.Error = errors.New(model.ErrorInvalidUsernameOrPasswordBahasa)
		}
		return member, http.StatusUnauthorized, fmt.Errorf(model.ErrorIncorrectMemberTypeMicrosite, data.MemberType)
	}

	member, ok := memberResult.Result.(sharedModel.B2BContactData)
	if !ok {
		return member, http.StatusUnauthorized, errors.New(msgResultNotMember)
	}

	// validate member login password
	err := au.validateMemberLoginPasswordCorporate(&member)
	if err != nil {
		return member, http.StatusBadRequest, err
	}

	// vallidate password
	member, httpStatus, err := au.ValidatePasswordCorporate(ctxReq, member, data)
	if err != nil {
		return member, httpStatus, err
	}

	key := fmt.Sprintf(keyAttempt, member.Email)

	<-au.LoginAttemptRepo.Delete(ctxReq, key)
	return member, 200, nil
}

func (au *AuthUseCaseImpl) ValidatePasswordCorporate(ctxReq context.Context, member sharedModel.B2BContactData, data *model.RequestToken) (sharedModel.B2BContactData, int, error) {
	if member.Salt != "" {
		// validate password
		passwordSaltDB := strings.TrimSpace(member.Salt)
		passwordHasher := memberModel.NewPBKDF2Hasher(memberModel.SaltSize, memberModel.SaltSize, memberModel.IterationsCount, sha1.New)
		passwordHasher.ParseSalt(passwordSaltDB)
		base64Data := base64.StdEncoding.EncodeToString(passwordHasher.Hash([]byte(data.Password)))

		if member.Password != base64Data {
			// check attempt and save it
			attemptResult := <-au.checkAttempt(ctxReq, data.Email)
			if attemptResult.Error != nil && attemptResult.Error.Error() == helper.ErrorLoginAttempt {
				err := errors.New(model.ErrorAccountBlockedBahasa)
				return member, http.StatusBadRequest, err
			}

			err := errors.New(model.ErrorInvalidUsernameOrPasswordBahasa)
			return member, http.StatusBadRequest, err
		}
	} else {
		passwordSlice := strings.Split(member.Password, ":")

		// check whether member has password or not
		if len(passwordSlice) < 2 || (len(passwordSlice[0]) == 0 && len(passwordSlice[1]) == 0) {
			err := errors.New(model.ErrorInvalidUsernameOrPasswordBahasa)
			return member, http.StatusBadRequest, err
		}

		// validate password
		passwordSaltDB := strings.TrimSpace(passwordSlice[1])
		// Converts users input password + generated salt to bytes
		passwordCombine := []byte(passwordSaltDB + data.Password)
		bytes := []byte(passwordCombine)

		// Converts string to sha2
		h := sha256.New()                   // new sha256 object
		h.Write(bytes)                      // data is now converted to hex
		code := h.Sum(nil)                  // code is now the hex sum
		codestr := hex.EncodeToString(code) // converts hex to string

		if passwordSlice[0] != codestr {
			err := errors.New(model.ErrorInvalidUsernameOrPasswordBahasa)
			return member, http.StatusBadRequest, err
		}
	}

	return member, 200, nil
}

// RequestTokenPasswordTypeMicrosite function for validate request token with grantType is password for microsite
func (au *AuthUseCaseImpl) RequestTokenPasswordTypeMicrosite(ctxReq context.Context, data *model.RequestToken) (sharedModel.B2BContactData, int, error) {
	ctx := "AuthUseCase-RequestTokenPasswordTypeMicrosite"

	trace := tracer.StartTrace(ctxReq, ctx)
	defer trace.Finish(nil)

	member := sharedModel.B2BContactData{}
	// validate request
	mErr := jsonschema.ValidateTemp(textGetTokenPass, data)
	if mErr != nil {
		return member, http.StatusBadRequest, mErr
	}
	// use came method as
	memberResult := <-au.CorporateContactQueryRead.FindContactMicrositeByEmail(ctxReq, data.Email, data.TransactionType, data.MemberType)
	if memberResult.Error != nil {

		if memberResult.Error == sql.ErrNoRows {
			memberResult.Error = fmt.Errorf(model.ErrorIncorrectMemberTypeMicrosite, data.MemberType)
		}
		return member, http.StatusUnauthorized, memberResult.Error
	}

	member, ok := memberResult.Result.(sharedModel.B2BContactData)
	if !ok {
		err := errors.New(msgResultNotMember)
		return member, http.StatusUnauthorized, err
	}

	// validate member login password
	err := au.validateMemberLoginPasswordMicrosite(member)
	if err != nil {
		return member, http.StatusBadRequest, err
	}

	// vallidate password
	member, httpStatus, err := au.ValidatePasswordCorporate(ctxReq, member, data)
	if err != nil {
		return member, httpStatus, err
	}

	key := fmt.Sprintf(keyAttempt, member.Email)

	<-au.LoginAttemptRepo.Delete(ctxReq, key)
	return member, 200, nil
}

// validateMemberLoginPasswordCorporate function for validate request token with grantType is password
func (au *AuthUseCaseImpl) validateMemberLoginPasswordCorporate(member *sharedModel.B2BContactData) error {
	if member.LoginType == model.LoginTypeShopcart {
		if len(member.Status) == 0 || strings.ToUpper(member.Status) != strings.ToUpper(corporateModel.StatusActivated) {
			return errors.New(model.ErrorAccountInActiveBahasa)
		}
	} else {
		// Find by contact id if contact corporate
		accountContactResult := <-au.CorporateAccContactQueryRead.FindByAccountContactID(member.ID)
		if accountContactResult.Error == nil {
			accountContactData := accountContactResult.Result.(sharedModel.B2BAccountContact)
			if strings.ToUpper(*accountContactData.Status) == strings.ToUpper(sharedModel.StatusNeedReviewAccount) {
				return errors.New(model.ErrorAccountInActiveBahasa)
			} else if strings.ToUpper(*accountContactData.Status) == strings.ToUpper(sharedModel.StatusDeactiveAccount) || accountContactData.IsDisabled {
				return errors.New(model.ErrorAccountDeactiveBahasa)
			} else if accountContactData.IsDelete {
				return errors.New(model.ErrorInvalidUsernameOrPasswordBahasa)
			}
			member.AccountID = *accountContactData.AccountID
		} else {
			return errors.New(model.ErrorInvalidUsernameOrPasswordBahasa)
		}
	}

	return nil
}

// validateMemberLoginPasswordMicrosite function for validate request token with grantType is password
func (au *AuthUseCaseImpl) validateMemberLoginPasswordMicrosite(member sharedModel.B2BContactData) error {
	if len(member.Status) == 0 || strings.ToUpper(member.Status) != strings.ToUpper(corporateModel.StatusActivated) {
		return errors.New(model.ErrorAccountInActiveBahasa)
	}

	return nil
}
