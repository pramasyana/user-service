package shared

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
	IsStaff        bool   `json:"staff"`
	jwt.StandardClaims
}

// JWTExtract function for extract token
func JWTExtract(rsaPublicKey *rsa.PublicKey, accessToken string) (*BearerClaims, error) {
	var errorStr error

	token, err := jwt.ParseWithClaims(accessToken, &BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	if claims, ok := token.Claims.(*BearerClaims); err == nil && token.Valid && ok {
		return claims, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {

		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			errorStr = fmt.Errorf("Invalid token format: %s", accessToken)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			errorStr = errors.New("Token has been expired")
		} else {
			errorStr = fmt.Errorf("Token Parsing Error: %s", err.Error())
		}

		return nil, errorStr
	}

	errorStr = errors.New("unknown errors")
	return nil, errorStr
}
