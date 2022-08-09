package middleware

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Bhinneka/user-service/src/shared/mocks"
	jwt "github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	expiredToken         = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0Ijp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiZXhwIjoxNTIxMDIyNTI0LCJpYXQiOjE1MjEwMjI0NjQsImlzcyI6ImJoaW5uZWthLmNvbSIsInN1YiI6IlVTUjE4MDIxODE3NSJ9.gw_tfbPHq6XIWuVT3ksHFovWkYZteUuuLGepkGeAnAGP41pF_AlEhcB120Jao1FOi74li7f3ab6kN_hBcMNMuYmhlkP4pa78QKFCY7uZpLUm6LIc2AOHf1VRm0poQnvH0AnDDw1_bU8NFe0GKr48Cf88934txTCJRQ75Sw4pnbk"
	anonymousAccessToken = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0aG9yaXNlZCI6ZmFsc2UsImRpZCI6ImMwYjRkMWI0YzQ0NzQiLCJkbGkiOiJXRUIiLCJleHAiOjE1NDYyNDE2ODUsImlhdCI6MTUyNDY0MTY4NSwiaXNzIjoiYmhpbm5la2EuY29tIiwic3ViIjoiYmhpbm5la2EtbWljcm9zZXJ2aWNlcy1iMTM3MTQtNTMxMjExNSJ9.U0XGYR_ZiYUDA1aRHgvbYqp0zF8saN-O4_H7793Ou8OfQKsU-t5NRqC6cyeImBN8ayh3o3s35_4AvAuGH9uLrlq3MuUxPZnsmU6SlJ0ODen0O8ak6i3PZF_6dltLA5ZsTgVt4YvOydSQusTqHoN-jyFsBFtZHKHCFoFVK9Hqu-Q"
	dummyToken           = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	invalidRSAToken      = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.POstGetfAytaZS82wHcjoTyoqhMyxXiWdR7Nn7A29DNSl0EiXLdwJ6xC6AfgZWF1bOsS_TuYI3OG85AmiExREkrS6tDfTQ2B3WXlrr-wp5AokiRbz3_oB4OxG-W9KcEEbDRcZc0nH3L7LzYptiy1PtAylQGxHTWZXtGz4ht0bAecBgmpdgXMguEIcoqPJ1n3pIWk_dUZegpqx0Lka21H6XxUTxiy8OcaarA8zdnPUnV6AmNP3ecFawIFYdvJB_cm-GvpCSbr8G8y_Mllj8f4x9nBH8pQux89_6gUY618iYv7tuPWBFfEbLxtF2pZS6YC1aSfLQxeNe8djT9YjpvRZA"
	noTokenKey           = "NO_TOKEN"
	validateRedisKey     = "VALIDATE_BEARER_REDIS"
)

var (
	jwtTokenValidTest = &jwt.Token{
		Claims: &BearerClaims{UserAuthorized: true, Adm: false},
		Valid:  true,
	}
	jwtTokenValidAdmTest = &jwt.Token{
		Claims: &BearerClaims{UserAuthorized: true, Adm: true},
		Valid:  true,
	}
	jwtTokenInvalidTest = &jwt.Token{
		Claims: &BearerClaims{UserAuthorized: true},
		Valid:  false,
	}
	jwtTokenUnauthorizedTest = &jwt.Token{
		Claims: &BearerClaims{UserAuthorized: false},
		Valid:  true,
	}
)

func setupRequest(headers map[string]string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func setupHandler() echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.String(http.StatusOK, "bhinneka.com"))
	})
}

var testVerify = []struct {
	name      string
	wantError bool
	setHeader bool
	token     string
}{
	{
		"BearerVerify#1",
		true,
		true,
		expiredToken,
	},
	{
		"BearerVerify#2",
		true,
		true,
		anonymousAccessToken,
	},
	{
		"BearerVerify#3",
		true,
		false,
		"",
	},
	{
		"BearerVerify#4",
		true,
		true,
		"Some",
	},
	{
		"BearerVerify#5",
		true,
		true,
		"Some thing",
	},
	{
		"BearerVerify#6",
		true,
		true,
		dummyToken,
	},
	{
		"BearerVerify#7",
		true,
		true,
		invalidRSAToken,
	},
	{
		"BearerVerify#8",
		true,
		true,
		"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..",
	},
}

func TestBearerVerify(t *testing.T) {
	for _, tt := range testVerify {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(map[string]string)
			if tt.setHeader {
				headers[echo.HeaderAuthorization] = tt.token
			}
			c := setupRequest(headers)
			verifyKey, _ := getPublicKey(ValidPublicKey)
			handler := setupHandler()
			fredis := mocks.InitFakeRedis()
			mw := BearerVerify(verifyKey, fredis, true, false)(handler)
			err := mw(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("POSITIVE_NO_TOKEN", func(t *testing.T) {
		os.Setenv(noTokenKey, "1")
		defer os.Unsetenv(noTokenKey)
		c := setupRequest(nil)
		verifyKey, _ := getPublicKey(ValidPublicKey)
		handler := setupHandler()
		fredis := mocks.InitFakeRedis()
		mw := BearerVerify(verifyKey, fredis, true, false)(handler)
		assert.NoError(t, mw(c))
	})
}

const ValidPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCoqzL5JrMzed4tb8uEoLKd42EO
sYmb0HpbicGt/OUeJxaHtt59Ew0BbpreBeiuugXweEa5xctQOxGYr27h4ZOnR0hW
Si+h5Y35CKzMEmZnzQwzQphgqww0U+e9/OAvVfCW1xWvVFr0WbhIRn+w/9DUvp+6
jKz3fIj3yQaHWVMMNQIDAQAB
-----END PUBLIC KEY-----`

func getPublicKey(publicKey string) (*rsa.PublicKey, error) {
	r := strings.NewReader(publicKey)
	verifyBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}
	return verifyKey, nil
}

func TestSetClaims(t *testing.T) {
	c := setupRequest(map[string]string{echo.HeaderAuthorization: anonymousAccessToken})
	next := setupHandler()
	cl := mocks.InitFakeRedis()

	type args struct {
		token          *jwt.Token
		mustAdmin      bool
		mustAuthorized bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "POSITIVE_SET_CLAIMS",
			args: args{
				token:          jwtTokenValidTest,
				mustAdmin:      false,
				mustAuthorized: false,
			},
			wantErr: false,
		},
		{
			name: "POSITIVE_SET_CLAIMS_MUST_ADMIN",
			args: args{
				token:          jwtTokenValidAdmTest,
				mustAdmin:      true,
				mustAuthorized: true,
			},
			wantErr: false,
		},
		{
			name: "NEGATIVE_SET_CLAIMS_MUST_ADMIN",
			args: args{
				token:          jwtTokenValidTest,
				mustAdmin:      true,
				mustAuthorized: true,
			},
			wantErr: true,
		},
		{
			name: "NEGATIVE_SET_CLAIMS_UNAUTHORIZED",
			args: args{
				token:          jwtTokenUnauthorizedTest,
				mustAdmin:      false,
				mustAuthorized: true,
			},
			wantErr: true,
		},
		{
			name: "NEGATIVE_SET_CLAIMS_INVALID_TOKEN",
			args: args{
				token:          jwtTokenInvalidTest,
				mustAdmin:      false,
				mustAuthorized: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setClaims(c, next, tt.args.token, cl, tt.args.mustAdmin, tt.args.mustAuthorized); (err != nil) != tt.wantErr {
				t.Errorf("setClaims() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("NEGATIVE_SET_CLAIMS_REDIS", func(t *testing.T) {
		os.Setenv(validateRedisKey, "true")
		defer os.Unsetenv(validateRedisKey)
		err := setClaims(c, next, jwtTokenValidTest, cl, false, false)
		assert.Error(t, err)
	})
}
