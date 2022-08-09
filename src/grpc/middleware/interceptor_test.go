package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	authorization = "authorization"
	methodTest    = "TestService.UnaryMethod"
	userService   = "user-service"
)

func TestAuthInterceptor(t *testing.T) {
	testsAuth := []struct {
		name, auth, req         string
		wantError, errorContext bool
	}{
		{
			name:      "PositiveTest",
			auth:      "user-service-9cc18122bcb544031798a8b1b9003c38",
			req:       "JUST TEST 1",
			wantError: false,
		},
		{
			name:      "NegativeTest",
			auth:      "invalid-auth",
			req:       "JUST TEST 2",
			wantError: true,
		},
		{
			name:      "MissingAuthorization",
			auth:      "",
			req:       "JUST TEST 3",
			wantError: true,
		},
		{
			name:         "MissingContext",
			auth:         "user-service-9cc18122bcb544031798a8b1b9003c38",
			req:          "JUST TEST 4",
			wantError:    true,
			errorContext: true,
		},
	}

	for _, tt := range testsAuth {

		t.Run(tt.name, func(t *testing.T) {
			md := metadata.Pairs(authorization, tt.auth)
			ctx := context.Background()
			if !tt.errorContext {
				ctx = metadata.NewIncomingContext(context.Background(), md)
			}

			unaryInfo := &grpc.UnaryServerInfo{
				FullMethod: methodTest,
			}

			unaryHandler := func(_ context.Context, _ interface{}) (interface{}, error) {
				return userService, nil
			}

			_, err := Auth(ctx, tt.req, unaryInfo, unaryHandler)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
