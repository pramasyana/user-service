package helper

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestError(t *testing.T) {
	var defError = fmt.Errorf("default error")
	var testData = []struct {
		name    string
		err     error
		payload interface{}
		active  string
	}{
		{
			name: "Test error #0",
			err:  nil,
		},
		{
			name:    "Test error #1",
			err:     defError,
			active:  "true",
			payload: map[string]interface{}{"1": "2"},
		},
	}

	for _, tc := range testData {
		os.Setenv("SENTRY", tc.active)
		SendErrorLog(context.Background(), "someContext", "some_scope", tc.err, tc.payload)
	}
}
