package usecase

import "github.com/Bhinneka/user-service/src/session/v1/model"

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
}

// SessionInfoUseCase interface abstraction
type SessionInfoUseCase interface {
	GetSessionInfoList(params *model.ParamList) <-chan ResultUseCase
}
