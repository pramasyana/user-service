package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared"
)

const (
	// TokenContextKey context token key
	TokenContextKey = shared.ContextKey("token")
)

// UploadService data structure
type UploadService struct {
	BaseURL    *url.URL
	AccessorID string
	Timeout    string
}

// NewUploadService function for initializing upload service
func NewUploadService() (*UploadService, error) {
	var (
		upload UploadService
		err    error
		ok     bool
	)

	baseURL, ok := os.LookupEnv("UPLOAD_SERVICE_URL")
	if !ok {
		return &upload, errors.New("you need to specify UPLOAD_SERVICE_URL in the environment variable")
	}

	accessorID, ok := os.LookupEnv("UPLOAD_SERVICE_ACCESSOR_ID")
	if !ok {
		return &upload, errors.New("you need to specify UPLOAD_SERVICE_URL in the environment variable")
	}

	timeout, ok := os.LookupEnv("UPLOAD_SERVICE_TIMEOUT")
	if !ok {
		return &upload, errors.New("you need to specify UPLOAD_SERVICE_TIMEOUT in the environment variable")
	}

	upload.BaseURL, err = url.Parse(baseURL)
	if err != nil {
		err := errors.New("error parsing upload services url")
		return &upload, err
	}
	upload.AccessorID = accessorID
	upload.Timeout = timeout

	return &upload, nil
}

// GetURLImage function for getting detail url image by url
func (p *UploadService) GetURLImage(ctxReq context.Context, key string, isAttachment string) <-chan serviceModel.ServiceResult {
	ctx := "Service-GetURLImage"
	var result = make(chan serviceModel.ServiceResult)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		resp := serviceModel.ResponseUploadService{}
		token, _ := shared.GetDataFromContext(ctxReq, TokenContextKey).(string)
		tags[helper.TextToken] = token
		tags[helper.TextURL] = key
		if token == "" {
			err := errors.New("token empty")
			helper.SendErrorLog(ctxReq, ctx, "get_token_from_context", err, nil)
			tags[helper.TextStmtError] = err
			result <- serviceModel.ServiceResult{Error: err}
			return
		}

		key, ok := p.replaceURLFile(key)
		if !ok {
			result <- serviceModel.ServiceResult{Result: nil}
			return
		}

		// generate uri
		uri := fmt.Sprintf("%s%s%s%s%s%s%s", p.BaseURL.String(), "/presigned?timeout=", p.Timeout, "&key=", key, "&isAttachment=", isAttachment)
		if _, err := url.Parse(uri); err != nil {
			result <- serviceModel.ServiceResult{Result: nil, Error: err}
			return
		}

		// generate headers
		headers := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": token,
			"Accessor-ID":   p.AccessorID,
		}

		tags[helper.TextParameter] = headers

		// request
		err := helper.GetHTTPNewRequestV2(ctxReq, "GET", uri, nil, &resp, headers)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "request_to_upload", err, headers)
			tags[helper.TextStmtError] = err
			err := errors.New("failed get upload service")
			result <- serviceModel.ServiceResult{Error: err}
			return
		}

		result <- serviceModel.ServiceResult{Result: resp}
	})

	return result
}

// replaceURLFile function for replace url if contains full url
func (p *UploadService) replaceURLFile(url string) (string, bool) {

	// check if the url from static is not from the upload service
	if strings.Contains(url, "static.bmdstatic.com") {
		return url, false
	}

	AWSMerchantDocumentURL, _ := os.LookupEnv("AWS_MERCHANT_DOCUMENT_URL")
	contains := strings.Replace(AWSMerchantDocumentURL, helper.TextHTTPS, "", -1)
	if strings.Contains(url, contains) {
		url = strings.Replace(url, AWSMerchantDocumentURL, "", -1)
		return url, true
	}

	AWSMerchantDocumentURLSalmon, ok := os.LookupEnv("AWS_MERCHANT_DOCUMENT_URL_SALMON")
	if ok {
		containsSalmon := strings.Replace(AWSMerchantDocumentURLSalmon, "https://", "", -1)
		if strings.Contains(url, containsSalmon) {
			url = strings.Replace(url, AWSMerchantDocumentURLSalmon, "", -1)
			return url, true
		}
	}

	return url, true
}
