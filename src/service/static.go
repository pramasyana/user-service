package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// StaticService data structure
type StaticService struct {
	BaseURL        *url.URL
	BaseGraphQLURL *url.URL
	BasicAuth      string
}

// NewStaticService function for initializing static service
func NewStaticService() (*StaticService, error) {
	var (
		statics StaticService
		err     error
		ok      bool
	)
	ctx := "NewStaticService"

	statics.BasicAuth, ok = os.LookupEnv("STATIC_SERVICE_AUTH")
	if !ok {
		return &statics, errors.New("you need to specify STATIC_SERVICE_AUTH in the environment variable")
	}

	baseURL, ok := os.LookupEnv("STATIC_SERVICE_URL")
	if !ok {
		return &statics, errors.New("you need to specify STATIC_SERVICE_URL in the environment variable")
	}

	statics.BaseURL, err = url.Parse(baseURL)
	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, "parse_static_url", err, baseURL)
		return &statics, errors.New("error parsing statics services url")
	}

	BaseGraphQLURL, ok := os.LookupEnv("GWS_GRAPHQL_URL")
	if !ok {
		return &statics, errors.New("you need to specify GWS_GRAPHQL_URL in the environment variable")
	}

	statics.BaseGraphQLURL, err = url.Parse(BaseGraphQLURL)
	if err != nil {
		return &statics, errors.New("error parsing gws graphql url")
	}

	return &statics, nil
}

// FindStaticsByID function for getting detail static by id
// Deprecated: moved to notification service for better performance
func (p *StaticService) FindStaticsByID(ctxReq context.Context, id string) <-chan serviceModel.ServiceResult {
	ctx := "Service-FindStaticsByID"
	var result = make(chan serviceModel.ServiceResult)

	go tracer.WithTrace(ctxReq, ctx, nil, func(ctxReq context.Context) {
		resp := serviceModel.StaticData{}
		// generate uri
		uri := fmt.Sprintf("%s%s%s", p.BaseURL.String(), "/v1/statics/", id)

		// generate headers
		headers := map[string]string{
			"Authorization": "Basic " + p.BasicAuth,
		}

		// request
		err := helper.GetHTTPNewRequest(ctxReq, "GET", uri, nil, &resp, headers)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "http_request_to_static", err, uri)
			err := errors.New("failed get static")
			result <- serviceModel.ServiceResult{Error: err}
			return
		}

		result <- serviceModel.ServiceResult{Result: resp}
	})

	return result
}

// FindStaticsGwsByID function for getting detail static by id
// Deprecated: moved to notification service for better performance
func (p *StaticService) FindStaticsGwsByID(ctxReq context.Context, id string) <-chan serviceModel.ServiceResult {
	ctx := "Service-FindStaticsGwsByID"
	var result = make(chan serviceModel.ServiceResult)

	go tracer.WithTrace(ctxReq, ctx, nil, func(ctxReq context.Context) {
		resp := serviceModel.ResponseGWSStatic{}
		uri := fmt.Sprintf("%s%s", p.BaseGraphQLURL.String(), "/graphql")

		staticID, _ := strconv.Atoi(id)

		params := fmt.Sprintf(`{
			"staticId": %d
		}`, staticID)
		jsonData := map[string]string{
			"query": `
				query	getStaticPageById ($staticId: Int!) {
						getStaticPageById(staticId: $staticId) {
							code
							success
							message
							result {
								id
								title
								subTitle
								metaTitle
								metaDescription
								placement
								contentType
								content
								reviveContent
								slug								
								isActive
								created
								lastModified
								zoneId
							}
						}
					}
				`,
			"variables": params,
		}

		gqlMarshalled, _ := json.Marshal(jsonData)
		payload := strings.NewReader(string(gqlMarshalled))
		headers := map[string]string{
			helper.TextAuthorization: "Basic " + p.BasicAuth,
			contentType:              "application/json",
		}

		if err := helper.GetHTTPNewRequestV2(ctxReq, http.MethodPost, uri, payload, &resp, headers); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "http_request", err, payload)
			result <- serviceModel.ServiceResult{Error: err}
			return
		}

		result <- serviceModel.ServiceResult{Result: resp.Data.GetStaticDetail.Result}
	})

	return result
}
