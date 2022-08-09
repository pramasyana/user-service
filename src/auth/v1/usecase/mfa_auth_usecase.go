package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	dgoogauth "github.com/dgryski/dgoogauth"
)

// verifyMFACode function for verify mfa otp
func (au *AuthUseCaseImpl) verifyMFACode(ctxReq context.Context, mfaKey string, otp string) <-chan ResultUseCase {
	ctx := "AuthUseCase-verifyMFACode"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags["args"] = mfaKey
		// only for exclude prod & key static
		if os.Getenv("ENV") != "PROD" && mfaKey == memberModel.StaticSharedMfaKeyForDev && otp == memberModel.StaticOTPMfaForDev {
			output <- ResultUseCase{Result: true}
			return
		}

		// The OTPConfig gets modified by otpc.Authenticate() to prevent passcode replay, etc.,
		// so allocate it once and reuse it for multiple calls.
		otpc := &dgoogauth.OTPConfig{
			Secret:      mfaKey,
			WindowSize:  3,
			HotpCounter: 0,
			UTC:         true,
		}

		// authentication OTPConfig
		val, err := otpc.Authenticate(otp)
		if err != nil {
			tags["response"] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if !val {
			err := errors.New(memberModel.ErrorMFAOTP)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
	})

	return output
}

// parseTokenMFA function for parsing token into email and real token
func (au *AuthUseCaseImpl) parseTokenMFA(token string) (string, string, error) {
	if len(token) == 0 {
		err := fmt.Errorf(helper.ErrorParameterRequired, "token")
		return "", "", err
	}

	splitTokens := strings.Split(token, "-")
	if len(splitTokens) < 2 {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "token")
		return "", "", err
	}

	memberID, err := base64.URLEncoding.DecodeString(splitTokens[1])
	if err != nil {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "token")

		return "", "", err
	}

	return string(memberID), splitTokens[0], nil
}
