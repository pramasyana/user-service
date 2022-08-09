package repo

import (
	"context"

	"github.com/Bhinneka/user-service/src/shipping_address/v2/model"
)

const keyShippingDetail = "STG-shipping-detail"

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// ShippingAddressRepository interface abstraction
type ShippingAddressRepository interface {
	Save(model.ShippingAddress) error
	Update(model.ShippingAddress) error
	Delete(model.ShippingAddress) error
	AddShippingAddress(ctxReq context.Context, data model.ShippingAddressData) <-chan ResultRepository
	UpdateShippingAddress(ctxReq context.Context, data model.ShippingAddressData) <-chan ResultRepository
	CountShippingAddressByUserID(ctxReq context.Context, id string) <-chan ResultRepository
	FindShippingAddressByID(ctxReq context.Context, id string, memberID string) <-chan ResultRepository
	FindShippingAddressPrimaryByID(ctxReq context.Context, memberID string) <-chan ResultRepository
	DeleteShippingAddressByID(ctxReq context.Context, id string) <-chan ResultRepository
	GetListShippingAddress(ctxReq context.Context, params *model.ParametersShippingAddress) <-chan ResultRepository
	GetTotalShippingAddress(ctxReq context.Context, params *model.ParametersShippingAddress) <-chan ResultRepository
	UpdatePrimaryShippingAddressByID(ctxReq context.Context, id string) <-chan ResultRepository
}

// ShippingAddressRepositoryRedis data structure
type ShippingAddressRepositoryRedis interface {
	SaveRedisMeta(memberID string, page string, limit string, shippingList model.ListShippingAddress) <-chan error
	LoadRedisMeta(memberID string, page string, limit string) <-chan ResultRepository
	DeleteMultipleRedis(memberID string) <-chan error
}
