package usecase

import (
	"context"

	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
)

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
	ErrorData  []model.MerchantError
	TotalData  interface{}
}

// MerchantUseCase interface abstraction
type MerchantUseCase interface {
	// Self usage
	AddMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	UpgradeMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	GetListMerchantBank(params *model.ParametersMerchantBank) <-chan ResultUseCase
	CheckMerchantName(ctxReq context.Context, merchantName string) <-chan ResultUseCase
	GetMerchantByUserID(ctxReq context.Context, userID, token string) <-chan ResultUseCase
	PublishToKafkaMerchant(ctxReq context.Context, data model.B2CMerchantDataV2, eventType string) <-chan ResultUseCase
	SelfUpdateMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	SelfUpdateMerchantPartial(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	ChangeMerchantName(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	ClearRejectUpgrade(ctxReq context.Context, memberID string, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase

	// CMS Related
	UpdateMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	DeleteMerchant(ctxReq context.Context, merchantID string, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	GetMerchants(ctxReq context.Context, params *model.QueryParameters) <-chan ResultUseCase
	GetMerchantByID(ctxReq context.Context, id string, privacy string, isAttachment string) <-chan ResultUseCase
	CreateMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	RejectMerchantRegistration(ctxReq context.Context, merchantID string, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	RejectMerchantUpgrade(ctxReq context.Context, merchantID string, userAttribute *model.MerchantUserAttribute, reasonReject string) <-chan ResultUseCase
	AddMerchantPIC(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase
	GetMerchantByVanityURL(ctxReq context.Context, vanityURL string) <-chan ResultUseCase

	//Public Related
	GetMerchantsPublic(ctxReq context.Context, params *model.QueryParametersPublic) <-chan ResultUseCase

	// Queue Worker Related
	SendEmailMerchantAdd(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string) <-chan ResultUseCase
	InsertLogMerchant(ctxReq context.Context, old model.B2CMerchantDataV2, new model.B2CMerchantDataV2, action string) error
	SendEmailMerchantRejectRegistration(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string) <-chan ResultUseCase
	SendEmailMerchantRejectUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string, reasonReject string) <-chan ResultUseCase
	SendEmailAdmin(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string, reasonReject string, adminCMS string) <-chan ResultUseCase
	SendEmailActivation(ctxReq context.Context, merchant model.B2CMerchantDataV2) <-chan ResultUseCase
	SendEmailApproval(ctxReq context.Context, old model.B2CMerchantDataV2) <-chan ResultUseCase
	SendEmailMerchantUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string) <-chan ResultUseCase
	SendEmailMerchantEmployeeLogin(ctxReq context.Context, dataMerchant model.B2CMerchantDataV2, dataMember memberModel.Member) <-chan ResultUseCase
	SendEmailMerchantEmployeeRegister(ctxReq context.Context, dataMerchant model.B2CMerchantDataV2, dataMember memberModel.Member) <-chan ResultUseCase

	// merchant employee
	AddEmployee(ctxReq context.Context, token, email, firstName string) <-chan ResultUseCase
	GetAllMerchantEmployee(ctxReq context.Context, token string, params *model.QueryMerchantEmployeeParameters) <-chan ResultUseCase
	GetMerchantEmployee(ctxReq context.Context, token string, params *model.QueryMerchantEmployeeParameters) <-chan ResultUseCase
	UpdateMerchantEmployee(ctxReq context.Context, token string, params *model.QueryMerchantEmployeeParameters) <-chan ResultUseCase

	// CMS merchant employee
	CmsGetAllMerchantEmployee(ctxReq context.Context, token string, params *model.QueryCmsMerchantEmployeeParameters) <-chan ResultUseCase
}

// MerchantAddressUseCase interface abstraction
type MerchantAddressUseCase interface {
	AddUpdateWarehouseAddress(ctxReq context.Context, data model.WarehouseData, memberID, action string) <-chan ResultUseCase
	UpdatePrimaryWarehouseAddress(ctxReq context.Context, data model.ParameterPrimaryWarehouse) <-chan ResultUseCase
	GetWarehouseAddresses(ctxReq context.Context, params *model.ParameterWarehouse) <-chan ResultUseCase
	GetDetailWarehouseAddress(ctxReq context.Context, addressID, memberID string) <-chan ResultUseCase
	GetWarehouseAddressByID(ctxReq context.Context, merchantID, addressID string) <-chan ResultUseCase
	DeleteWarehouseAddress(ctxReq context.Context, addressID, memberID string) <-chan ResultUseCase
}
