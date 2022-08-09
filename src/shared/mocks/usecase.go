package mocks

import (
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

// CreateUsecaseResult shared usecase result
func CreateUsecaseResult(data sharedModel.ResultUseCase) <-chan sharedModel.ResultUseCase {
	output := make(chan sharedModel.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}
