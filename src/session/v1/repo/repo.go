package repo

import "github.com/Bhinneka/user-service/src/session/v1/model"

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// SessionInfoRepository interface abstraction
type SessionInfoRepository interface {
	SaveSessionInfo(params *model.SessionInfoRequest) <-chan ResultRepository
}
