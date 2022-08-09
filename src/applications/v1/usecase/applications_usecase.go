package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/applications/v1/model"
	"github.com/Bhinneka/user-service/src/applications/v1/repo"
)

const (
	msgErrorSaveApp = "failed to save application"
)

// ApplicationsUseCaseImpl data structure
type ApplicationsUseCaseImpl struct {
	ENVKey          string
	ApplicationRepo repo.ApplicationRepository
}

// NewApplicationsUseCase function for initialise applications use case implementation mo el
func NewApplicationsUseCase(envKey string, applicationRepo repo.ApplicationRepository) ApplicationsUseCase {
	return &ApplicationsUseCaseImpl{
		ENVKey:          envKey,
		ApplicationRepo: applicationRepo,
	}
}

// GetApplicationsList function for get list of session info
func (ap *ApplicationsUseCaseImpl) GetApplicationsList() <-chan ResultUseCase {
	ctx := "ApplicationUseCase-GetApplicationList"
	output := make(chan ResultUseCase)
	go func() {
		resp, err := http.Get(os.Getenv("APPLICATIONS_JSON"))
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextStmtError, err, resp)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respByte := buf.Bytes()

		var data model.ApplicationList
		errorUnmarshal := json.Unmarshal([]byte(respByte), &data)
		if errorUnmarshal != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextStmtError, err, string(respByte))
			output <- ResultUseCase{Error: errorUnmarshal, HTTPStatus: http.StatusBadRequest}
			return
		}
		output <- ResultUseCase{Result: data}
	}()

	return output
}

// AddUpdateApplication function for add new address
func (ap *ApplicationsUseCaseImpl) AddUpdateApplication(ctxReq context.Context, data model.Application) <-chan ResultUseCase {
	ctx := "ApplicationUseCase-AddUpdateApplication"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags["attr"] = data

		statusCode, err := ap.validateApplication(ctxReq, data)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: statusCode}
			return
		}

		data.Created = time.Now()
		data.LastModified = time.Now()

		var saveResult repo.ResultRepository
		if len(data.ID) > 0 {
			// update application repository process to database
			saveResult = <-ap.ApplicationRepo.Update(ctxReq, data)
		} else {
			// add application repository process to database
			saveResult = <-ap.ApplicationRepo.Save(ctxReq, data)
		}

		if saveResult.Error != nil {
			err := errors.New(msgErrorSaveApp)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		resultData, ok := saveResult.Result.(model.Application)
		if !ok {
			err := errors.New(msgErrorSaveApp)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		tags[helper.TextResponse] = resultData
		output <- ResultUseCase{Result: resultData}
	})

	return output
}

// validateApplication function for validate new address
func (ap *ApplicationsUseCaseImpl) validateApplication(ctxReq context.Context, data model.Application) (int, error) {
	if len(data.ID) > 0 {
		// find app by id
		findResult := <-ap.ApplicationRepo.FindApplicationByID(ctxReq, data.ID)
		if findResult.Error != nil {
			err := errors.New("cannot find application")
			return http.StatusNotFound, err
		}
	}

	_, err := url.ParseRequestURI(data.URL)
	if err != nil {
		err := errors.New("Url not valid")
		return http.StatusBadRequest, err
	}

	if !helper.ValidateDocumentFileURL(data.Logo) {
		err := errors.New("Logo not valid")
		return http.StatusBadRequest, err
	}

	return 0, nil
}

// DeleteApplication function for delete application
func (ap *ApplicationsUseCaseImpl) DeleteApplication(ctxReq context.Context, id string) <-chan ResultUseCase {
	ctx := "ApplicationUseCase-DeleteApplication"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags["id"] = id

		// find app by id
		findResult := <-ap.ApplicationRepo.FindApplicationByID(ctxReq, id)
		if findResult.Error != nil {
			err := errors.New("cannot find application")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		// delete application from database
		result := <-ap.ApplicationRepo.Delete(ctxReq, id)
		if result.Error != nil {
			err := errors.New("failed to delete application")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: nil}
	})
	return output
}

// GetListApplication function for getting list of shipping address
func (ap *ApplicationsUseCaseImpl) GetListApplication(ctxReq context.Context, params *model.ParametersApplication) <-chan ResultUseCase {
	ctx := "ApplicationUseCase-GetListApplication"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var err error

		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1, // default
				StrPage:  params.StrPage,
				Limit:    20, // default
				StrLimit: params.StrLimit,
			})

		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		params.Page = paging.Page
		params.Limit = paging.Limit
		params.Offset = paging.Offset
		tags[helper.TextParameter] = params

		applicationResult := <-ap.ApplicationRepo.GetListApplication(ctxReq, params)
		if applicationResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			// when data is not found
			if applicationResult.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				applicationResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "application")
			}

			output <- ResultUseCase{Error: applicationResult.Error, HTTPStatus: httpStatus}
			return
		}

		application := applicationResult.Result.(model.ListApplication)

		totalResult := <-ap.ApplicationRepo.GetTotalApplication(ctxReq, params)
		if totalResult.Error != nil {
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		total := totalResult.Result.(int)
		application.TotalData = total

		output <- ResultUseCase{Result: application}
	})

	return output
}
