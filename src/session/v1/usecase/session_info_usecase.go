package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bhinneka/golib/jsonschema"
	goString "github.com/Bhinneka/golib/string"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/session/v1/model"
	"github.com/Bhinneka/user-service/src/session/v1/query"
	"github.com/Bhinneka/user-service/src/session/v1/repo"
)

// SessionInfoUseCaseImpl data structure
type SessionInfoUseCaseImpl struct {
	SessionInfoQuery query.SessionInfoQuery
	SessionInfoRepo  repo.SessionInfoRepository
}

// NewSessionInfoUseCase function for initialise session info use case implementation mo el
func NewSessionInfoUseCase(sessionInfoQuery query.SessionInfoQuery, sessionInfoRepository repo.SessionInfoRepository) SessionInfoUseCase {
	return &SessionInfoUseCaseImpl{
		SessionInfoQuery: sessionInfoQuery,
		SessionInfoRepo:  sessionInfoRepository,
	}
}

// GetSessionInfoList function for get list of session info
func (gs *SessionInfoUseCaseImpl) GetSessionInfoList(params *model.ParamList) <-chan ResultUseCase {
	ctx := "SessionInfoUseCase-GetSessionInfoList"

	output := make(chan ResultUseCase)
	go func() {
		ctxReq := context.Background()
		// validate request
		mErr := jsonschema.ValidateTemp("session_get_list_params", params)
		if mErr != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, mErr, params)
			output <- ResultUseCase{Error: mErr, HTTPStatus: http.StatusBadRequest}
			return
		}

		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("%v", r)
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, nil)
			}
			close(output)
		}()
		var err error

		params, err := gs.validateParams(params)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		queryResult := <-gs.SessionInfoQuery.GetListSessionInfo(ctxReq, params)

		if queryResult.Error != nil {
			httpStatus := http.StatusInternalServerError
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, queryResult.Error, nil)

			// when data is not found
			if queryResult.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				queryResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "session_info")
			}

			output <- ResultUseCase{Error: queryResult.Error, HTTPStatus: httpStatus}
			return

		}

		sessionInfoList := queryResult.Result.(model.SessionInfoList)

		totalResult := <-gs.SessionInfoQuery.GetTotalSessionInfo(params)
		if totalResult.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_total", err, totalResult)
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		totalRecord := totalResult.Result.(int)
		sessionInfoList.TotalData = totalRecord

		output <- ResultUseCase{Result: sessionInfoList}

	}()
	return output
}

// validateParams function for validate params
func (gs *SessionInfoUseCaseImpl) validateParams(params *model.ParamList) (*model.ParamList, error) {
	// validate all parameters
	paging, err := helper.ValidatePagination(
		helper.PaginationParameters{
			Page:     1, // default
			StrPage:  params.StrPage,
			Limit:    10, // default
			StrLimit: params.StrLimit,
		})

	if err != nil {
		return params, err
	}
	params.Page = paging.Page
	params.Limit = paging.Limit
	params.Offset = paging.Offset

	if len(params.OrderBy) > 0 {
		if !helper.StringInSlice(params.OrderBy, model.AllowedSortFields) {
			return params, fmt.Errorf(helper.ErrorParameterInvalid, "order by")
		}
	} else {
		params.OrderBy = "id"
	}

	if len(params.Sort) > 0 {
		if !helper.StringInSlice(params.Sort, []string{"asc", "desc"}) {
			return params, fmt.Errorf(helper.ErrorParameterInvalid, "sort")
		}
	} else {
		params.Sort = "desc"
	}
	if params.Range != "" {
		params.RangeInt, _ = strconv.Atoi(params.Range)
	}
	if goString.IsValidEmail(params.Query) {
		params.Email = params.Query
		params.Query = ""
	}

	if strings.HasPrefix(params.Query, "USR") {
		params.MemberID = params.Query
		params.Query = ""
	}

	return params, nil
}
