package mocks

import (
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	merchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
)

// MemberRepoResult mock member repositoty result
func MemberRepoResult(data memberRepo.ResultRepository) <-chan memberRepo.ResultRepository {
	output := make(chan memberRepo.ResultRepository)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

// MerchantRepoResult mock merchant repositoty result
func MerchantRepoResult(data merchantRepo.ResultRepository) <-chan merchantRepo.ResultRepository {
	output := make(chan merchantRepo.ResultRepository, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

// MemberQueryResult mock member query result
func MemberQueryResult(data memberQuery.ResultQuery) <-chan memberQuery.ResultQuery {
	output := make(chan memberQuery.ResultQuery)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}
