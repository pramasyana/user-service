package usecase

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	goString "github.com/Bhinneka/golib/string"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
)

// ForgotPassword function for requesting token of forgot password
func (mu *MemberUseCaseImpl) ForgotPassword(ctxReq context.Context, email string) <-chan ResultUseCase {
	ctx := "MemberUseCase-ForgotPassword"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		member, err := mu.validateForgotPassword(ctxReq, email)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// get data first from redis when it does not exist
		tokenRedisResult := <-mu.MemberRepoRedis.Load(member.ID)
		if tokenRedisResult.Error != nil && tokenRedisResult.Error.Error() != helper.ErrorRedis {
			output <- ResultUseCase{Error: tokenRedisResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		// check member password
		if len(member.Password) == 0 {
			err := errors.New("you do not have password before")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		tokenRedis, ok := tokenRedisResult.Result.(model.MemberRedis)
		if !ok || tokenRedis.Count < 10 {
			err := errors.New("result is not member from redis or counter less than 10")
			tags[helper.TextResponse] = err

			mix := member.Email + "-" + member.ID
			tokenRedis.Token = helper.GenerateTokenByString(mix)
			tokenRedis.ID = member.ID
			tokenRedis.Count++
			tokenRedis.TTL = 24 * 3600 * time.Second

			saveResult := <-mu.MemberRepoRedis.Save(&tokenRedis)
			if saveResult.Error != nil {
				output <- ResultUseCase{Error: saveResult.Error, HTTPStatus: http.StatusBadRequest}
				return
			}

			// just for returning to member
			base64EncodeEmail := base64.URLEncoding.EncodeToString([]byte(email))
			combined := tokenRedis.Token + "-" + base64EncodeEmail

			response := model.SuccessResponse{
				ID:          member.ID,
				Message:     helper.SuccessMessage,
				Token:       combined,
				HasPassword: true,
				Email:       strings.ToLower(member.Email),
				FirstName:   member.FirstName,
				LastName:    member.LastName,
			}

			plEmail := model.MemberEmailQueue{
				Member: &member,
				Data:   response,
			}

			go mu.QPublisher.QueueJob(ctxReq, plEmail, member.ID, "SendEmailForgotPassword")

			output <- ResultUseCase{Result: response}
			return
		}

		output <- ResultUseCase{Error: errors.New("you have already done the same things 10 times"), HTTPStatus: http.StatusBadRequest}

	})

	return output
}

// validateForgotPassword function for validate forgot password
func (mu *MemberUseCaseImpl) validateForgotPassword(ctxReq context.Context, email string) (model.Member, error) {
	member := model.Member{}
	if len(email) == 0 {
		err := fmt.Errorf(helper.ErrorParameterRequired, "email")
		return member, err
	}

	if err := goString.ValidateEmail(email); err != nil {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "email")
		return member, err
	}

	memberResult := <-mu.MemberQueryRead.FindByEmail(ctxReq, email)
	if memberResult.Error != nil {
		// when data not found
		if memberResult.Error == sql.ErrNoRows {
			memberResult.Error = fmt.Errorf("Email belum terdaftar")
		}

		return member, memberResult.Error
	}

	member, ok := memberResult.Result.(model.Member)
	if !ok {
		err := errors.New(msgErrorResultMember)
		return member, err
	}
	return member, nil
}

// ChangeForgotPassword function for replacing old forgotten password
func (mu *MemberUseCaseImpl) ChangeForgotPassword(ctxReq context.Context, token, newPassword, rePassword, requestFrom string) <-chan ResultUseCase {
	ctx := "MemberUseCase-ChangeForgotPassword"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// validate password and its confirmation
		if err := mu.validateChangeForgotPassword(token, newPassword, rePassword); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		email, newToken, err := mu.parseToken(token)
		tags[helper.TextParameter] = email
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, errCode, err := mu.getMemberForgotPassword(ctxReq, email, newToken)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: errCode}
			return
		}

		// generate salt and password
		// encode the new password then replace the old password and salt
		member.Salt = mu.Hash.GenerateSalt()
		err = mu.Hash.ParseSalt(member.Salt)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		member.Password = base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(newPassword)))
		member.LastPasswordModified = time.Now()

		memberSaveResult := <-mu.MemberRepoWrite.Save(ctxReq, member)
		if memberSaveResult.Error != nil {
			output <- ResultUseCase{Error: memberSaveResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		// delete redis key after saving new password
		<-mu.MemberRepoRedis.Delete(member.ID)

		response := model.SuccessResponse{
			ID:          golib.RandomString(8),
			Message:     helper.SuccessMessage,
			HasPassword: true,
			Email:       strings.ToLower(email),
			FirstName:   member.FirstName,
			LastName:    member.LastName,
		}

		// check b2b_contact & get data
		mu.SyncPasswordContact(ctxReq, member)

		// flush existing token user
		mu.flushAllTokenUser(ctxReq, member)
		if requestFrom == model.Sturgeon {
			plEmail := model.MemberEmailQueue{
				Member: &member,
				Data:   response,
			}
			go mu.QPublisher.QueueJob(ctxReq, plEmail, member.ID, "SendEmailSuccessForgotPassword")
		}

		output <- ResultUseCase{Result: response}

	})

	return output
}

// getMemberForgotPassword function for get member detail forgot password
func (mu *MemberUseCaseImpl) getMemberForgotPassword(ctxReq context.Context, email, newToken string) (model.Member, int, error) {
	member := model.Member{}
	memberResult := <-mu.MemberQueryRead.FindByEmail(ctxReq, email)
	if memberResult.Error != nil {
		// when data not found
		if memberResult.Error == sql.ErrNoRows {
			memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, msgErrorMember+email)
		}

		return member, http.StatusBadRequest, memberResult.Error
	}

	member, ok := memberResult.Result.(model.Member)
	if !ok {
		err := errors.New(msgErrorResultMember)
		return member, http.StatusBadRequest, err
	}

	// check password existence first before replacing the forgotten password
	// this process can only be done by member who has password
	if len(member.Password) == 0 {
		err := errors.New("you do not have password before")
		return member, http.StatusBadRequest, err
	}

	tokenResult := <-mu.MemberRepoRedis.Load(member.ID)
	if tokenResult.Error != nil {
		err := errors.New(msgErrorResultMember)
		return member, http.StatusBadRequest, err
	}

	tokenRedis, ok := tokenResult.Result.(model.MemberRedis)
	if !ok {
		err := errors.New("result is not member from redis")
		return member, http.StatusBadRequest, err
	}

	// validate the token from redis and form
	if tokenRedis.Token != newToken {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "token")
		return member, http.StatusUnauthorized, err
	}

	return member, http.StatusOK, nil
}

// validateChangeForgotPassword function for validate change password
func (mu *MemberUseCaseImpl) validateChangeForgotPassword(token, newPassword, rePassword string) error {
	if len(newPassword) == 0 {
		err := fmt.Errorf(helper.ErrorParameterRequired, "new password")
		return err
	}

	// validate password format
	if err := helper.ValidatePassword(newPassword); err != nil {
		return err
	}

	if len(rePassword) == 0 {
		err := fmt.Errorf(helper.ErrorParameterRequired, "confirmation password")
		return err
	}

	if newPassword != rePassword {
		err := errors.New(msgErrorPasswordMatch)
		return err
	}

	return nil
}

// ActivateNewPassword function for activating new password for member who is registered from dolphin
func (mu *MemberUseCaseImpl) ActivateNewPassword(ctxReq context.Context, token, password, rePassword string) <-chan ResultUseCase {
	ctx := "MemberUseCase-ActivateNewPassword"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// validate data first
		member, err := mu.validateActiveNewPassword(ctxReq, token, password, rePassword)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// generate salt and password
		// encode the new password then replace the old password and salt
		member.Salt = mu.Hash.GenerateSalt()
		err = mu.Hash.ParseSalt(member.Salt)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		member.Password = base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(password)))
		member.LastPasswordModified = time.Now()

		// append the new data
		member.Status = model.StringToStatus(model.ActiveString)
		member.Token = ""

		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, member)
		if saveResult.Error != nil {
			err := errors.New(msgErrorSaveMember)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		response := model.SuccessResponse{
			ID:          golib.RandomString(8),
			Message:     helper.SuccessMessage,
			HasPassword: true,
			Email:       strings.ToLower(member.Email),
			FirstName:   member.FirstName,
			LastName:    member.LastName,
		}
		//Publish Kafka if From Dolphin
		if member.SignUpFrom == model.Dolphin {
			// find member by token given
			memberResult := <-mu.MemberQueryRead.FindByEmail(ctxReq, member.Email)
			if memberResult.Error != nil {
				err := errors.New(msgErrorSaveMember)
				output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
				return
			}

			member, ok := memberResult.Result.(model.Member)
			if !ok {
				err := errors.New(msgErrorSaveMember)
				output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
				return
			}
			go mu.PublishToKafkaUser(ctxReq, &member, "activation")
		}

		// flush existing token user
		mu.flushAllTokenUser(ctxReq, member)

		output <- ResultUseCase{Result: response}

	})

	return output
}

// validateActiveNewPassword function for validate new password
func (mu *MemberUseCaseImpl) validateActiveNewPassword(ctxReq context.Context, token, password, rePassword string) (model.Member, error) {
	member := model.Member{}
	if len(token) == 0 {
		err := fmt.Errorf(helper.ErrorParameterRequired, "token")
		return member, err
	}
	if len(password) == 0 {
		err := fmt.Errorf(helper.ErrorParameterRequired, "new password")
		return member, err
	}
	if password != rePassword {
		err := errors.New(msgErrorPasswordMatch)
		return member, err
	}

	// validate password format
	if err := helper.ValidatePassword(password); err != nil {
		return member, err
	}

	// find member by token given
	memberResult := <-mu.MemberQueryRead.FindByToken(ctxReq, token)
	if memberResult.Error != nil {
		// when data not found
		if memberResult.Error == sql.ErrNoRows {
			memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "member with token "+token)
		}

		return member, memberResult.Error
	}

	member, ok := memberResult.Result.(model.Member)
	if !ok {
		err := errors.New(msgErrorResultMember)
		return member, err
	}

	// validate if token to activete employee
	member, err := mu.changeStatusEmployee(ctxReq, token, member)
	return member, err
}

func (mu *MemberUseCaseImpl) changeStatusEmployee(ctxReq context.Context, token string, member model.Member) (model.Member, error) {
	// change status employee to active
	if len(token) > 0 {
		data, err := hex.DecodeString(token) // convert to byte
		if err == nil {
			err, decode := golib.Decrypt(data, "SECRET")
			if err == nil {
				// merchantId-memberId
				extractToken := strings.Split(decode, "-")

				params := merchantModel.B2CMerchantEmployee{}
				params.MerchantID = extractToken[0]
				params.MemberID = extractToken[1]
				params.Status = helper.TextActive
				params.ModifiedAt = time.Now()
				params.ModifiedBy = &member.ID
				err = mu.MerchantEmployeeRead.ChangeStatus(ctxReq, params)
				if err != nil {
					err := errors.New("Failed change status employee")
					return member, err
				}
			}
		}
	}

	return member, nil
}
