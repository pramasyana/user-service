package shared

import "context"

// ContextKey type
type ContextKey string

// SetDataToContext set token to context
func SetDataToContext(ctx context.Context, key ContextKey, data interface{}) context.Context {
	return context.WithValue(ctx, key, data)
}

// GetDataFromContext set token to context
func GetDataFromContext(ctx context.Context, key ContextKey) interface{} {
	return ctx.Value(key)
}
