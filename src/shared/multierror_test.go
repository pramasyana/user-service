package shared

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	name       string
	errMessage string
	wantError  bool
}{
	{
		"#1",
		"something",
		true,
	},
}

func TestMultiError(t *testing.T) {
	for _, tc := range testCases {
		me := NewMultiError()
		if tc.wantError {
			me.Append("somekey", errors.New(tc.errMessage))
			assert.Contains(t, me.Error(), tc.errMessage)
			me.Append("somekey", errors.New("other message"))
		}

		assert.Equal(t, me.HasError(), tc.wantError)
		me.Clear()
	}
}
