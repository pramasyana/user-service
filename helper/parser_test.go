package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	var testDataParse = []struct {
		input  int32
		parsed string
	}{
		{
			input:  16463,
			parsed: "2015-01-28",
		},
	}

	for _, tc := range testDataParse {
		parsed := DateSinceEpoch(tc.input)
		assert.Equal(t, tc.parsed, parsed.Format(FormatDateDB))
	}
}
