package usecase

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/config/rsa"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/golang-jwt/jwt"
)

const (
	errTokenMismatch = "token mismatch"
)

// VerifyTokenMember usecase
func (au *AuthUseCaseImpl) VerifyTokenMember(ctxReq context.Context, tokenStr string) (result ResultUseCase) {
	ctx := "Auth-VerifyTokenMember"

	trace := tracer.StartTrace(ctxReq, ctx)
	defer trace.Finish(nil)

	redisKey, dataToken, err := au.GetJTIToken(ctxReq, tokenStr, "")
	if err != nil {
		tracer.Log(ctxReq, ctx, err)
		result.Error = err
		return
	}

	getTokenRedis := <-au.LoginSessionRepo.Load(ctxReq, redisKey)

	if getTokenRedis.Error != nil {
		result.Error = errors.New(msgTokenExpired)
		return
	}
	tokenFromRedis, ok := getTokenRedis.Result.(model.LoginSessionRedis)
	if !ok {
		result.Error = errors.New("load result is not token")
		return
	}
	if tokenFromRedis.Token != tokenStr {
		result.Error = errors.New(errTokenMismatch)
		return
	}

	claims := dataToken.(jwt.MapClaims)

	responseVerify := au.GenerateResponseVerifyByID(ctxReq, claims, tokenStr)

	result.Result = responseVerify
	return
}

// VerifyTokenMemberB2b usecase
func (au *AuthUseCaseImpl) VerifyTokenMemberB2b(ctxReq context.Context, tokenStr, transaction_type, member_type string) <-chan ResultUseCase {
	ctx := "Auth-VerifyTokenMemberB2b"

	trace := tracer.StartTrace(ctxReq, ctx)
	defer trace.Finish(nil)

	output := make(chan ResultUseCase, 1)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		claims := jwt.MapClaims{}
		jwtResult, err := jwt.ParseWithClaims(tokenStr, claims, func(tkn *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if (jwtResult == nil && err != nil) || len(claims) == 0 {
			output <- ResultUseCase{Error: errors.New("invalid token")}
			return
		}

		redisLpseIDKey := strings.Join([]string{"STG-BELA", claims["sub"].(string), claims["email"].(string)}, "-")
		getLpseIDRedis := <-au.LoginSessionRepo.Load(ctxReq, redisLpseIDKey)
		if getLpseIDRedis.Error != nil {
			output <- ResultUseCase{Error: errors.New("LpseID not found")}
			return
		}

		lpseID := ""
		if existingLpseID, ok := getLpseIDRedis.Result.(model.LoginSessionRedis); ok {
			lpseID = existingLpseID.Token
		}

		// generate new token
		data := model.RequestToken{
			Email:           strings.ToLower(claims["email"].(string)),
			TransactionType: transaction_type,
			MemberType:      member_type,
			Audience:        claims["aud"].(string),
			DeviceID:        claims["did"].(string),
			DeviceLogin:     claims["dli"].(string),
			LpseID:          lpseID,
			TokenBela:       claims["customToken"].(string),
		}

		generateTokenB2b := <-au.GenerateTokenB2b(ctxReq, "", data)
		reqToken, ok := generateTokenB2b.Result.(model.RequestToken)
		if !ok {
			output <- ResultUseCase{Error: errors.New("result is not token")}
			return
		}
		output <- ResultUseCase{Result: reqToken}
	})

	return output
}

// GenerateResponseVerify function for generate response
func (au *AuthUseCaseImpl) GenerateResponseVerify(ctxReq context.Context, claims jwt.MapClaims, tokenStr string) interface{} {
	ctx := "Auth-GenerateResponseVerify"

	trace := tracer.StartTrace(ctxReq, ctx)
	defer trace.Finish(nil)

	responseVerify := &model.VerifyResponse{}
	responseVerify.Adm = claims["adm"].(bool)
	responseVerify.Aud = claims["aud"].(string)
	responseVerify.Authorised = claims["authorised"].(bool)
	responseVerify.Did = claims["did"].(string)
	responseVerify.Dli = claims["dli"].(string)
	responseVerify.Email = claims["email"].(string)
	responseVerify.Exp = claims["exp"].(float64)
	responseVerify.Iat = claims["iat"].(float64)
	responseVerify.Iss = claims["iss"].(string)
	responseVerify.Jti = claims["jti"].(string)
	responseVerify.MemberType = claims["memberType"].(string)
	responseVerify.Staff = claims["staff"].(bool)
	responseVerify.Sub = claims["sub"].(string)
	responseVerify.SignUpFrom = claims["signUpFrom"].(string)
	responseVerify.Token = tokenStr
	responseVerify.IsMerchant = false

	t := time.Unix(int64(responseVerify.Exp), 0)
	responseVerify.ExpiredTime = t

	pubKey, err := rsa.InitPublicKey()
	if err != nil {
		tracer.Log(ctxReq, ctx, err)
	}

	// validate old token
	oldClaims, _ := token.VerifyTokenIgnoreExpiration(pubKey, tokenStr)
	RTKey := getRefreshTokenKey(oldClaims.Subject, oldClaims.DeviceID, oldClaims.DeviceLogin)
	refreshResult := <-au.RefreshTokenRepo.Load(ctxReq, RTKey)
	rt, _ := refreshResult.Result.(model.RefreshToken)
	responseVerify.RefreshToken = rt.Token

	if responseVerify.Email != "" {
		responseVerify = au.adjustVerifyData(ctxReq, responseVerify)
	}
	return responseVerify
}

// adjustVerifyData for added detail member
func (au *AuthUseCaseImpl) adjustVerifyData(ctxReq context.Context, responseVerify *model.VerifyResponse) *model.VerifyResponse {
	if responseVerify.MemberType == model.UserTypeCorporate {
		emailCorporateResult := <-au.CorporateContactQueryRead.FindByEmail(ctxReq, responseVerify.Email)
		if emailCorporateResult.Result != nil {
			corporateContact := emailCorporateResult.Result.(sharedModel.B2BContactData)
			responseVerify.FirstName = corporateContact.FirstName
			responseVerify.LastName = corporateContact.LastName
			responseVerify.UserID = strconv.Itoa(corporateContact.ID)
			if len(corporateContact.PhoneNumber) != 0 {
				responseVerify.Mobile = corporateContact.PhoneNumber
			}
		}
	} else {
		emailResult := <-au.MemberQueryRead.FindByEmail(ctxReq, responseVerify.Email)
		if emailResult.Result != nil {
			member := emailResult.Result.(memberModel.Member)
			responseVerify.FirstName = member.FirstName
			responseVerify.LastName = member.LastName
			responseVerify.UserID = member.ID
			responseVerify.Mobile = member.Mobile
			merchantResult := au.MerchantRepoRead.FindMerchantByUser(ctxReq, responseVerify.UserID)
			if merchantResult.Result != nil {
				responseVerify.IsMerchant = true
			}
		}
	}
	return responseVerify
}

// GenerateResponseVerifyByID function for generate response
func (au *AuthUseCaseImpl) GenerateResponseVerifyByID(ctxReq context.Context, claims jwt.MapClaims, tokenStr string) interface{} {
	ctx := "Auth-GenerateResponseVerifyByID"

	trace := tracer.StartTrace(ctxReq, ctx)
	defer trace.Finish(nil)

	responseVerify := &model.VerifyResponse{}
	responseVerify.Adm = claims["adm"].(bool)
	responseVerify.Aud = claims["aud"].(string)
	responseVerify.Authorised = claims["authorised"].(bool)
	responseVerify.Did = claims["did"].(string)
	responseVerify.Dli = claims["dli"].(string)
	responseVerify.Email = strings.ToLower(claims["email"].(string))
	responseVerify.Exp = claims["exp"].(float64)
	responseVerify.Iat = claims["iat"].(float64)
	responseVerify.Iss = claims["iss"].(string)
	responseVerify.Jti = claims["jti"].(string)
	responseVerify.MemberType = claims["memberType"].(string)
	responseVerify.Staff = claims["staff"].(bool)
	responseVerify.Sub = claims["sub"].(string)
	responseVerify.SignUpFrom = claims["signUpFrom"].(string)
	responseVerify.UserID = responseVerify.Sub
	responseVerify.Token = tokenStr
	responseVerify.IsMerchant = false
	responseVerify.AccountID = ""
	responseVerify.CustomToken = claims["customToken"].(string)

	t := time.Unix(int64(responseVerify.Exp), 0)
	responseVerify.ExpiredTime = t

	pubKey, err := rsa.InitPublicKey()
	if err != nil {
		tracer.Log(ctxReq, ctx, err)
	}

	// validate old token
	oldClaims, err := token.VerifyTokenIgnoreExpiration(pubKey, tokenStr)
	if err == nil {
		RTKey := getRefreshTokenKey(oldClaims.Subject, oldClaims.DeviceID, oldClaims.DeviceLogin)
		refreshResult := <-au.RefreshTokenRepo.Load(ctxReq, RTKey)
		rt, _ := refreshResult.Result.(model.RefreshToken)
		responseVerify.RefreshToken = rt.Token
	}

	if responseVerify.Email != "" && responseVerify.UserID != "" {
		responseVerify = au.AdjustMemberData(ctxReq, responseVerify)
	}
	return responseVerify
}

// Logout usecase function for user logout
func (au *AuthUseCaseImpl) Logout(ctxReq context.Context, token string) <-chan ResultUseCase {
	ctx := "AuthUseCase-Logout"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		tags["args"] = token

		redisKey, dataToken, err := au.GetJTIToken(ctxReq, token, "logout")
		if err != nil || dataToken == nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		removeToken := <-au.LoginSessionRepo.Delete(ctxReq, redisKey)

		if removeToken.Error != nil {
			output <- ResultUseCase{Error: errors.New("invalid token format"), HTTPStatus: http.StatusBadRequest}
			return
		}

		claims := dataToken.(jwt.MapClaims)
		subject := claims["sub"].(string)
		deviceID := claims["did"].(string)
		deviceLogin := claims["dli"].(string)
		RTokenKey := getRefreshTokenKey(subject, deviceID, deviceLogin)

		mm := <-au.RefreshTokenRepo.Delete(ctxReq, RTokenKey)
		if mm.Error != nil {
			output <- ResultUseCase{Error: errors.New("error delete refresh token"), HTTPStatus: http.StatusBadRequest}
			return
		}
		response := model.Logout{Token: token}
		tags[helper.TextResponse] = response
		output <- ResultUseCase{Result: response}

	})
	return output
}
