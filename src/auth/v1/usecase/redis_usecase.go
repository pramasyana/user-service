package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/golang-jwt/jwt"
)

const (
	expiredToken = "this token has been expired"
	invalidToken = "malformed token"
)

// GetJTIToken function for get jti from token
func (au *AuthUseCaseImpl) GetJTIToken(ctxReq context.Context, token, request string) (string, interface{}, error) {
	ctx := "AuthUseCase-GetJTIToken"

	trace := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		trace.Finish(tags)
	}()
	claims := jwt.MapClaims{}

	jwtResult, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != tkn.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", tkn.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if (jwtResult == nil && err != nil) || len(claims) == 0 {
		return "", nil, errors.New("invalid token")
	}
	if request != "logout" {
		// if time now is between in expired date
		expired, ok := claims["exp"].(float64)
		if !ok {
			return "", nil, errors.New(invalidToken)
		}
		nowDate, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		expiredDate, err := time.Parse(time.RFC3339, time.Unix(int64(expired), 0).Format(time.RFC3339))
		if err != nil {
			return "", nil, errors.New(invalidToken)
		}

		// check token expired
		if nowDate.After(expiredDate) {
			tags["nowDate"] = nowDate.String()
			tags["expiredDate"] = expiredDate.String()
			err := errors.New(expiredToken)
			return "", nil, err
		}
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", nil, errors.New(invalidToken)
	}
	deviceID, ok := claims["did"].(string)
	if !ok {
		return "", nil, errors.New(invalidToken)
	}
	deviceLogin, ok := claims["dli"].(string)
	if !ok {
		return "", nil, errors.New(invalidToken)
	}

	// REDIS KEY FORMAT LOGIN
	redisKey := strings.Join([]string{"STG", sub, deviceID, deviceLogin}, "-")
	tags["key"] = redisKey

	return redisKey, claims, nil
}

func getRefreshTokenKey(subject, deviceID, deviceLogin string) string {
	return strings.Join([]string{"RT", subject, deviceID, deviceLogin}, "-")
}
