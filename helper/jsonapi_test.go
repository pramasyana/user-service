package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testJSONAPIDatas = []struct {
	name      string
	wantError bool
	payload   interface{}
}{
	{
		"1",
		false,
		&ModelTest{Name: "pzm"},
	},
	{
		"2",
		true,
		ModelTest{Name: "pzm"},
	},
}

type ModelTest struct {
	Name string `json:"name"`
}

func TestJSONAPI(t *testing.T) {
	for _, tt := range testJSONAPIDatas {
		_, err := MarshalConvertOnePayload(tt.payload)
		if tt.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}
