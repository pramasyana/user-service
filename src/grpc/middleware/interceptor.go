package middleware

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

//Auth function,
//or Unary interceptor
//additional security for our GRPC server
func Auth(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	meta, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	if len(meta["authorization"]) != 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid authorization")
	}

	authorization := meta["authorization"][0]

	if authorization != "user-service-9cc18122bcb544031798a8b1b9003c38" {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid authorization")
	}

	return handler(ctx, req)
}
