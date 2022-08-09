package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/labstack/echo"
)

// VerifyCaptcha usecase function for verify email member
func (au *AuthUseCaseImpl) VerifyCaptcha(ctxReq context.Context, data model.GoogleCaptcha) <-chan ResultUseCase {
	ctx := "AuthQuery-VerifyCaptcha"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		uri := fmt.Sprintf("%s?secret=%s&response=%s&remoteip=%s", au.GoogleVerifyCaptchaURL.String(), data.Secret, data.Response, data.RemoteIP)
		tags["uri"] = uri
		if _, err := url.Parse(uri); err != nil {
			output <- ResultUseCase{Error: err}
			return
		}

		headers := map[string]string{
			echo.HeaderContentType: echo.MIMEApplicationJSON,
		}

		resp := model.GoogleCaptchaResponse{}

		if err := helper.GetHTTPNewRequest(ctxReq, "POST", uri, nil, &resp, headers); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_captcha_response", err, data)
			output <- ResultUseCase{Error: err}
			return
		}

		responseResult := model.GoogleCaptchaResponseResult{}
		responseResult.Success = resp.Success
		responseResult.Score = resp.Score
		responseResult.Action = resp.Action
		responseResult.ChallengeTs = resp.ChallengeTs
		responseResult.HostName = resp.HostName
		responseResult.ErrorCodes = resp.ErrorCodes

		if !resp.Success {
			err := errors.New("failed to verify captcha")
			output <- ResultUseCase{Result: responseResult, Error: err}
			return
		}

		tags[helper.TextResponse] = responseResult
		output <- ResultUseCase{Result: responseResult}
	})

	return output
}
