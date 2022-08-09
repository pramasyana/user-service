package shared

import (
	"net/http"
	"reflect"

	"github.com/Bhinneka/user-service/helper"
	"github.com/labstack/echo"
)

// HTTPResponse abstract interface
type HTTPResponse interface {
	SetSuccess(isSuccess bool)
	JSON(c echo.Context) error
}

type (
	// httpResponse model
	httpResponse struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Meta    interface{} `json:"meta,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Errors  interface{} `json:"errors,omitempty"`
	}

	// Meta model
	Meta struct {
		Page         int `json:"page"`
		Limit        int `json:"limit"`
		TotalRecords int `json:"totalRecords"`
		TotalPages   int `json:"totalPages"`
	}
)

// NewHTTPResponse for create common response, data must in first params and meta in second params
func NewHTTPResponse(code int, message string, params ...interface{}) HTTPResponse {
	commonResponse := new(httpResponse)

	for _, param := range params {
		// get value param if type is pointer
		refValue := reflect.ValueOf(param)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		param = refValue.Interface()

		switch param.(type) {
		case Meta:
			commonResponse.Meta = param
		case MultiError:
			multiError := param.(MultiError)
			commonResponse.Errors = multiError.ToMap()
		default:
			commonResponse.Data = param
		}
	}

	if code < http.StatusBadRequest && message != helper.ErrorDataNotFound {
		commonResponse.Success = true
	} else {
		commonResponse.Success = false
	}
	commonResponse.Code = code
	commonResponse.Message = message
	return commonResponse
}

// SetSuccess for set custom success
func (resp *httpResponse) SetSuccess(isSuccess bool) {
	resp.Success = isSuccess
}

// JSON for set http JSON response (Content-Type: application/json)
func (resp *httpResponse) JSON(c echo.Context) error {
	if resp.Data == nil {
		resp.Data = struct{}{}
	}
	return c.JSON(resp.Code, resp)
}
