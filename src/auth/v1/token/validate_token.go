package token

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

//BearerClaims data structure
type BearerClaims struct {
	DeviceID       string `json:"did"`
	DeviceLogin    string `json:"dli"`
	Email          string `json:"email"`
	UserAuthorized bool   `json:"authorised"`
	IsAdmin        bool   `json:"adm"`
	MemberType     string `json:"memberType"`
	jwt.StandardClaims
}

// VerifyTokenIgnoreExpiration function for verifying token and ignore expiration
func VerifyTokenIgnoreExpiration(rsaPublicKey *rsa.PublicKey, oldAccessToken string) (*BearerClaims, error) {
	var errorStr error

	token, err := jwt.ParseWithClaims(oldAccessToken, &BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			errorStr = fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, errorStr
		}
		return rsaPublicKey, nil
	})

	if token == nil {
		errorStr = errors.New("invalid old token")
		return nil, errorStr
	}

	if claims, ok := token.Claims.(*BearerClaims); ok {
		return claims, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {

		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			errorStr = fmt.Errorf("invalid token format: %s", oldAccessToken)
		} else {
			errorStr = fmt.Errorf("token parsing error: %s", err.Error())
		}

		return nil, errorStr
	}

	errorStr = errors.New("unknown errors")
	return nil, errorStr
}
