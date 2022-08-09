package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// BarracudaService data structure
type BarracudaService struct {
	BaseURL   *url.URL
	BasicAuth string
}

// NewBarracudaService function for initializing barracuda service
func NewBarracudaService() (*BarracudaService, error) {
	var (
		barracuda BarracudaService
		err       error
		ok        bool
	)

	barracuda.BasicAuth, ok = os.LookupEnv("BARRACUDA_SERVICE_AUTH")
	if !ok {
		return &barracuda, errors.New("you need to specify BARRACUDA_SERVICE_AUTH in the environment variable")
	}

	baseURL, ok := os.LookupEnv("BARRACUDA_SERVICE_URL")
	if !ok {
		return &barracuda, errors.New("you need to specify BARRACUDA_SERVICE_URL in the environment variable")
	}

	barracuda.BaseURL, err = url.Parse(baseURL)
	if err != nil {
		return &barracuda, errors.New("error parsing barracuda services url")
	}

	return &barracuda, nil
}

// FindZipcode function for getting detail by zipcode
func (p *BarracudaService) FindZipcode(ctxReq context.Context, params serviceModel.ZipCodeQueryParameter) <-chan serviceModel.ServiceResult {
	ctx := "Service-FindZipcode"
	output := make(chan serviceModel.ServiceResult)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var response serviceModel.ResponseZipCode

		// generate headers
		headers := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "basic " + p.BasicAuth,
		}

		uri := fmt.Sprintf("%s/area/zipcode/?provinceId=%s&provinceName=%s&cityId=%s&cityName=%s&districtId=%s&districtName=%s&subDistrictId=%s&subDistrictName=%s&zipcode=%s",
			p.BaseURL.String(), params.ProvinceID, params.ProvinceName, params.CityID, params.CityName, params.DistrictID, params.DistrictName, params.SubDistrictID, params.SubDistrictName, params.ZipCode)
		uri = strings.Replace(uri, " ", "%20", -1)

		tags["uri"] = uri

		// request
		err := helper.GetHTTPNewRequest(ctxReq, http.MethodGet, uri, nil, &response, headers)
		if err != nil {
			tags["error"] = err
			e := errors.New("error getting data area zipcode")
			helper.SendErrorLog(ctxReq, ctx, "get_barracuda_area_zipcode", err, uri)
			output <- serviceModel.ServiceResult{Error: e}
			return
		}
		tags[helper.TextResponse] = response
		output <- serviceModel.ServiceResult{Result: response}
	})

	return output
}
