package helper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/labstack/echo"

	"github.com/xeipuuv/gojsonschema"
)

//Empty struct
type Empty struct{}

//EmptySlice type empty slice
type EmptySlice []Empty

// Meta map
type Meta map[string]interface{}

// JSONSchemaTemplate data structure
type JSONSchemaTemplate struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// SetData function
func (t *JSONSchemaTemplate) SetData(data interface{}) {
	t.Data = data
}

// ShowHTTPResponse function
func (t *JSONSchemaTemplate) ShowHTTPResponse(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(t.Code)
	return json.NewEncoder(c.Response()).Encode(t)
}

// ValidateJSONSchema function, function that validate JSON Schema structure
func ValidateJSONSchema(schema, data string) error {
	doc := gojsonschema.NewStringLoader(data)
	sc := gojsonschema.NewStringLoader(schema)

	result, err := gojsonschema.Validate(sc, doc)
	if err != nil {
		fmt.Println(err)
	}

	if !result.Valid() {
		return errors.New("Invalid json schema")
	}

	return nil
}
