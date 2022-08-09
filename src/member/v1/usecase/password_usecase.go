package usecase

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	corporateModel "github.com/Bhinneka/user-service/src/corporate/v2/model"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/golang-jwt/jwt"
)

// UpdatePassword function for updating password
func (mu *MemberUseCaseImpl) UpdatePassword(ctxReq context.Context, token, uid, oldPassword, newPassword string) <-chan ResultUseCase {
	ctx := "MemberUseCase-UpdatePassword"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// validate request
		params := model.PayloadUpdate{
			OldPassword: oldPassword,
			NewPassword: newPassword,
		}
		tags[helper.TextMemberIDCamel] = uid

		memberOld, member, err := mu.validateUpdatePassword(ctxReq, params, uid)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, member)
		if saveResult.Error != nil {
			output <- ResultUseCase{Error: saveResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		// logout all session login
		err = mu.revokeAllAccessProccess(ctxReq, member.ID, token, false)
		if err != nil {
			tags[helper.TextResponse] = err
		}

		response := model.SuccessResponse{
			ID:          member.ID,
			Message:     helper.SuccessMessage,
			HasPassword: true,
			Email:       strings.ToLower(member.Email),
			FirstName:   member.FirstName,
			LastName:    member.LastName,
		}
		plLog := model.MemberLog{
			Before: &memberOld,
			After:  &member,
		}

		go mu.QPublisher.QueueJob(ctxReq, plLog, member.ID, "InsertLogUpdateMember")

		plEmail := model.MemberEmailQueue{
			Member: &member,
			Data:   response,
		}
		go mu.QPublisher.QueueJob(ctxReq, plEmail, member.ID, "SendEmailSuccessForgotPassword")

		output <- ResultUseCase{Result: response}
	})

	return output
}

func (mu *MemberUseCaseImpl) validateUpdatePassword(ctxReq context.Context, params model.PayloadUpdate, uid string) (model.Member, model.Member, error) {
	var (
		memberVUP    model.Member
		memberOldVUP model.Member
	)
	err := jsonschema.ValidateTemp("change_password_params", params)
	if err != nil {
		return memberOldVUP, memberVUP, err
	}

	if !strings.Contains(uid, usrFormat) {
		err := fmt.Errorf(helper.ErrorParameterInvalid, msgErrorMemberID)
		return memberOldVUP, memberVUP, err
	}

	memberResultVUP := <-mu.MemberRepoRead.Load(ctxReq, uid)
	if memberResultVUP.Error != nil {
		if memberResultVUP.Error == sql.ErrNoRows {
			memberResultVUP.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
		}

		return memberOldVUP, memberVUP, memberResultVUP.Error
	}

	memberVUP, ok := memberResultVUP.Result.(model.Member)
	if !ok {
		err := errors.New(msgErrorResultMember)
		return memberOldVUP, memberVUP, err
	}

	memberOldVUP = memberVUP

	// validate password format
	if err := helper.ValidatePassword(params.NewPassword); err != nil {
		return memberOldVUP, memberVUP, err
	}

	// matching old password with new password
	if params.OldPassword == params.NewPassword {
		err := errors.New("new password and previous password cannot be same")
		return memberOldVUP, memberVUP, err
	}

	passwordSaltDB := strings.TrimSpace(memberVUP.Salt)
	mu.Hash.ParseSalt(passwordSaltDB)
	base64Data := base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(params.OldPassword)))

	if memberVUP.Password != base64Data {
		err := errors.New(model.ErrorOldPasswordInvalid)
		return memberOldVUP, memberVUP, err
	}

	// encode the new password then replace the old password and salt
	memberVUP.Salt = mu.Hash.GenerateSalt()
	err = mu.Hash.ParseSalt(memberVUP.Salt)
	if err != nil {
		return memberOldVUP, memberVUP, err
	}

	memberVUP.Password = base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(params.NewPassword)))
	memberVUP.LastPasswordModified = time.Now()

	return memberOldVUP, memberVUP, nil
}

func (mu *MemberUseCaseImpl) validateUpdateSyncPassword(ctxReq context.Context, params model.PayloadUpdate, email, memberType string) (model.Member, model.Member, error) {
	var (
		member    model.Member
		memberOld model.Member
	)
	err := jsonschema.ValidateTemp("change_password_params", params)
	if err != nil {
		return memberOld, member, err
	}

	memberResult := <-mu.MemberQueryRead.FindByEmail(ctxReq, email)
	if memberResult.Error != nil {
		if memberResult.Error == sql.ErrNoRows {
			memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
		}

		return memberOld, member, memberResult.Error
	}

	member, ok := memberResult.Result.(model.Member)
	if !ok {
		err := errors.New(msgErrorResultMember)
		return memberOld, member, err
	}

	memberOld = member

	// validate password format
	if err := helper.ValidatePassword(params.NewPassword); err != nil {
		return memberOld, member, err
	}

	// matching old password with new password
	if params.OldPassword == params.NewPassword {
		err := errors.New("new password and previous password cannot be same")
		return memberOld, member, err
	}

	passwordSaltDB := strings.TrimSpace(member.Salt)
	mu.Hash.ParseSalt(passwordSaltDB)
	base64Data := base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(params.OldPassword)))

	if memberOld.Password != "" {
		if member.Password != base64Data {
			err := errors.New(model.ErrorOldPasswordInvalid)
			return memberOld, member, err
		}
	}

	// encode the new password then replace the old password and salt
	member.Salt = mu.Hash.GenerateSalt()
	err = mu.Hash.ParseSalt(member.Salt)
	if err != nil {
		return memberOld, member, err
	}

	member.Password = base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(params.NewPassword)))
	member.LastPasswordModified = time.Now()

	return memberOld, member, nil
}

// AddNewPassword function for adding new password
func (mu *MemberUseCaseImpl) AddNewPassword(ctxReq context.Context, data model.Member) <-chan ResultUseCase {
	ctx := "MemberUseCase-AddNewPassword"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		errorNewPassword := mu.validateNewPassword(data)
		if errorNewPassword != nil {
			tags[helper.TextResponse] = errorNewPassword
			output <- ResultUseCase{Error: errorNewPassword, HTTPStatus: http.StatusBadRequest}
			return
		}

		// encode the new password then replace the old password and salt
		data.Salt = mu.Hash.GenerateSalt()
		err := mu.Hash.ParseSalt(data.Salt)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		data.Password = base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(data.NewPassword)))
		data.LastPasswordModified = time.Now()

		memberResult := <-mu.MemberRepoWrite.Save(ctxReq, data)
		if memberResult.Error != nil {
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusInternalServerError}
			return
		}

		// check b2b_contact & get data
		mu.SyncPasswordContact(ctxReq, data)

		// flush existing token user
		mu.flushAllTokenUser(ctxReq, data)

	})

	return output
}

// validateNewPassword function for validate new password
func (mu *MemberUseCaseImpl) validateNewPassword(data model.Member) error {
	// validate request
	params := model.PayloadUpdate{
		NewPassword: data.NewPassword,
		RePassword:  data.RePassword,
	}

	mErr := jsonschema.ValidateTemp("add_password_params", params)
	if mErr != nil {
		return mErr
	}
	// matching password and password confirmation
	if data.NewPassword != data.RePassword {
		err := errors.New(msgErrorPasswordMatch)
		return err
	}

	// validate password format
	if err := helper.ValidatePassword(data.NewPassword); err != nil {
		return err
	}

	return nil
}

func (mu *MemberUseCaseImpl) ClaimsToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token = strings.Replace(token, "Bearer ", "", -1)
	jwtResult, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if (jwtResult == nil && err != nil) || len(claims) == 0 {
		return claims, errors.New("invalid token")
	}

	return claims, nil
}

func (mu *MemberUseCaseImpl) SyncPassword(ctxReq context.Context, token, oldPassword, newPassword string) <-chan ResultUseCase {
	ctx := "MemberUseCase-SyncPassword"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		claims, err := mu.ClaimsToken(token)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// validate request
		paramsSP := model.PayloadUpdate{
			OldPassword: oldPassword,
			NewPassword: newPassword,
		}
		tags[helper.TextMemberIDCamel] = claims["sub"].(string)

		// check member & validate password
		memberOldSP, memberSP, errSP := mu.validateUpdateSyncPassword(ctxReq, paramsSP, claims["email"].(string), claims["memberType"].(string))
		if errSP != nil {
			tags[helper.TextResponse] = errSP
			output <- ResultUseCase{Error: errSP, HTTPStatus: http.StatusBadRequest}
			return
		}

		// check b2b_contact & get data
		corporateContactSP := <-mu.CorporateContactQueryRead.FindContactByEmail(ctxReq, memberSP.Email)
		if corporateContactSP.Error != nil {
			errSP := errors.New("user not exists in contact")
			tags[helper.TextResponse] = errSP
			output <- ResultUseCase{Error: errSP, HTTPStatus: http.StatusBadRequest}
			return
		}

		contactSP, ok := corporateContactSP.Result.(sharedModel.B2BContactData)
		if !ok {
			errSP := errors.New("contact not result")
			tags[helper.TextResponse] = errSP
			output <- ResultUseCase{Error: errSP, HTTPStatus: http.StatusBadRequest}
			return
		}

		// check password bbcom if bcom password not set
		if memberOldSP.Password == "" {
			passwordSaltDB := strings.TrimSpace(contactSP.Salt)
			mu.Hash.ParseSalt(passwordSaltDB)
			base64Data := base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(paramsSP.OldPassword)))

			if contactSP.Password != base64Data {
				errSP := errors.New(model.ErrorOldPasswordInvalid)
				tags[helper.TextResponse] = errSP
				output <- ResultUseCase{Error: errSP, HTTPStatus: http.StatusBadRequest}
				return
			}
		}

		// save member
		saveResultSP := <-mu.MemberRepoWrite.Save(ctxReq, memberSP)
		if saveResultSP.Error != nil {
			output <- ResultUseCase{Error: saveResultSP.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		// logout all session login
		errSP = mu.revokeAllAccessProccess(ctxReq, memberSP.ID, token, false)
		if errSP != nil {
			tags[helper.TextResponse] = errSP
		}

		// update flag is_sync true member
		saveFlagResultSP := <-mu.MemberRepoRead.UpdateFlagIsSyncMember(ctxReq, memberSP)
		if saveFlagResultSP.Error != nil {
			tags[helper.TextResponse] = saveFlagResultSP.Error
			output <- ResultUseCase{Error: saveFlagResultSP.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		// send to kafka contact
		transactionType, _ := json.Marshal(contactSP.TransactionType)
		payload := corporateModel.ContactPayload{
			ID:              contactSP.ID,
			Email:           memberSP.Email,
			FirstName:       contactSP.FirstName,
			LastName:        contactSP.LastName,
			PhoneNumber:     contactSP.PhoneNumber,
			Password:        memberSP.Password,
			Salt:            memberSP.Salt,
			TransactionType: string(transactionType),
			CreatedAt:       time.Now(),
		}
		go mu.PublishToKafkaContact(ctxReq, payload, "syncb2b")

		output <- ResultUseCase{Error: nil}
	})

	return output
}

func (mu *MemberUseCaseImpl) ChangePassword(ctxReq context.Context, token, oldPassword, newPassword string) <-chan ResultUseCase {
	ctx := "MemberUseCase-ChangePassword"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		claims, err := mu.ClaimsToken(token)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// validate request
		params := model.PayloadUpdate{
			OldPassword: oldPassword,
			NewPassword: newPassword,
		}
		tags[helper.TextMemberIDCamel] = claims["sub"].(string)

		// check member
		_, member, err := mu.validateUpdateSyncPassword(ctxReq, params, claims["email"].(string), claims["memberType"].(string))
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, member)
		if saveResult.Error != nil {
			output <- ResultUseCase{Error: saveResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		// logout all session login
		err = mu.revokeAllAccessProccess(ctxReq, member.ID, token, false)
		if err != nil {
			tags[helper.TextResponse] = err
		}

		// check b2b_contact & get data
		mu.SyncPasswordContact(ctxReq, member)

		response := model.SuccessResponse{
			ID:          member.ID,
			Message:     helper.SuccessMessage,
			HasPassword: true,
			Email:       strings.ToLower(member.Email),
			FirstName:   member.FirstName,
			LastName:    member.LastName,
		}

		output <- ResultUseCase{Result: response}
	})

	return output
}

func (mu *MemberUseCaseImpl) SyncPasswordContact(ctxReq context.Context, member model.Member) {
	ctx := "MemberUseCase-SyncPasswordContact"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		corporateContact := <-mu.CorporateContactQueryRead.FindContactByEmail(ctxReq, member.Email)
		if corporateContact.Error != nil {
			err := errors.New("user not exists in contact")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		contact, ok := corporateContact.Result.(sharedModel.B2BContactData)
		if !ok {
			err := errors.New("contact not result")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// check is_sync is true
		if contact.IsSync {
			// send to kafka contact
			transactionType, _ := json.Marshal(contact.TransactionType)
			payload := corporateModel.ContactPayload{
				ID:              contact.ID,
				Email:           member.Email,
				FirstName:       contact.FirstName,
				LastName:        contact.LastName,
				PhoneNumber:     contact.PhoneNumber,
				Password:        member.Password,
				Salt:            member.Salt,
				TransactionType: string(transactionType),
				CreatedAt:       time.Now(),
			}
			go mu.PublishToKafkaContact(ctxReq, payload, "syncb2b")
		}
	})
}
