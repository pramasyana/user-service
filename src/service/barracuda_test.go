package service

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/Bhinneka/bhinneka-go-sdk"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	defaultAuth   = "bearer token"
	defaultURL    = "https://barracuda.bhinnekatesting.com"
	badURL        = "https://segment%%2815197306101420000%29.ts"
	getZipCodeURL = "https://barracuda.bhinnekatesting.com/area/zipcode/?provinceId=11\u0026provinceName=\u0026cityId=\u0026cityName=\u0026districtId=\u0026districtName=\u0026subDistrictId=\u0026subDistrictName=\u0026zipcode="
)

func TestInitBarracuda(t *testing.T) {
	_, err := NewBarracudaService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "specify BARRACUDA_SERVICE_AUTH")

	os.Setenv("BARRACUDA_SERVICE_AUTH", defaultAuth)
	_, err = NewBarracudaService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "specify BARRACUDA_SERVICE_URL")

	os.Setenv("BARRACUDA_SERVICE_URL", badURL)
	b, errm := NewBarracudaService()
	assert.Nil(t, b.BaseURL)
	assert.Error(t, errm)

	os.Setenv("BARRACUDA_SERVICE_URL", defaultURL)
	b, errm = NewBarracudaService()
	assert.NotNil(t, b.BaseURL)
	assert.NoError(t, errm)
}

func TestBarracuda(t *testing.T) {
	os.Setenv("BARRACUDA_SERVICE_AUTH", defaultAuth)
	os.Setenv("BARRACUDA_SERVICE_URL", defaultURL)
	b, _ := NewBarracudaService()
	ctx := context.Background()
	var testData = []struct {
		name            string
		serviceResponse interface{}
		wantError       bool
		statusCode      int
	}{
		{
			name:            "Find Zip Code #1",
			serviceResponse: serviceModel.ResponseZipCode{Status: true},
			wantError:       false,
			statusCode:      http.StatusOK,
		},
		{
			name:            "Find Zip Code #2",
			serviceResponse: `{"code":200}`,
			wantError:       true,
			statusCode:      http.StatusBadRequest,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tc := range testData {
		bhinneka.MockHTTP(http.MethodGet, getZipCodeURL, tc.statusCode, tc.serviceResponse)
		sr := <-b.FindZipcode(ctx, serviceModel.ZipCodeQueryParameter{ProvinceID: "11"})
		if tc.wantError {
			assert.Error(t, sr.Error)
		} else {
			assert.NoError(t, sr.Error)
		}
	}

}
