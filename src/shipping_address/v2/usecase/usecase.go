package usecase

import (
	"context"

	"github.com/Bhinneka/user-service/src/shipping_address/v2/model"
)

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
	ErrorData  []model.ShippingAddressError
}

// ShippingAddressUseCase interface abstraction
type ShippingAddressUseCase interface {
	AddShippingAddress(ctxReq context.Context, data model.ShippingAddressData) <-chan ResultUseCase
	DeleteShippingAddressByID(ctxReq context.Context, shippingID string, memberID string) <-chan ResultUseCase
	UpdateShippingAddress(ctxReq context.Context, data model.ShippingAddressData) <-chan ResultUseCase
	GetListShippingAddress(ctxReq context.Context, params *model.ParametersShippingAddress) <-chan ResultUseCase
	GetDetailShippingAddress(ctxReq context.Context, shippingID string, memberID string) <-chan ResultUseCase
	GetAllListShippingAddress(ctxReq context.Context, params *model.ParametersShippingAddress, memberID string) <-chan ResultUseCase
	GetPrimaryShippingAddress(ctxReq context.Context, memberID string) <-chan ResultUseCase
	UpdatePrimaryShippingAddressByID(ctxReq context.Context, params model.ParamaterPrimaryShippingAddress) <-chan ResultUseCase

	// specific for logging
	InsertLogShipping(ctxReq context.Context, oldData, newData *model.ShippingAddressData, action string) error
}
