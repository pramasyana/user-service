package usecase

import (
	"context"

	"github.com/Bhinneka/user-service/src/phone_area/v1/model"
)

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
	ErrorData  []model.PhoneAreaError
}

// PhoneAreaUseCase interface abstraction
type PhoneAreaUseCase interface {
	GetAllPhoneArea(ctxReq context.Context) <-chan ResultUseCase
	GetTotalPhoneArea(ctxReq context.Context) <-chan ResultUseCase
}
