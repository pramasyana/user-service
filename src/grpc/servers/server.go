package servers

import (
	"fmt"
	"net"

	"github.com/Bhinneka/user-service/src/grpc/middleware"
	"google.golang.org/grpc"

	healthPB "github.com/Bhinneka/user-service/protogo/health"
	healthDelivery "github.com/Bhinneka/user-service/src/health/delivery"

	memberPB "github.com/Bhinneka/user-service/protogo/member"
	memberDelivery "github.com/Bhinneka/user-service/src/member/v1/delivery"
)

// Server data structure
type Server struct {
	healthGRPCHandler *healthDelivery.GRPCHandler
	memberGRPCHandler *memberDelivery.GRPCHandler
}

// NewGRPCServer function for creating GRPC server
func NewGRPCServer(healthGrpcHandler *healthDelivery.GRPCHandler, memberGrpcHandler *memberDelivery.GRPCHandler) *Server {
	return &Server{
		healthGRPCHandler: healthGrpcHandler,
		memberGRPCHandler: memberGrpcHandler,
	}
}

// Serve insecure server/ no server side encryption
func (s *Server) Serve(port uint) error {
	address := fmt.Sprintf(":%d", port)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		//Unary interceptor
		grpc.UnaryInterceptor(middleware.Auth),
	)

	//Register all sub server here
	healthPB.RegisterPingPongServiceServer(server, s.healthGRPCHandler)
	memberPB.RegisterMemberServiceServer(server, s.memberGRPCHandler)
	//end register server

	err = server.Serve(l)

	if err != nil {
		return err
	}

	return nil
}
