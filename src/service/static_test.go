package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Bhinneka/bhinneka-go-sdk"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	constStaticAuth = "STATIC_SERVICE_AUTH"
	constStaticURL  = "STATIC_SERVICE_URL"
	constStaticGWS  = "GWS_GRAPHQL_URL"
	staticURL       = "http://static.bhinnekatesting.com"
)

var (
	defaultGwsURL   = fmt.Sprintf("%s/graphql", staticURL)
	defaultStaticID = fmt.Sprintf("%s/v1/statics/1", staticURL)
)

func TestInitStatic(t *testing.T) {
	_, err := NewStaticService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "specify STATIC_SERVICE_AUTH")

	os.Setenv(constStaticAuth, defaultAuth)
	_, err = NewStaticService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "specify STATIC_SERVICE_URL")

	os.Setenv(constStaticURL, badURL)
	_, errm := NewStaticService()
	assert.Error(t, errm)

	os.Setenv(constStaticURL, staticURL)
	_, errm = NewStaticService()
	assert.Error(t, errm)

	os.Setenv(constStaticGWS, badURL)
	_, errm = NewStaticService()
	assert.Error(t, errm)

	os.Setenv(constStaticGWS, defaultGwsURL)
	_, errm = NewStaticService()
	assert.NoError(t, errm)
}

func TestFindStatic(t *testing.T) {
	os.Setenv(constStaticAuth, defaultAuth)
	os.Setenv(constStaticURL, staticURL)
	os.Setenv(constStaticGWS, defaultGwsURL)

	ctx := context.Background()
	var testDatas = []struct {
		name            string
		wantError       bool
		serviceResponse interface{}
		statusCode      int
	}{
		{
			name:            "Test Find Static By ID #1",
			wantError:       false,
			serviceResponse: nil,
			statusCode:      http.StatusOK,
		},
		{
			name:            "Test Find Static By ID #2",
			wantError:       true,
			serviceResponse: []byte(``),
			statusCode:      http.StatusBadRequest,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tc := range testDatas {
		ss, _ := NewStaticService()
		bhinneka.MockHTTP(http.MethodGet, defaultStaticID, tc.statusCode, tc.serviceResponse)
		sr := <-ss.FindStaticsByID(ctx, "1")
		if tc.wantError {
			assert.Error(t, sr.Error)
		} else {
			assert.NoError(t, sr.Error)
		}
	}
}

func TestFindStaticGws(t *testing.T) {
	os.Setenv(constStaticAuth, defaultAuth)
	os.Setenv(constStaticURL, staticURL)
	os.Setenv(constStaticGWS, staticURL)

	ctx := context.Background()
	var testDatas = []struct {
		name            string
		wantError       bool
		serviceResponse interface{}
		statusCode      int
	}{
		{
			name:            "Test Find Static GWS By ID #1",
			wantError:       false,
			serviceResponse: nil,
			statusCode:      http.StatusOK,
		},
		{
			name:            "Test Find Static GWS By ID #2",
			wantError:       true,
			serviceResponse: []byte(``),
			statusCode:      http.StatusBadRequest,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tc := range testDatas {
		ss, _ := NewStaticService()
		bhinneka.MockHTTP(http.MethodPost, defaultGwsURL, tc.statusCode, tc.serviceResponse)
		sr := <-ss.FindStaticsGwsByID(ctx, "1")
		if tc.wantError {
			assert.Error(t, sr.Error)
		} else {
			assert.NoError(t, sr.Error)
		}
	}
}
