package middleware

import (
	"errors"

	"github.com/Bhinneka/user-service/src/document/v2/model"
	jwt "github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

// ExtractMemberIDFromToken function for extracting user ID
// UserID can be DeviceID or MemberID
func ExtractMemberIDFromToken(c echo.Context) (string, error) {

	var userID string
	claims, err := ExtractClaimsFromToken(c)

	if err != nil {
		return userID, err
	}

	if claims.UserAuthorized {
		userID = claims.StandardClaims.Subject
	} else {
		userID = claims.DeviceID
	}

	if len(userID) == 0 {
		return userID, errors.New("a non authorised call should have DeviceId in the payload")
	}

	return userID, nil
}

// ExtractClaimsFromToken function for extracting claims
func ExtractClaimsFromToken(c echo.Context) (*BearerClaims, error) {
	tokenContext := c.Get("token")

	token, ok := tokenContext.(*jwt.Token)
	if !ok {
		return nil, errors.New("wrong token format")
	}
	if token == nil {
		return nil, errors.New("empty token")
	}

	claims, ok := token.Claims.(*BearerClaims)
	if !ok {
		return nil, errors.New("claims is in wrong type")
	}

	return claims, nil
}

// ExtractClaimsIsAdmin function for extracting claims
func ExtractClaimsIsAdmin(c echo.Context) error {
	claims, err := ExtractClaimsFromToken(c)
	if err != nil {
		return err
	}

	if !claims.Adm {
		return errors.New(model.Unauthorized)
	}

	return nil
}
