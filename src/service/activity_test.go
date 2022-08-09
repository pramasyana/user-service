package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Bhinneka/bhinneka-go-sdk"
	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	actURL        = `http://staging.bhinnekalocal.com/activity-service/v2/logs?dateFrom=2020-10-21&dateTo=2020-10-22&limit=1&orderBy=created&page=1&phrase%5Baction%5D=UPDATE&phrase%5BcreatorId%5D=myUserID&phrase%5Bmodule%5D=LKPP&phrase%5BobjectType%5D=users&phrase%5Bpack%5D=Sturgeon&phrase%5Btarget%5D=sometarget&sort=asc&viewType=public`
	getLogURLByID = `https://activity.bhinneka.com/v1/logs/%s`
	createLogURL  = `http://dev.bhinnekalocal.com/activity-service/v2/logs`
	bearerAuth    = `Bearer someToken`
)

var testDataGetLogs = []struct {
	name            string
	wantError       bool
	serviceResponse interface{}
	statusCode      int
	param           *sharedModel.Parameters
}{
	{
		name:       "get logs #1",
		wantError:  false,
		statusCode: http.StatusOK,
		param: &sharedModel.Parameters{
			Page:       1,
			Module:     "LKPP",
			Action:     "UPDATE",
			Pack:       "Sturgeon",
			Limit:      1,
			Sort:       "asc",
			OrderBy:    "created",
			DateFrom:   "2020-10-21",
			DateTo:     "2020-10-22",
			ViewType:   "public",
			Target:     "sometarget",
			Creator:    "myUserID",
			ObjectType: "users",
		},
	},
	{
		name:       "get logs #2",
		wantError:  true,
		statusCode: http.StatusBadRequest,
		param:      nil,
	},
}

var testDataGetLogByID = []struct {
	name            string
	wantError       bool
	serviceResponse interface{}
	statusCode      int
	param           *sharedModel.Parameters
}{
	{
		name:       "get single log #1",
		wantError:  false,
		statusCode: http.StatusOK,
		param:      &sharedModel.Parameters{ID: "someID"},
	},
	{
		name:       "get single log #2",
		wantError:  true,
		statusCode: http.StatusBadRequest,
		param:      &sharedModel.Parameters{ID: "otherID"},
	},
}

func TestGetLogs(t *testing.T) {
	os.Setenv("ENV", "STAGING")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tc := range testDataGetLogs {
		bhinneka.MockHTTP(http.MethodGet, actURL, tc.statusCode, tc.serviceResponse)
		a := NewActivityService("v2")
		ctx := context.WithValue(context.Background(), helper.TextAuthorization, bearerAuth)
		sr := <-a.GetAll(ctx, tc.param)
		if tc.wantError {
			assert.Error(t, sr.Error)
		} else {
			assert.NoError(t, sr.Error)
		}
	}
}

func TestGetLogByID(t *testing.T) {
	os.Setenv("ENV", "PROD")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tc := range testDataGetLogByID {
		fetchURL := fmt.Sprintf(getLogURLByID, tc.param.ID)
		bhinneka.MockHTTP(http.MethodGet, fetchURL, tc.statusCode, tc.serviceResponse)
		a := NewActivityService("")
		ctx := context.WithValue(context.Background(), helper.TextAuthorization, bearerAuth)
		sr := <-a.GetLogByID(ctx, tc.param.ID)
		if tc.wantError {
			assert.Error(t, sr.Error)
		} else {
			assert.NoError(t, sr.Error)
		}
	}
}

var testDataCreateLog = []struct {
	name            string
	wantError       bool
	serviceResponse interface{}
	statusCode      int
	param           serviceModel.Payload
}{
	{
		name:       "create log #1",
		wantError:  false,
		statusCode: http.StatusOK,
		param:      serviceModel.Payload{ID: "createLogByID"},
	},
	{
		name:       "create log #2",
		wantError:  true,
		statusCode: http.StatusBadRequest,
		param:      serviceModel.Payload{ID: "myOtherID"},
	},
}

func TestCreateLog(t *testing.T) {
	os.Setenv("ENV", "DEV")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tc := range testDataCreateLog {
		bhinneka.MockHTTP(http.MethodPost, createLogURL, tc.statusCode, tc.serviceResponse)
		a := NewActivityService("v2")
		ctx := context.WithValue(context.Background(), helper.TextAuthorization, bearerAuth)
		ctx = context.WithValue(ctx, clientIP, "127.0.0.1")
		sr := <-a.CreateLog(ctx, tc.param)
		if tc.wantError {
			assert.Error(t, sr.Error)
		} else {
			assert.NoError(t, sr.Error)
		}
	}
}

func TestInsertLog(t *testing.T) {
	old := sharedModel.DBOperator{Field: "email", Operator: "="}
	new := sharedModel.DBOperator{Field: "email", Operator: ">="}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	bhinneka.MockHTTP(http.MethodPost, createLogURL, http.StatusOK, nil)

	a := NewActivityService("v2")
	ctx := context.WithValue(context.Background(), helper.TextAuthorization, bearerAuth)
	ctx = context.WithValue(ctx, clientIP, "127.0.0.1")
	err := a.InsertLog(ctx, old, new, serviceModel.Payload{Module: "new module"})
	assert.NoError(t, err)
}
