package usecase

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	sessionModel "github.com/Bhinneka/user-service/src/session/v1/model"
	"github.com/golang-jwt/jwt"
)

// parseToken function for parsing token into email and real token
func (mu *MemberUseCaseImpl) parseToken(token string) (string, string, error) {
	if len(token) == 0 {
		return "", "", fmt.Errorf(helper.ErrorParameterRequired, "token")
	}

	splitTokens := strings.Split(token, "-")
	if len(splitTokens) < 2 {
		return "", "", fmt.Errorf(helper.ErrorParameterInvalid, "token")
	}

	emailTmp, err := base64.URLEncoding.DecodeString(splitTokens[1])
	if err != nil {
		return "", "", fmt.Errorf(helper.ErrorParameterInvalid, "token")
	}

	return string(emailTmp), splitTokens[0], nil
}

func (mu *MemberUseCaseImpl) flushAllTokenUser(ctxReq context.Context, member model.Member) <-chan ResultUseCase {
	ctx := "MemberUseCase-flushAllTokenUser"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		// Remove token user
		qSession := sessionModel.ParamList{
			StrPage:  "1",
			Page:     1,
			StrLimit: "999",
			Limit:    999,
			Offset:   0,
			Sort:     "asc",
			OrderBy:  "id",
			Email:    member.Email,
		}

		tags["param"] = qSession

		sessionData := <-mu.SessionQueryRead.GetListSessionInfo(ctxReq, &qSession)
		if sessionData.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_data_session", sessionData.Error, member)
			output <- ResultUseCase{Error: sessionData.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		sData, ok := sessionData.Result.(sessionModel.SessionInfoList)
		if !ok {
			output <- ResultUseCase{Error: errors.New("result is not session info list"), HTTPStatus: http.StatusBadRequest}
			return
		}

		for _, session := range sData.Data {
			if *session.JTI != "" {
				rediskey := strings.Join([]string{"STG", *session.UserID, *session.DeviceID, *session.ClientType}, "-")
				mu.LoginSessionRedis.Delete(ctxReq, rediskey)
			}
		}

		tags[helper.TextArgs] = member

		output <- ResultUseCase{Result: member}
	})

	return output
}

// RevokeAllAccess function for getting detail member based on member id
func (mu *MemberUseCaseImpl) RevokeAllAccess(ctxReq context.Context, uid, token string) <-chan ResultUseCase {
	ctx := "MemberUseCase-RevokeAllAccess"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags["args"] = uid
		if !strings.Contains(uid, usrFormat) {
			err := fmt.Errorf(helper.ErrorParameterInvalid, msgErrorMemberID)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberRepoRead.Load(ctxReq, uid)
		if memberResult.Error != nil {
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}

			helper.SendErrorLog(ctxReq, ctx, "member_repo_load_redis", memberResult.Error, memberResult.Result)
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		err := mu.revokeAllAccessProccess(ctxReq, uid, token, true)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "revoke_all_session", err, token)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: nil}

	})

	return output
}

func (mu *MemberUseCaseImpl) revokeAllAccessProccess(ctxReq context.Context, uid, token string, excludeActiveToken bool) error {
	activeKey := ""
	activeRefreshKey := ""
	if excludeActiveToken {
		activeKey, activeRefreshKey, _, _ = mu.ParseToken(ctxReq, token)
	}

	// revoke all session redis
	rediskey := strings.Join([]string{"STG", uid}, "-")
	revokeResult := <-mu.MemberRepoRedis.RevokeAllAccess(ctxReq, rediskey, activeKey)
	if revokeResult.Error != nil {
		return errors.New(model.ErrorRevokeAllAccess)
	}

	// revoke all session refresh token redis
	refreshKey := strings.Join([]string{"RT", uid}, "-")
	revokeRefreshResult := <-mu.MemberRepoRedis.RevokeAllAccess(ctxReq, refreshKey, activeRefreshKey)
	if revokeRefreshResult.Error != nil {
		return errors.New(model.ErrorRevokeAllAccess)
	}

	return nil
}

// ParseToken function for get jti from token
func (mu *MemberUseCaseImpl) ParseToken(ctxReq context.Context, token string) (string, string, interface{}, error) {
	ctx := "MemberUseCase-ParseToken"

	trace := tracer.StartTrace(ctxReq, ctx)
	defer trace.Finish(nil)

	claims := jwt.MapClaims{}

	jwtResult, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != tkn.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", tkn.Header["alg"])
		}

		return []byte(textSecret), nil
	})

	if jwtResult == nil && err != nil {
		err := errors.New("invalid token")
		return "", "", nil, err
	}

	sub := claims["sub"].(string)
	deviceID := claims["did"].(string)
	deviceLogin := claims["dli"].(string)

	// REDIS KEY FORMAT LOGIN
	redisKey := strings.Join([]string{"STG", sub, deviceID, deviceLogin}, "-")

	// REDIS REFREH KEY FORMAT LOGIN
	refreshKey := strings.Join([]string{"RT", sub, deviceID, deviceLogin}, "-")

	return redisKey, refreshKey, claims, nil
}

// checkAttemptResendActivation function for saving resend activation attempt data to redis based on email for 24 hours
func (mu *MemberUseCaseImpl) checkAttemptResendActivation(ctxReq context.Context, email string, lastTokenAttempt time.Time) <-chan ResultUseCase {
	ctx := "AuthUseCase-checkAttemptResendActivation"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// validate time from lastTokenAttempt (minutes)
		loc, _ := time.LoadLocation("UTC")
		now := time.Now().In(loc)
		diff := now.Sub(lastTokenAttempt)
		mins := int(diff.Minutes())
		ageReq, _ := strconv.Atoi(mu.ResendActivationAttemptAgeRequest)
		tags[helper.TextEmail] = email
		if mins < ageReq {
			err := fmt.Errorf(model.ErrorResendActivationTime)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err}
			return
		}

		attemptAge, err := time.ParseDuration(mu.ResendActivationAttemptAge)
		if err != nil {
			output <- ResultUseCase{Error: errors.New(model.ErrorResendActivation)}
			return
		}

		attempt := 1

		// get attempt first then increment the number
		key := fmt.Sprintf("ATTEMPT_RESEND_ACTIVATION:%s", email)
		data, err := mu.saveAttemptResendActivation(ctxReq, key, attempt, attemptAge)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err}
			return
		}

		if len(data.Attempt) > 0 && data.Attempt != "0" {
			intAttempt, _ := strconv.Atoi(data.Attempt)

			// by pass when attempt reach 10 times
			if intAttempt == 10 {
				err := errors.New(model.ErrorMaxAttemptResendActivation)
				tags[helper.TextResponse] = err
				output <- ResultUseCase{Error: err}
				return
			}
			attempt = intAttempt + 1
		}

		strAttempt := strconv.Itoa(attempt)
		newData := model.ResendActivationAttempt{
			Key:                        key,
			Attempt:                    strAttempt,
			ResendActivationAttemptAge: attemptAge,
		}

		saveResult := <-mu.MemberRepoRedis.SaveResendActivationAttempt(ctxReq, &newData)
		if saveResult.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_resend_activation_attempt", saveResult.Error, newData)
			output <- ResultUseCase{Error: saveResult.Error}
			return
		}

		output <- ResultUseCase{Result: saveResult}

	})

	return output
}

func (mu *MemberUseCaseImpl) saveAttemptResendActivation(ctxReq context.Context, key string, attempt int, attemptAge time.Duration) (model.ResendActivationAttempt, error) {
	ctx := "AuthUseCase-saveAttemptResendActivation"
	data := model.ResendActivationAttempt{}
	attResult := <-mu.MemberRepoRedis.LoadByKey(ctxReq, key)
	if attResult.Error != nil {
		if attResult.Error.Error() != helper.ErrorRedis {
			helper.SendErrorLog(ctxReq, ctx, "get_login_attempt", attResult.Error, key)
			return data, attResult.Error
		}

		// when redis nil
		strAttempt := strconv.Itoa(attempt)
		newData := model.ResendActivationAttempt{
			Key:                        key,
			Attempt:                    strAttempt,
			ResendActivationAttemptAge: attemptAge,
		}

		saveResult := <-mu.MemberRepoRedis.SaveResendActivationAttempt(ctxReq, &newData)
		if saveResult.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_login_attempt", attResult.Error, newData)
			return data, saveResult.Error
		}

		return data, nil
	}

	data, ok := attResult.Result.(model.ResendActivationAttempt)
	if !ok {
		err := errors.New(model.ErrorResendActivation)
		return data, err
	}
	return data, nil
}

// GetLoginActivity function for getting detail login activity based on member id
func (mu *MemberUseCaseImpl) GetLoginActivity(ctxReq context.Context, params *model.ParametersLoginActivity) <-chan ResultUseCase {
	ctx := "MemberUseCase-GetLoginActivity"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextMemberIDCamel] = params.MemberID

		memberResult := <-mu.MemberRepoRead.Load(ctxReq, params.MemberID)
		if memberResult.Error != nil {
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		// GET CURRENT ACTIVE TOKEN
		activeID, activeSession, err := mu.GetActiveSession(ctxReq, params)
		if err != nil {
			err := errors.New(model.ErrorGetLoginActivity)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// GET SESSION HISTORY
		historySession := <-mu.GetHistorySession(ctxReq, activeID, params)
		if historySession.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_history_session", historySession.Error, params)
			output <- ResultUseCase{Error: historySession.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		resultHistory := historySession.Result.(model.SessionHistoryInfoList)

		sessionInfo := model.SessionInfo{
			ActiveSession:  activeSession,
			HistorySession: resultHistory.Data,
		}

		resultSessionInfo := model.SessionInfoList{
			Data:      sessionInfo,
			TotalData: resultHistory.TotalData,
		}

		output <- ResultUseCase{Result: resultSessionInfo}

	})

	return output
}

// GetActiveSession function for getting active session login activity
func (mu *MemberUseCaseImpl) GetActiveSession(ctxReq context.Context, params *model.ParametersLoginActivity) ([]string, []model.SessionInfoDetail, error) {
	var activeID []string
	activeSession := []model.SessionInfoDetail{}

	// GET ACTIVE SESSION LOGIN ON REDIS
	rediskey := strings.Join([]string{"STG", params.MemberID}, "-")
	getLoginActive := <-mu.LoginSessionRedis.GetLoginActive(ctxReq, rediskey)
	if getLoginActive.Error != nil {
		return activeID, activeSession, getLoginActive.Error
	}
	dataLoginActive := getLoginActive.Result.([]authModel.LoginSessionRedis)

	activeKey, _, activeClaims, _ := mu.ParseToken(ctxReq, params.Token)
	claims := activeClaims.(jwt.MapClaims)
	data, err := mu.GetDetailSession(ctxReq, claims)
	if err == nil {
		data.ActiveNow = true
		activeID = append(activeID, fmt.Sprintf("'%s'", data.ID))
		activeSession = append(activeSession, data)
	}

	// GET OTHERS ACTIVE SESSION
	for _, value := range dataLoginActive {
		key, _, activeClaims, _ := mu.ParseToken(ctxReq, value.Token)
		if activeKey == key {
			continue
		}

		claims := activeClaims.(jwt.MapClaims)
		data, err := mu.GetDetailSession(ctxReq, claims)
		if data.GrantType == authModel.AuthTypeLDAP || data.GrantType == authModel.AuthTypeAzure || err != nil {
			continue
		}

		activeID = append(activeID, fmt.Sprintf("'%s'", data.ID))
		activeSession = append(activeSession, data)
	}

	return activeID, activeSession, nil
}

// GetDetailSession  function for getting detail session by param
func (mu *MemberUseCaseImpl) GetDetailSession(ctxReq context.Context, claims jwt.MapClaims) (model.SessionInfoDetail, error) {
	data := model.SessionInfoDetail{}
	param := sessionModel.ParametersGetSession{
		DeviceID:   claims["did"].(string),
		ClientType: claims["dli"].(string),
		UserID:     claims["sub"].(string),
	}
	session := <-mu.SessionQueryRead.GetDetailSessionInfo(ctxReq, param)

	if session.Error != nil {
		return data, session.Error
	}

	resultSession := session.Result.(sessionModel.SessionInfoResponse)
	browser, isMobile, isApp := helper.ParseUserAgent(*resultSession.UserAgent)

	data.ID = *resultSession.ID
	data.DeviceType = *resultSession.ClientType
	data.IP = *resultSession.IP
	data.UserAgent = browser
	data.ActiveNow = false
	data.LastLogin = resultSession.CreatedAt
	data.GrantType = *resultSession.GrantType
	data.IsMobile = isMobile
	data.IsApp = isApp

	return data, nil
}

// RevokeAccess function for revoke specific jti session
func (mu *MemberUseCaseImpl) RevokeAccess(ctxReq context.Context, memberID, sessionID string) <-chan ResultUseCase {
	ctx := "MemberUseCase-RevokeAccess"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags["args"] = sessionID

		memberData := <-mu.MemberRepoRead.Load(ctxReq, memberID)
		if memberData.Error != nil {
			if memberData.Error == sql.ErrNoRows {
				memberData.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}

			helper.SendErrorLog(ctxReq, ctx, "member_repo_revoke_access", memberData.Error, memberData.Result)
			output <- ResultUseCase{Error: memberData.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		sessionInfo := <-mu.SessionQueryRead.GetDetailSessionInfo(ctxReq, sessionModel.ParametersGetSession{UserID: memberID, SessionID: sessionID})
		if sessionInfo.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "member_repo_redis_get_detail", sessionInfo.Error, sessionInfo.Result)
			output <- ResultUseCase{Error: sessionInfo.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		revokeSession := sessionInfo.Result.(sessionModel.SessionInfoResponse)
		mu.revokeJTIToken(ctxReq, revokeSession)

		output <- ResultUseCase{Result: nil}
	})

	return output
}

func (mu *MemberUseCaseImpl) revokeJTIToken(ctxReq context.Context, params sessionModel.SessionInfoResponse) {
	// revoke all session redis
	mu.LoginSessionRedis.Delete(ctxReq, strings.Join([]string{"STG", *params.UserID, *params.DeviceID, *params.ClientType}, "-"))
	mu.LoginSessionRedis.Delete(ctxReq, strings.Join([]string{"RT", *params.UserID, *params.DeviceID, *params.ClientType}, "-"))
}
