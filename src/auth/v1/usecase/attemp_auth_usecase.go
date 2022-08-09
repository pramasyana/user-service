package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

// checkAttempt function for saving login attempt data to redis based on email for 24 hours
func (au *AuthUseCaseImpl) checkAttempt(ctxReq context.Context, email string) <-chan ResultUseCase {
	ctx := "AuthUseCase-checkAttempt"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		key := fmt.Sprintf(keyAttempt, email)
		tags[helper.TextParameter] = key
		attemptAge, err := time.ParseDuration(au.LoginAttemptAge)
		if err != nil {
			errC := errors.New("failed to check login attempt")
			output <- ResultUseCase{Error: errC}
			return
		}

		attempt, err := au.checkAttemptCount(ctxReq, email, attemptAge, key)
		if err != nil {
			output <- ResultUseCase{Error: err}
			return
		}

		strAttempt := strconv.Itoa(attempt)
		newData := model.LoginAttempt{
			Key:             key,
			Attempt:         strAttempt,
			LoginAttemptAge: attemptAge,
		}

		saveResult := <-au.LoginAttemptRepo.Save(ctxReq, &newData)
		if saveResult.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_login_attemp", err, email)
			output <- ResultUseCase{Error: saveResult.Error}
			return
		}

		output <- ResultUseCase{Error: errors.New("failed to login")}

	})

	return output
}

// checkAttemptCount function for get attempt count
func (au *AuthUseCaseImpl) checkAttemptCount(ctxReq context.Context, email string, attemptAge time.Duration, key string) (int, error) {
	ctx := "AuthUseCase-checkAttemptCount"
	attempt := 1
	// get attempt first then increment the number
	getAttempt := <-au.getAttempt(ctxReq, key, attempt, attemptAge)
	if getAttempt.Error != nil {
		if getAttempt.Error.Error() == helper.ErrorRedis {
			helper.SendErrorLog(ctxReq, ctx, "parse_login_attempt", getAttempt.Error, email)
		}

		err := errors.New("result is not login attempt")
		return attempt, err
	}

	data := getAttempt.Result.(model.LoginAttempt)
	if len(data.Attempt) > 0 && data.Attempt != "0" {
		intAttempt, _ := strconv.Atoi(data.Attempt)

		// by pass when attempt reach 10 times
		if intAttempt == 9 {
			updateResult := <-au.MemberQueryWrite.UpdateBlockedMember(email)
			if updateResult.Error != nil {
				helper.SendErrorLog(ctxReq, ctx, "update_member_login_attempt", updateResult.Error, email)
				return attempt, updateResult.Error
			}

			err := errors.New(model.ErrorAccountBlockedBahasa)
			return attempt, err
		}
		attempt = intAttempt + 1
	}
	return attempt, nil
}

// getAttempt function for get attempt exist
func (au *AuthUseCaseImpl) getAttempt(ctxReq context.Context, key string, attempt int, attemptAge time.Duration) <-chan ResultUseCase {
	ctx := "AuthUseCase-getAttempt"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		attResult := <-au.LoginAttemptRepo.Load(ctxReq, key)
		if attResult.Error != nil {
			if attResult.Error.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "get_login_attempt", attResult.Error, key)
				output <- ResultUseCase{Error: attResult.Error}
				return
			}

			// when redis nil
			strAttempt := strconv.Itoa(attempt)
			newData := model.LoginAttempt{
				Key:             key,
				Attempt:         strAttempt,
				LoginAttemptAge: attemptAge,
			}

			tags[helper.TextParameter] = newData
			saveResult := <-au.LoginAttemptRepo.Save(ctxReq, &newData)
			if saveResult.Error != nil {
				helper.SendErrorLog(ctxReq, ctx, "save_login_attempt", saveResult.Error, newData)
				output <- ResultUseCase{Error: saveResult.Error}
				return
			}

			output <- ResultUseCase{Error: attResult.Error}
			return
		}

		data := attResult.Result.(model.LoginAttempt)
		output <- ResultUseCase{Result: data}
	})
	return output
}
