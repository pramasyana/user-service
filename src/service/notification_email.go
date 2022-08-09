package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	staticshared "github.com/Bhinneka/user-service/src/shared/static"
	"github.com/labstack/echo"
	"github.com/spf13/cast"
)

// GetTemplateByID function for getting detail template by id
func (em *NotificationService) GetTemplateByID(ctxReq context.Context, templateID, envKey string) <-chan serviceModel.ServiceResult {
	ctx := "NotifServiceData-GetTemplateByID"
	var result = make(chan serviceModel.ServiceResult)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		resp := serviceModel.ResponseGetTemplate{}
		tokenWithBearer, _ := em.auth(ctxReq)

		uri := fmt.Sprintf("%s%s%s", em.BaseURL.String(), "/template/", templateID)
		tags[helper.TextURL] = uri

		headers := map[string]string{
			echo.HeaderAuthorization: tokenWithBearer,
		}

		// request
		err := helper.GetHTTPNewRequestV2(ctxReq, http.MethodGet, uri, nil, &resp, headers)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "http_request_to_notification", err, uri)
			resp.Data.Content = staticshared.GetFallbackEmailContent(envKey)
			result <- serviceModel.ServiceResult{Result: resp.Data}
			return
		}
		tags[helper.TextResponse] = resp

		result <- serviceModel.ServiceResult{Result: resp.Data}
	})

	return result
}

// SendEmail function for send
func (em *NotificationService) SendEmail(ctxReq context.Context, email serviceModel.Email) (string, error) {
	ctx := "NotifServiceData-SendEmail"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextParameter: email,
		helper.TextEmail:     email.To[0],
	}
	defer tr.Finish(tags)

	var status string

	// get auth first
	tokenWithBearer, err := em.auth(ctxReq)
	if err != nil {
		return status, err
	}

	// set payload
	payload := serviceModel.PayloadEmail{}
	payload.Data.Attributes = email
	jsonPayload, _ := json.Marshal(payload)

	method := echo.POST
	url := fmt.Sprintf("%s/email/send", em.BaseURL.String())
	resp := &serviceModel.SuccessMessage{}
	respError := &serviceModel.ErrorMessage{}

	// set connection
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "init_http_request", err, payload)
		return "", err
	}
	req.Header.Set(echo.HeaderContentType, "application/vnd.api+json")
	req.Header.Set(echo.HeaderAuthorization, tokenWithBearer)

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "do_http_request", err, payload)
		return "", err
	}
	tags[helper.TextResponse] = resp
	defer r.Body.Close()

	if r.StatusCode == http.StatusUnauthorized || r.StatusCode == http.StatusBadRequest {
		if err = json.NewDecoder(r.Body).Decode(respError); err != nil {
			return status, err
		}

		return status, errors.New(respError.Errors[0].Detail)
	} else if r.StatusCode == http.StatusOK {
		if err = json.NewDecoder(r.Body).Decode(resp); err != nil {
			return status, err
		}

		status = resp.Data.Attributes.Message
	} else {
		err = fmt.Errorf("error code : %s", cast.ToString(r.StatusCode))
		helper.SendErrorLog(ctxReq, ctx, "error_other", err, nil)
	}

	return status, nil
}
