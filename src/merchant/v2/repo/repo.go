package repo

import (
	"context"

	"github.com/Bhinneka/user-service/src/merchant/v2/model"
)

// ResultRepository data structure
type ResultRepository struct {
	Result    interface{}
	TotalData int
	Error     error
}

// MerchantRepository interface abstraction
type MerchantRepository interface {
	Save(model.B2CMerchant) error
	SaveMerchantGWS(ctxReq context.Context, data model.B2CMerchant) error
	UpdateMerchantGWS(ctxReq context.Context, data model.B2CMerchant) error
	Delete(string) error
	AddUpdateMerchant(ctxReq context.Context, data model.B2CMerchantDataV2) <-chan ResultRepository
	LoadMerchant(ctxReq context.Context, uid string, privacy string) ResultRepository
	LoadMerchantByVanityURL(ctxReq context.Context, vanityURL string) ResultRepository
	FindMerchantByEmail(ctxReq context.Context, uid string) ResultRepository
	FindMerchantByUser(ctxReq context.Context, uid string) ResultRepository
	FindMerchantByName(ctxReq context.Context, uid string) ResultRepository
	FindMerchantBySlug(ctxReq context.Context, slug string) ResultRepository
	FindMerchantByID(ctxReq context.Context, id, uid string) ResultRepository
	SoftDelete(ctxReq context.Context, merchantID string) <-chan ResultRepository
	GetMerchants(ctxReq context.Context, params *model.QueryParameters) <-chan ResultRepository
	GetMerchantsPublic(ctxReq context.Context, params *model.QueryParameters) <-chan ResultRepository
	GetTotalMerchant(ctxReq context.Context, params *model.QueryParameters) <-chan ResultRepository
	LoadLegalEntity(ctxReq context.Context, id int) ResultRepository
	LoadCompanySize(ctxReq context.Context, id int) ResultRepository

	RejectUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2, reasonReject string) <-chan ResultRepository
	ClearRejectUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2) <-chan ResultRepository
}

// MerchantDocumentRepository interface abstraction
type MerchantDocumentRepository interface {
	Save(model.B2CMerchantDocument) error
	SaveMerchantDocumentGWS(model.B2CMerchantDocument) error
	Delete(string) error
	FindMerchantDocumentByParam(ctxReq context.Context, param *model.B2CMerchantDocumentQueryInput) <-chan ResultRepository
	GetListMerchantDocument(ctxReq context.Context, params *model.B2CMerchantDocumentQueryInput) <-chan ResultRepository
	InsertNewMerchantDocument(ctxReq context.Context, param *model.B2CMerchantDocumentData) <-chan ResultRepository
	UpdateMerchantDocument(ctxReq context.Context, id string, param *model.B2CMerchantDocumentData) <-chan ResultRepository
	ResetRejectedDocument(ctxReq context.Context, param model.B2CMerchantDocumentData) <-chan ResultRepository
}

// MerchantBankRepository interface abstraction
type MerchantBankRepository interface {
	Save(model.B2CMerchantBank) error
	Delete(model.B2CMerchantBank) error
	SaveMasterBankGWS(ctxReq context.Context, data model.B2CMerchantBank) error
	Load(uid string) ResultRepository
	FindActiveMerchantBankByID(ctxReq context.Context, bankID int) <-chan ResultRepository
	GetListMerchantBank(params *model.ParametersMerchantBank) <-chan ResultRepository
	GetTotalMerchantBank(params *model.ParametersMerchantBank) <-chan ResultRepository
}

// MerchantAddressRepository interface abstraction
type MerchantAddressRepository interface {
	CountAddress(ctxReq context.Context, relationID, relationName string, params *model.ParameterWarehouse) <-chan ResultRepository
	AddUpdateAddress(ctxReq context.Context, data model.AddressData) <-chan ResultRepository
	AddPhoneAddress(ctxReq context.Context, data model.PhoneData) <-chan ResultRepository
	UpdatePhoneAddress(ctxReq context.Context, data model.PhoneData) <-chan ResultRepository
	DeleteWarehouseAddress(ctxReq context.Context, addressID string) <-chan ResultRepository
	DeletePhoneAddress(ctxReq context.Context, relationID, relationName string) <-chan ResultRepository
	FindMerchantAddress(ctxReq context.Context, id string) <-chan ResultRepository
	UpdatePrimaryAddressByRelationID(ctxReq context.Context, relationID, relationName string) <-chan ResultRepository
	GetListAddress(ctxReq context.Context, params *model.ParameterWarehouse) <-chan ResultRepository
	AddUpdateAddressMaps(ctxReq context.Context, data model.Maps) <-chan ResultRepository
	FindAddressMaps(ctxReq context.Context, realationID, relationName string) <-chan ResultRepository
}

// MerchantEmployeeRepository interface abstraction
type MerchantEmployeeRepository interface {
	Save(model.B2CMerchantEmployee) error
	ChangeStatus(ctxReq context.Context, params model.B2CMerchantEmployee) error

	GetAllMerchantEmployees(ctxReq context.Context, params *model.QueryMerchantEmployeeParameters) <-chan ResultRepository
	GetTotalMerchantEmployees(ctxReq context.Context, params *model.QueryMerchantEmployeeParameters) <-chan ResultRepository

	GetMerchantEmployees(ctxReq context.Context, params *model.QueryMerchantEmployeeParameters) <-chan ResultRepository
}
