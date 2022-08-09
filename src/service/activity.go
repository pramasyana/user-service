package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Bhinneka/bhinneka-go-sdk"
	"github.com/Bhinneka/bhinneka-go-sdk/activity-service"
	"github.com/Bhinneka/bhinneka-go-sdk/general"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/spf13/cast"
)

//ActivityServiceImpl constructor
type ActivityServiceImpl struct {
	act            *activity.BasicActivity
	TokenGenerator token.AccessTokenGenerator
}

const (
	logPath       = "logs"
	clientIP      = string("ClientIP")
	defaultTmeout = 10 * time.Second
)

// NewActivityService function for initializing activity service
func NewActivityService(version string) *ActivityServiceImpl {
	var env string

	switch os.Getenv("ENV") {
	case "STAGING":
		env = activity.Staging
	case "PROD":
		env = activity.Production
	default:
		env = activity.Devel
	}
	if version == "" {
		version = "v1"
	}

	param := activity.Request{
		BasicToken: "Basic " + os.Getenv("STATIC_SERVICE_AUTH"),
		Version:    version,
		Env:        env,
		Protocol:   bhinneka.REST,
		Timeout:    defaultTmeout,
	}

	rev, _ := activity.NewActivityRequest(param)
	return &ActivityServiceImpl{
		act: rev,
	}
}

func (as *ActivityServiceImpl) setAuthArgument(ctx context.Context) *general.Args {
	args := general.NewHTTPArgument()
	args.SetContentType(bhinneka.ContentJSON)
	authVal := ctx.Value(helper.TextAuthorization)
	if authVal != nil {
		if authorization, ok := authVal.(string); ok {
			args.SetAuth(authorization)
		}
	}

	return args
}

// CreateLog function for write log to activity services
func (as *ActivityServiceImpl) CreateLog(ctxReq context.Context, param serviceModel.Payload) <-chan serviceModel.ServiceResult {
	ctx := "ActivityService-CreateLog"

	var (
		err       error
		creatorIP = "127.0.0.1"
	)
	output := make(chan serviceModel.ServiceResult)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		ctxReq, cancel := context.WithTimeout(ctxReq, defaultTmeout)
		defer cancel()

		defer close(output)
		param.Pack = "sturgeon"
		if ipAddr := ctxReq.Value(middleware.ContextKeyClientIP); ipAddr != nil {
			if ipAddress, ok := ipAddr.(string); ok {
				creatorIP = ipAddress
			}
		}
		param.CreatorIP = creatorIP
		args := as.setAuthArgument(ctxReq)
		args.SetMethod(http.MethodPost)
		args.SetURI(logPath)

		args.SetParam(param)
		var respExec activity.ResponseExecV2
		tags[helper.TextHeader] = args
		tags[helper.TextParameter] = param
		if err = as.act.Exec(args, &respExec); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "http_request_activity", err, args)
			output <- serviceModel.ServiceResult{Error: err}
			return
		}
		tags[helper.TextResponse] = respExec
		output <- serviceModel.ServiceResult{Result: respExec}
	})
	return output
}

// InsertLog function for insert log
func (as *ActivityServiceImpl) InsertLog(ctxReq context.Context, oldData, newData interface{}, payload serviceModel.Payload) error {
	oldDataJSON, _ := json.Marshal(oldData)
	newDataJSON, _ := json.Marshal(newData)
	// Declared an empty interface
	var resultNew map[string]interface{}
	var resultOld map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(newDataJSON), &resultNew)
	json.Unmarshal([]byte(oldDataJSON), &resultOld)

	logs := []serviceModel.Log{}
	for key, data := range resultNew {
		if cast.ToString(data) == cast.ToString(resultOld[key]) {
			continue
		}

		logData := serviceModel.Log{
			Field:    key,
			NewValue: cast.ToString(data),
			OldValue: cast.ToString(resultOld[key]),
		}
		logs = append(logs, logData)
	}

	if len(logs) > 0 {
		payload.Logs = logs
		<-as.CreateLog(ctxReq, payload)
	}

	return nil
}

// GetAll return all log from activity service
func (as *ActivityServiceImpl) GetAll(ctxReq context.Context, param *sharedModel.Parameters) <-chan serviceModel.ServiceResult {
	ctx := "ActivityService-GetAllLogs"
	sr := make(chan serviceModel.ServiceResult)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		ctxReq, cancel := context.WithTimeout(ctxReq, defaultTmeout)
		defer cancel()

		defer close(sr)
		args := as.setAuthArgument(ctxReq)
		args.SetMethod(http.MethodGet)

		uri := fmt.Sprintf("%s?%s", logPath, as.parseParameter(param))
		args.SetURI(uri)

		var respExec serviceModel.ResponseService
		if err := as.act.Exec(args, &respExec); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "http_request_get_all", err, args)
			sr <- serviceModel.ServiceResult{Error: err}
			return
		}
		tags[helper.TextResponse] = respExec
		sr <- serviceModel.ServiceResult{Result: respExec.Data, Meta: respExec.Meta}
	})

	return sr
}

func (as *ActivityServiceImpl) parseParameter(p *sharedModel.Parameters) string {
	if p == nil {
		return ""
	}
	qParam := url.Values{}
	var (
		sort    = "desc"
		orderBy = "created"
		page    = 1
		limit   = 10
	)
	if p.Module != "" {
		qParam.Add("phrase[module]", p.Module)
	}
	if p.Action != "" {
		qParam.Add("phrase[action]", p.Action)
	}
	if p.Pack != "" {
		qParam.Add("phrase[pack]", p.Pack)
	}
	if p.Creator != "" {
		qParam.Add("phrase[creatorId]", p.Creator)
	}
	if p.Target != "" {
		qParam.Add("phrase[target]", p.Target)
	}
	if p.ObjectType != "" {
		qParam.Add("phrase[objectType]", p.ObjectType)
	}
	if p.Page > 0 {
		page = p.Page
	}

	if p.Limit > 0 {
		limit = p.Limit
	}

	if p.Sort != "" {
		sort = p.Sort
	}

	if p.OrderBy != "" {
		orderBy = p.OrderBy
	}

	if p.DateFrom != "" {
		qParam.Add("dateFrom", p.DateFrom)
	}
	if p.DateTo != "" {
		qParam.Add("dateTo", p.DateTo)
	}
	if p.ViewType != "" {
		qParam.Add("viewType", p.ViewType)
	}

	qParam.Add("orderBy", orderBy)
	qParam.Add("sort", sort)
	qParam.Add("limit", strconv.Itoa(limit))
	qParam.Add("page", strconv.Itoa(page))
	return qParam.Encode()
}

// GetLogByID return single log from activity service
func (as *ActivityServiceImpl) GetLogByID(ctxReq context.Context, logID string) <-chan serviceModel.ServiceResult {
	ctx := "ActivityService-GetLogByID"
	sr := make(chan serviceModel.ServiceResult)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		ctxReq, cancel := context.WithTimeout(ctxReq, defaultTmeout)
		defer cancel()

		defer close(sr)
		args := as.setAuthArgument(ctxReq)
		args.SetMethod(http.MethodGet)

		uri := fmt.Sprintf("%s/%s", logPath, logID)
		args.SetURI(uri)

		var respExec serviceModel.ResponseService
		if err := as.act.Exec(args, &respExec); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "http_request_get_by_id", err, args)
			sr <- serviceModel.ServiceResult{Error: err}
			return
		}
		tags[helper.TextResponse] = respExec
		sr <- serviceModel.ServiceResult{Result: respExec.Data}
	})

	return sr
}
