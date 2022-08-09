package shared

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := SetDataToContext(context.Background(), ContextKey("test"), "test")
	result := GetDataFromContext(ctx, ContextKey("test"))
	assert.Equal(t, "test", result)
}
