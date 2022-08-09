package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Bhinneka/bhinneka-go-sdk"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	dolpinURL          = "https://dolphin.bhinnekatesting.com"
	dolphinRegisterURL = "https://dolphin.bhinnekatesting.com/index.php?entryPoint=syncAccountEntryPoint"
	dolphinGetURL      = "https://dolphin.bhinnekatesting.com/index.php?entryPoint=syncAccountEntryPoint&id=%s"
)

func TestInitDolphin(t *testing.T) {
	_, err := NewDolphinService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "specify DOLPHIN_BASIC_AUTH")

	os.Setenv("DOLPHIN_BASIC_AUTH", defaultAuth)
	_, err = NewDolphinService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "specify DOLPHIN_BASE_URL")

	os.Setenv("DOLPHIN_BASE_URL", badURL)
	b, errm := NewDolphinService()
	assert.Nil(t, b.BaseURL)
	assert.Error(t, errm)

	os.Setenv("DOLPHIN_BASE_URL", dolpinURL)
	b, errm = NewDolphinService()
	assert.NotNil(t, b.BaseURL)
	assert.NoError(t, errm)
}

var (
	memberDolphin   = serviceModel.MemberDolphin{ID: "1", Email: "pian.mutakin@bhinneka.com"}
	defaultMemberID = "179"
)

var testDataDolphinRegister = []struct {
	name            string
	param           serviceModel.MemberDolphin
	serviceResponse interface{}
	wantError       bool
	statusCode      int
}{
	{
		name:            "Test Register #1",
		param:           memberDolphin,
		wantError:       false,
		statusCode:      http.StatusOK,
		serviceResponse: serviceModel.Response{Data: serviceModel.Data{Attributes: serviceModel.Attributes{IsSuccess: true}}},
	},
	{
		name:       "Test Register #2",
		param:      memberDolphin,
		wantError:  true,
		statusCode: http.StatusBadRequest,
	},
	{
		name:            "Test Register #3",
		param:           memberDolphin,
		wantError:       true,
		statusCode:      http.StatusBadRequest,
		serviceResponse: []byte(`ss`),
	},
	{
		name:            "Test Register #4",
		param:           memberDolphin,
		wantError:       true,
		statusCode:      http.StatusBadRequest,
		serviceResponse: serviceModel.Response{Data: serviceModel.Data{Attributes: serviceModel.Attributes{IsSuccess: false, Message: "some error message"}}},
	},
}

var testDataGetMember = []struct {
	name            string
	memberID        string
	serviceResponse interface{}
	wantError       bool
	statusCode      int
}{
	{
		name:            "Get Member #1",
		memberID:        defaultMemberID,
		wantError:       false,
		statusCode:      http.StatusOK,
		serviceResponse: &serviceModel.MemberResponse{},
	},
	{
		name:            "Get Member #2",
		memberID:        defaultMemberID,
		wantError:       true,
		statusCode:      http.StatusOK,
		serviceResponse: []byte(`some value `),
	},
	{
		name:            "Get Member #3",
		memberID:        defaultMemberID,
		wantError:       true,
		statusCode:      http.StatusOK,
		serviceResponse: &serviceModel.MemberResponse{Data: serviceModel.DataMember{Attributes: serviceModel.MemberDolphin{Message: DolphinDataNotFound}}},
	},
}

func TestDolphin(t *testing.T) {
	os.Setenv("DOLPHIN_BASIC_AUTH", defaultAuth)
	os.Setenv("DOLPHIN_BASE_URL", dolpinURL)
	ds, _ := NewDolphinService()
	ctx := context.Background()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("RegisterMember", func(t *testing.T) {
		for _, tc := range testDataDolphinRegister {
			bhinneka.MockHTTP(http.MethodPost, dolphinRegisterURL, tc.statusCode, tc.serviceResponse)
			err := ds.RegisterMember(ctx, tc.param)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		}
	})

	t.Run("UpdateMember", func(t *testing.T) {
		for _, tc := range testDataDolphinRegister {
			bhinneka.MockHTTP(http.MethodPut, dolphinRegisterURL, tc.statusCode, tc.serviceResponse)
			err := ds.UpdateMember(ctx, tc.param)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		}
	})
	t.Run("ActivateMember", func(t *testing.T) {
		for _, tc := range testDataDolphinRegister {
			bhinneka.MockHTTP(http.MethodPut, dolphinRegisterURL, tc.statusCode, tc.serviceResponse)
			err := ds.ActivateMember(ctx, tc.param)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		}
	})

	t.Run("GetMember", func(t *testing.T) {
		for _, tc := range testDataGetMember {
			bhinneka.MockHTTP(http.MethodGet, fmt.Sprintf(dolphinGetURL, tc.memberID), tc.statusCode, tc.serviceResponse)
			sr, err := ds.GetMember(ctx, tc.memberID)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sr)
			}
		}
	})
}
