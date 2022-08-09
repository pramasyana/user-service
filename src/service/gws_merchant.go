package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// GetMerchantServiceGraphQL function for getting from merchant service
func (m *MerchantService) GetMerchantServiceGraphQL(ctxReq context.Context, token, merchantID string) (interface{}, error) {
	ctx := "MerchantService-GetMerchantServiceGraphQL"
	resp := serviceModel.ResponseGWSMerchant{}

	uri := fmt.Sprintf("%s%s", m.BaseGraphQLURL.String(), "/graphql")

	params := fmt.Sprintf(`{
		"code": "%s"
	}`, merchantID)
	jsonData := map[string]string{
		"query": `
			query	getMerchantDetail ($code: String) {
					getMerchantDetail(code: $code) {
						code
						success
						message
						result {
							id
							businessType
							code
							name
							isActive
							createdAt
							updatedAt
						}
					}
				}
			`,
		"variables": params,
	}

	gqlMarshalled, _ := json.Marshal(jsonData)
	payload := strings.NewReader(string(gqlMarshalled))

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}
	if err := helper.GetHTTPNewRequestV2(ctxReq, http.MethodPost, uri, payload, &resp, headers); err != nil {
		helper.SendErrorLog(ctxReq, ctx, "get_http_request", err, resp)
		return "", errors.New(helper.ErrorStatusCode)
	}

	if resp.Data.GetMerchantDetail.Code != 200 || !resp.Data.GetMerchantDetail.Success {
		helper.SendErrorLog(ctxReq, ctx, "decode_response_code", errors.New(resp.Message), resp)
		return "", errors.New(helper.ErrorStatusCode)
	}

	return resp, nil
}
