package mocks

import serviceModel "github.com/Bhinneka/user-service/src/service/model"

// ServiceResult mock merchant service result
func ServiceResult(data serviceModel.ServiceResult) <-chan serviceModel.ServiceResult {
	output := make(chan serviceModel.ServiceResult, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}
