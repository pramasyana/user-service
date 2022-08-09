package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	m := GenerateMemberIDv2()
	assert.NotEqual(t, "USR21041757518556", m)
	assert.Less(t, len(m), 20)
	// USR21049401747183

	n := GenerateDocumentID()
	assert.Less(t, len(n), 20)
}
