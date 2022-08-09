package delivery

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Bhinneka/user-service/protogo/health"
	"github.com/Bhinneka/user-service/src/health/model"
	"github.com/Bhinneka/user-service/src/health/usecase"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
)

// GRPCHandler data structure
type GRPCHandler struct {
	healthUseCase usecase.HealthUseCase
}

// NewGRPCHandler function for initializing grpc handler object
func NewGRPCHandler(healthUseCase usecase.HealthUseCase) *GRPCHandler {
	return &GRPCHandler{healthUseCase}
}

// PingPong function for implementing function from health protobuf
func (h *GRPCHandler) PingPong(c context.Context, arg *google_protobuf.Empty) (*pb.PongResponse, error) {
	result := <-h.healthUseCase.Ping()

	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	pong, ok := result.Result.(*model.Health)

	if !ok {
		err := errors.New("result is not health")
		return nil, status.Error(codes.Internal, err.Error())
	}

	pongRes := &pb.PongResponse{
		State:        int32(pong.State),
		Dependencies: pong.Dependencies,
		ErrorCount:   int32(pong.ErrorCount),
	}

	return pongRes, nil
}
