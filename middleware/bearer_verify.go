package middleware

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/config/redis"
	jwt "github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

const (
	msgAuthEmpty      = "authorization is empty"
	eventErrorToken   = "error_token"
	unknownErrorToken = "Unknown token error"
)

// BearerClaims data structure for claims
type BearerClaims struct {
	Adm            bool   `json:"adm"`
	DeviceID       string `json:"did"`
	DeviceLogin    string `json:"dli"`
	Email          string `json:"email"`
	UserAuthorized bool   `json:"authorised"`
	JTI            string `json:"jti"`
	jwt.StandardClaims
}

// BearerVerify function to verify token
func BearerVerify(rsaPublicKey *rsa.PublicKey, cl redis.Client, mustAuthorized bool, mustAdmin bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if os.Getenv("NO_TOKEN") == "1" {
				return next(c)
			}

			req := c.Request()
			header := req.Header
			auth := header.Get("Authorization")

			ctx := "verifyToken"
			ctxReq := c.Request().Context()

			trace := tracer.StartTrace(ctxReq, ctx)
			trace.InjectHTTPHeader(req)

			var responseBody *echo.HTTPError

			tags := make(map[string]interface{})
			defer func() {
				tags["http.headers"] = req.Header
				tags["http.method"] = req.Method
				tags["http.url"] = req.URL.String()
				if responseBody != nil {
					tags["response.status_code"] = http.StatusUnauthorized
					tags["response.body"] = responseBody
				}
				trace.Finish(tags)
			}()

			tokenStr, err := getTokenString(auth)
			if err != nil {
				tracer.Log(ctxReq, eventErrorToken, err)
				return err
			}

			token, err := parseToken(ctxReq, rsaPublicKey, tokenStr)

			if err != nil {
				responseBody = getTokenError(tokenStr, err)
				tracer.Log(ctxReq, eventErrorToken, responseBody)
				return getTokenError(tokenStr, err)
			}

			return setClaims(c, next, token, cl, mustAdmin, mustAuthorized)
		}
	}
}

func parseToken(ctxReq context.Context, rsaPublicKey *rsa.PublicKey, tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			responseBody := echo.NewHTTPError(http.StatusUnauthorized, err)
			tracer.Log(ctxReq, eventErrorToken, responseBody)
			return nil, err
		}
		return rsaPublicKey, nil
	})
}
func getTokenString(auth string) (tokenStr string, err error) {
	responseBody := echo.NewHTTPError(http.StatusUnauthorized, msgAuthEmpty)
	if len(auth) <= 0 {
		return tokenStr, responseBody
	}
	splitToken := strings.Split(auth, " ")
	if len(splitToken) < 2 {
		return tokenStr, responseBody
	}
	if splitToken[0] != "Bearer" || splitToken[1] == "" {
		return tokenStr, responseBody
	}
	tokenStr = splitToken[1]
	return
}

func getTokenError(tokenStr string, err error) *echo.HTTPError {
	responseBody := echo.NewHTTPError(http.StatusUnauthorized, unknownErrorToken)
	if ve, ok := err.(*jwt.ValidationError); ok {
		var errorStr string
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			errorStr = fmt.Sprintf("Invalid token format: %s", tokenStr)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			errorStr = "Token has been expired"
		} else {
			errorStr = fmt.Sprintf("Token Parsing Error: %s", err.Error())
		}
		responseBody = echo.NewHTTPError(http.StatusUnauthorized, errorStr)
	}
	return responseBody
}

func setClaims(c echo.Context, next echo.HandlerFunc, token *jwt.Token, cl redis.Client, mustAdmin, mustAuthorized bool) error {
	responseBody := echo.NewHTTPError(http.StatusUnauthorized, unknownErrorToken)
	ctxReq := c.Request().Context()
	if claims, ok := token.Claims.(*BearerClaims); token.Valid && ok {

		if err := validateClaims(cl, mustAdmin, claims); err != nil {
			tracer.Log(ctxReq, eventErrorToken, err)
			return err
		}

		if mustAuthorized {
			if claims.UserAuthorized {
				c.Set("token", token)
				return next(c)
			}
			responseBody = echo.NewHTTPError(http.StatusUnauthorized, "Resource need an authorised user")
			tracer.Log(ctxReq, eventErrorToken, responseBody)
			return responseBody
		}

		c.Set("token", token)
		return next(c)
	}

	tracer.Log(ctxReq, eventErrorToken, responseBody)
	return responseBody
}

func validateClaims(cl redis.Client, mustAdmin bool, claims *BearerClaims) error {
	if os.Getenv("VALIDATE_BEARER_REDIS") == "true" {
		sub := claims.Subject
		deviceID := claims.DeviceID
		deviceLogin := claims.DeviceLogin

		redisKey := strings.Join([]string{"STG", sub, deviceID, deviceLogin}, "-")
		val, err := cl.Get(redisKey)

		if err != nil || val == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token has been expired")
		}

	}
	if mustAdmin && !claims.Adm {
		return echo.NewHTTPError(http.StatusUnauthorized, "Resource need an authorised user")
	}
	return nil
}
