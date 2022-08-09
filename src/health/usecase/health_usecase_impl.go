package usecase

import (
	"errors"

	"github.com/Bhinneka/user-service/src/health/model"
	"github.com/Bhinneka/user-service/src/health/query"
)

// healthUseCaseImpl model
type healthUseCaseImpl struct {
	healthQuery query.HealthQuery
}

// NewHealthUseCase function for initializing health use case implementation
func NewHealthUseCase(healthQuery query.HealthQuery) HealthUseCase {
	return &healthUseCaseImpl{
		healthQuery: healthQuery,
	}
}

// Ping function for checking service
func (hu *healthUseCaseImpl) Ping() <-chan ResultUseCase {
	output := make(chan ResultUseCase)

	go func() {
		defer close(output)

		healthResult := <-hu.healthQuery.Ping()

		if healthResult.Error != nil {
			output <- ResultUseCase{Error: healthResult.Error}
			return
		}

		health, ok := healthResult.Result.(*model.Health)

		if !ok {
			output <- ResultUseCase{Error: errors.New("result is not health")}
			return
		}

		output <- ResultUseCase{Result: health}

	}()

	return output
}
