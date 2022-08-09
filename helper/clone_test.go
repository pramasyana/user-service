package helper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cloner struct {
	password string
}

func TestClone(t *testing.T) {
	var testData = []struct {
		name   string
		input  interface{}
		output interface{}
	}{
		{
			name: "test 1",
			input: &cloner{
				password: "something",
			},
			output: &cloner{
				password: "something",
			},
		},
	}
	for _, tc := range testData {
		CloneStruct(tc.input, tc.output)
		memAddress := fmt.Sprintf("%p", tc.input)
		memAddress2 := fmt.Sprintf("%p", tc.output)
		assert.NotEqual(t, memAddress, memAddress2)
	}
}
