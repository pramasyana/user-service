package query

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/session/v1/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

const (
	userIDparams = `"userId" = '%s'`
	whereParams  = " WHERE %s"
	errDefault   = "failed executing request"
)

// SessionInfoQueryPostgres data struct
type SessionInfoQueryPostgres struct {
	repo *repository.Repository
}

// NewSessionInfoQueryPostgres function for initializing auth query
func NewSessionInfoQueryPostgres(repo *repository.Repository) *SessionInfoQueryPostgres {
	return &SessionInfoQueryPostgres{repo: repo}
}

// GetListSessionInfo for save session info user login
func (qp *SessionInfoQueryPostgres) GetListSessionInfo(ctxReq context.Context, params *model.ParamList) <-chan ResultQuery {
	ctx := "SessionInfoQuery-GetListSessionInfo"
	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		strQuery, queryValues := qp.generateFilter(params)

		if len(params.OrderBy) > 0 {
			params.OrderBy = fmt.Sprintf(`"%s"`, params.OrderBy)
		}
		Limit := fmt.Sprintf(`LIMIT %d OFFSET %d`, params.Limit, params.Offset)
		if params.Limit == 999 {
			Limit = ``
		}

		q := fmt.Sprintf(`SELECT id, "userId", "userName", ip, "userAgent", "deviceId", "clientType", "grantType", jti, "createdAt" from session_info
						%s
						ORDER BY %s %s
						%s`, strQuery, params.OrderBy, params.Sort, Limit)

		tags[helper.TextQuery] = q
		tags["parameters"] = queryValues
		rows, err := qp.repo.Query(q, queryValues...)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, q)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: nil}
			return
		}
		defer rows.Close()

		var datalist model.SessionInfoList

		for rows.Next() {
			var SessionInfoResponse model.SessionInfoResponse

			err = rows.Scan(
				&SessionInfoResponse.ID, &SessionInfoResponse.UserID, &SessionInfoResponse.UserName, &SessionInfoResponse.IP,
				&SessionInfoResponse.UserAgent, &SessionInfoResponse.DeviceID, &SessionInfoResponse.ClientType,
				&SessionInfoResponse.GrantType, &SessionInfoResponse.JTI, &SessionInfoResponse.CreatedAt,
			)

			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, q)
				tags[helper.TextResponse] = err
				output <- ResultQuery{Error: err}
				return
			}

			datalist.Data = append(datalist.Data, SessionInfoResponse)
		}

		tags[helper.TextResponse] = datalist
		output <- ResultQuery{Error: nil, Result: datalist}
	})

	return output
}

// GetTotalSessionInfo function for getting total of members
func (qp *SessionInfoQueryPostgres) GetTotalSessionInfo(params *model.ParamList) <-chan ResultQuery {
	ctx := "SessionInfoQuery-GetTotalSessionInfo"

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		var totalData int
		ctxReq := context.Background()

		strQuery, queryValues := qp.generateFilter(params)

		if len(params.OrderBy) > 0 {
			params.OrderBy = fmt.Sprintf(`"%s"`, params.OrderBy)
		}

		sq := fmt.Sprintf(`SELECT count(id) FROM session_info %s`, strQuery)

		stmt, err := qp.repo.Prepare(ctxReq, sq)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, sq)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		err = stmt.QueryRow(queryValues...).Scan(&totalData)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, sq)
			output <- ResultQuery{Error: err}
			return
		}

		output <- ResultQuery{Result: totalData}
	}()

	return output
}

// generateFilter function for generating filter query
func (qp *SessionInfoQueryPostgres) generateFilter(params *model.ParamList) (string, []interface{}) {
	var (
		strQuery    string
		queryStrOR  []string
		queryStrAND []string
		queryValues []interface{}
	)

	if len(params.Query) > 0 {
		queryStrOR = append(queryStrOR, `"userAgent" ilike $`+strconv.Itoa(len(queryStrOR)+1))
		queryValues = append(queryValues, "%"+params.Query+"%")
	}

	if len(params.MemberID) > 0 {
		queryStrAND = append(queryStrAND, `"userId" = $`+strconv.Itoa(len(queryStrAND)+1))
		queryValues = append(queryValues, params.MemberID)
	}
	if len(params.Email) > 0 {
		queryStrAND = append(queryStrAND, `"userName" = $`+strconv.Itoa(len(queryStrAND)+1))
		queryValues = append(queryValues, params.Email)
	}
	if len(params.ClientType) > 0 {
		queryStrAND = append(queryStrAND, `"clientType" = $`+strconv.Itoa(len(queryStrAND)+1))
		queryValues = append(queryValues, params.ClientType)
	}
	if params.RangeInt > 0 {
		queryStrAND = append(queryStrAND, fmt.Sprintf(`"createdAt" > (CURRENT_DATE - INTERVAL '%d DAY')::DATE`, params.RangeInt))
	}

	if len(queryStrOR) > 0 || len(queryStrAND) > 0 {
		if len(queryStrOR) > 0 {
			strQuery = fmt.Sprintf(`(%s)`, strings.Join(queryStrOR, " OR "))
			queryStrAND = append(queryStrAND, strQuery)
		}

		if len(queryStrAND) > 0 {
			strQuery = strings.Join(queryStrAND, " AND ")
		}
	}

	if len(strQuery) > 0 {
		strQuery = fmt.Sprintf(whereParams, strQuery)
	}

	return strQuery, queryValues
}

func (qp *SessionInfoQueryPostgres) mapParams(param model.ParametersGetSession) ([]string, []interface{}) {
	var (
		output     []string
		bindVar    []interface{}
		bindNumber int
	)

	if param.DeviceID != "" {
		bindNumber++
		output = append(output, fmt.Sprintf(`"deviceId" = $%s`, strconv.Itoa(bindNumber)))
		bindVar = append(bindVar, param.DeviceID)
	}

	if param.ClientType != "" {
		bindNumber++
		output = append(output, fmt.Sprintf(`"clientType" = $%s`, strconv.Itoa(bindNumber)))
		bindVar = append(bindVar, param.ClientType)
	}

	if param.UserID != "" {
		bindNumber++
		output = append(output, fmt.Sprintf(`"userId" = $%s`, strconv.Itoa(bindNumber)))
		bindVar = append(bindVar, param.UserID)
	}

	if param.Jti != "" {
		bindNumber++
		output = append(output, fmt.Sprintf(`"jti" = $%s`, strconv.Itoa(bindNumber)))
		bindVar = append(bindVar, param.Jti)
	}

	if param.SessionID != "" {
		bindNumber++
		output = append(output, fmt.Sprintf(`"id" = $%s`, strconv.Itoa(bindNumber)))
		bindVar = append(bindVar, param.SessionID)
	}
	return output, bindVar
}

// GetDetailSessionInfo function for getting detail sessuib ubfi
func (qp *SessionInfoQueryPostgres) GetDetailSessionInfo(ctxReq context.Context, param model.ParametersGetSession) <-chan ResultQuery {
	ctx := "SessionInfoQueryPostgres-GetDetailSessionInfo"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			queryParam string
		)
		queryList, bindVar := qp.mapParams(param)

		if len(queryList) > 0 {
			queryParam = fmt.Sprintf(whereParams, strings.Join(queryList, " AND "))
		}

		q := fmt.Sprintf(`SELECT id, "userId", "userName", ip, "userAgent",
						"deviceId", "clientType", "grantType", jti, "createdAt" 
						from session_info %s ORDER BY "createdAt" DESC LIMIT 1`, queryParam)

		tags[helper.TextQuery] = q

		stmt, err := qp.repo.Prepare(ctxReq, q)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: errors.New(errDefault)}
			return
		}
		defer stmt.Close()

		var session model.SessionInfoResponse

		err = stmt.QueryRow(bindVar...).Scan(
			&session.ID,
			&session.UserID,
			&session.UserName,
			&session.IP,
			&session.UserAgent,
			&session.DeviceID,
			&session.ClientType,
			&session.GrantType,
			&session.JTI,
			&session.CreatedAt,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}

		output <- ResultQuery{Result: session}
	})
	return output
}

// GetHistorySessionInfo function for getting history session
func (qp *SessionInfoQueryPostgres) GetHistorySessionInfo(ctxReq context.Context, params *memberModel.ParametersLoginActivity) <-chan ResultQuery {
	ctx := "SessionInfoQueryPostgres-GetHistorySessionInfo"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			queryParam string
			queryList  []string
		)

		queryList = append(queryList, `"createdAt" > (CURRENT_DATE - INTERVAL '30 DAY')::DATE`)
		if params.ExcludeID != "" {
			queryList = append(queryList, fmt.Sprintf(`"id" NOT IN (%s)`, params.ExcludeID))
		}

		if params.MemberID != "" {
			queryList = append(queryList, fmt.Sprintf(userIDparams, params.MemberID))
		}

		queryList = append(queryList, fmt.Sprintf(`"grantType" NOT IN ('%s','%s')`, authModel.AuthTypeLDAP, authModel.AuthTypeAzure))

		if len(queryList) > 0 {
			queryParam = fmt.Sprintf(whereParams, strings.Join(queryList, " AND "))
		}

		q := fmt.Sprintf(`SELECT "id","userId", "userName", ip, "userAgent",
						"deviceId", "clientType", "grantType", jti, "createdAt" 
						from session_info %s
						ORDER BY "createdAt" DESC LIMIT %d OFFSET %d `,
			queryParam, params.Limit, params.Offset)
		tags[helper.TextQuery] = q
		stmt, err := qp.repo.Query(q)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		var listSessionInfo model.SessionInfoList

		for stmt.Next() {
			var session model.SessionInfoResponse

			err = stmt.Scan(
				&session.ID,
				&session.UserID,
				&session.UserName,
				&session.IP,
				&session.UserAgent,
				&session.DeviceID,
				&session.ClientType,
				&session.GrantType,
				&session.JTI,
				&session.CreatedAt,
			)

			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
				tags[helper.TextResponse] = err
				output <- ResultQuery{Error: err}
				return
			}

			listSessionInfo.Data = append(listSessionInfo.Data, session)
		}

		output <- ResultQuery{Result: listSessionInfo}
	})
	return output
}

// GetTotalHistorySessionInfo function for getting total of history session
func (qp *SessionInfoQueryPostgres) GetTotalHistorySessionInfo(ctxReq context.Context, params *memberModel.ParametersLoginActivity) <-chan ResultQuery {
	ctx := "SessionInfoQueryPostgres-GetTotalHistorySessionInfo"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		var totalData int
		var (
			queryParam string
			queryList  []string
		)

		queryList = append(queryList, `"createdAt" > (CURRENT_DATE - INTERVAL '30 DAY')::DATE`)
		if params.ExcludeID != "" {
			queryList = append(queryList, fmt.Sprintf(`"id" NOT IN (%s)`, params.ExcludeID))
		}

		if params.MemberID != "" {
			queryList = append(queryList, fmt.Sprintf(userIDparams, params.MemberID))
		}

		if len(queryList) > 0 {
			queryParam = fmt.Sprintf(whereParams, strings.Join(queryList, " AND "))
		}
		sq := fmt.Sprintf(`SELECT count(id) FROM session_info %s`, queryParam)

		tags[helper.TextQuery] = sq
		stmt, err := qp.repo.Prepare(ctxReq, sq)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, params)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		err = stmt.QueryRow().Scan(&totalData)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, sq)
			output <- ResultQuery{Error: err}
			return
		}

		tags[helper.TextResponse] = totalData
		output <- ResultQuery{Result: totalData}
	})

	return output
}
