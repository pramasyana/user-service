package helper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/stretchr/testify/assert"
)

func TestJsonSchema(t *testing.T) {
	a := JSONSchemaTemplate{}
	a.SetData(map[string]string{"a": "1"})
	b := a.Data.(map[string]string)
	assert.Equal(t, "1", b["a"])

	app := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := app.NewContext(req, rec)
	err := a.ShowHTTPResponse(c)
	assert.NoError(t, err)
}

const (
	loginParam = `{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"description": "Request body for login",
		"type": "object",
		"properties": {
			"appId": {
				"type": "string"
			},
			"appSecret": {
				"type": "string"
			}
		},
		"required": ["appId","appSecret"]
	}`
	badLoginParam = `{
		"description": "Request body for login",
		"type": "object",
		"properties": {
			"appId": {
				"type": "string"
			},
			"appSecret": {
				"type": "string"
			}
		},
		"required": ["appId","appSecret"]
	}`

	loginBody = `{
		"appId": "someAppId",
		"appSecret":"someAppSecret"
	}`
)

func TestValidateJSONSchema(t *testing.T) {
	err := ValidateJSONSchema(loginParam, loginBody)
	assert.NoError(t, err)

	err = ValidateJSONSchema(loginParam, `{"appId":"ss", "appSecrets":"som"}`)
	assert.Error(t, err)
}
