package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var BaseTestExtractClaims = []struct {
	name      string
	wantError bool
	token     interface{}
}{
	{
		name:      "#1 All OK",
		wantError: false,
		token: &jwt.Token{
			Claims: &BearerClaims{
				DeviceID:       "pzm",
				Adm:            true,
				StandardClaims: jwt.StandardClaims{Subject: "ui"},
			},
		},
	},
	{
		name:      "#2 Error format token",
		wantError: true,
		token: jwt.Token{
			Claims: &BearerClaims{
				DeviceID:       "pzm-pp",
				Adm:            false,
				StandardClaims: jwt.StandardClaims{Subject: "ui-op"},
			},
		},
	},
	{
		name:      "#3 error claims",
		wantError: true,
		token: &jwt.Token{
			Claims: BearerClaims{
				DeviceID:       "pzm-po",
				Adm:            true,
				StandardClaims: jwt.StandardClaims{Subject: "ui-opo"},
			},
		},
	},
	{
		name:      "#4 error claims",
		wantError: true,
		token:     "some",
	},
}

func TestExtractClaims(t *testing.T) {
	for _, tt := range BaseTestExtractClaims {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("token", tt.token)
		_, err := ExtractClaimsFromToken(c)
		if tt.wantError {
			assert.Error(t, err)
		}
	}
}

var BaseTestExtractMember = []struct {
	name      string
	wantError bool
	token     interface{}
	result    interface{}
}{
	{
		name:      "#1 All OK",
		wantError: false,
		token: &jwt.Token{
			Claims: &BearerClaims{
				DeviceID:       "pzm",
				UserAuthorized: true,
				Adm:            true,
				StandardClaims: jwt.StandardClaims{Subject: "mm"},
			},
		},
	},
	{
		name:      "#2 claim deviceID",
		wantError: false,
		token: &jwt.Token{
			Claims: &BearerClaims{
				DeviceID:       "pm",
				UserAuthorized: false,
				Adm:            true,
				StandardClaims: jwt.StandardClaims{Subject: "ss"},
			},
		},
	},
	{
		name:      "#3 error claim",
		wantError: true,
		token:     map[string]interface{}{"name": "value"},
	},
	{
		name:      "#4",
		wantError: true,
		token: &jwt.Token{
			Claims: &BearerClaims{
				UserAuthorized: true,
			},
		},
	},
}

func TestExtractMemberID(t *testing.T) {
	for _, tt := range BaseTestExtractMember {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("token", tt.token)
		result, err := ExtractMemberIDFromToken(c)
		if tt.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, result)
		}

	}
}

func TestExtractIsAdmin(t *testing.T) {
	for _, tt := range BaseTestExtractMember {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("token", tt.token)
		err := ExtractClaimsIsAdmin(c)
		if tt.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
