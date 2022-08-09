package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/spf13/cast"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

// MerchantService data structure
type MerchantService struct {
	BaseURL         *url.URL
	BaseGraphQLURL  *url.URL
	Qpublisher      QPublisher
	ActivityService ActivityServices
}

// NewMerchantService function for initializing static service
func NewMerchantService(
	qPublisher QPublisher,
	activityService ActivityServices,
) (*MerchantService, error) {
	var (
		merchant MerchantService
		err      error
		ok       bool
	)

	baseURL, ok := os.LookupEnv("MERCHANT_SERVICE_URL")
	if !ok {
		return &merchant, errors.New("you need to specify MERCHANT_SERVICE_URL in the environment variable")
	}

	BaseGraphQLURL, ok := os.LookupEnv("MERCHANT_SERVICE_GRAPHQL_URL")
	if !ok {
		return &merchant, errors.New("you need to specify MERCHANT_SERVICE_GRAPHQL_URL in the environment variable")
	}

	merchant.BaseURL, err = url.Parse(baseURL)
	if err != nil {
		return &merchant, errors.New("error parsing merchant services url")
	}

	merchant.BaseGraphQLURL, err = url.Parse(BaseGraphQLURL)
	if err != nil {
		return &merchant, errors.New("error parsing merchant services graphql url")
	}

	merchant.Qpublisher = qPublisher
	merchant.ActivityService = activityService

	return &merchant, nil
}

// FindMerchantServiceByID function for getting detail static by id
func (m *MerchantService) FindMerchantServiceByID(ctxReq context.Context, userID, token, merchantID string) <-chan serviceModel.ServiceResult {
	ctx := "MerchantUseCase-FindMerchantServiceByID"
	var result = make(chan serviceModel.ServiceResult)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		// request
		var (
			response interface{}
			err      error
		)
		if os.Getenv("GWS_MERCHANT_ACTIVE") == "true" {
			response, err = m.GetMerchantServiceGraphQL(ctxReq, token, merchantID)
		} else {
			response, err = m.GetMerchantService(ctxReq, token)
		}

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_merchant_by_id", err, []string{token, merchantID})
			result <- serviceModel.ServiceResult{Error: errors.New(model.FailedGetMerchant)}
			return
		}

		tags[helper.TextResponse] = response

		result <- serviceModel.ServiceResult{Result: response}
	})

	return result
}

// GetMerchantService function for getting from merchant service
func (m *MerchantService) GetMerchantService(ctxReq context.Context, token string) (interface{}, error) {
	ctx := "MerchantService-GetMerchantService"

	// generate uri
	uri := fmt.Sprintf("%s%s", m.BaseURL.String(), "/api/v1/merchant/me")
	trace := tracer.StartTrace(ctxReq, ctx)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": token,
	}

	defer trace.Finish(map[string]interface{}{
		"Auth": token,
	})
	resp := serviceModel.ResponseMerchantService{}

	if err := helper.GetHTTPNewRequestV2(ctxReq, http.MethodGet, uri, nil, &resp, headers); err != nil {
		helper.SendErrorLog(ctxReq, ctx, "request_to_merchant", err, headers)
		return "", errors.New("failed get merchant")
	}

	return resp, nil
}

// RestructMerchant function for restruct from cdc
func RestructMerchant(pl serviceModel.MerchantPayloadCDC) model.B2CMerchant {
	var merchantModel model.B2CMerchant
	merchantModel.ID = pl.Payload.After.ID
	merchantModel.UserID = pl.Payload.After.UserID
	merchantModel.MerchantName = pl.Payload.After.MerchantName
	merchantModel.VanityURL = pl.Payload.After.VanityURL
	merchantModel.MerchantCategory = pl.Payload.After.MerchantCategory
	merchantModel.CompanyName = pl.Payload.After.CompanyName
	merchantModel.Pic = pl.Payload.After.Pic
	merchantModel.PicOccupation = pl.Payload.After.PicOccupation
	merchantModel.DailyOperationalStaff = pl.Payload.After.DailyOperationalStaff
	merchantModel.StoreClosureDate = pl.Payload.After.StoreClosureDate
	merchantModel.StoreReopenDate = pl.Payload.After.StoreReopenDate
	merchantModel.StoreActiveShippingDate = pl.Payload.After.StoreActiveShippingDate
	merchantModel.MerchantAddress = pl.Payload.After.MerchantAddress
	merchantModel.MerchantVillage = pl.Payload.After.MerchantVillage
	merchantModel.MerchantDistrict = pl.Payload.After.MerchantDistrict
	merchantModel.MerchantCity = pl.Payload.After.MerchantCity
	merchantModel.MerchantProvince = pl.Payload.After.MerchantProvince
	merchantModel.ZipCode = pl.Payload.After.ZipCode
	merchantModel.StoreAddress = pl.Payload.After.StoreAddress
	merchantModel.StoreVillage = pl.Payload.After.StoreVillage
	merchantModel.StoreDistrict = pl.Payload.After.StoreDistrict
	merchantModel.StoreCity = pl.Payload.After.StoreCity
	merchantModel.StoreProvince = pl.Payload.After.StoreProvince
	merchantModel.StoreZipCode = pl.Payload.After.StoreZipCode
	merchantModel.PhoneNumber = pl.Payload.After.PhoneNumber
	merchantModel.MobilePhoneNumber = pl.Payload.After.MobilePhoneNumber
	merchantModel.AdditionalEmail = pl.Payload.After.AdditionalEmail
	merchantModel.MerchantDescription = pl.Payload.After.MerchantDescription
	merchantModel.MerchantLogo = pl.Payload.After.MerchantLogo
	merchantModel.AccountHolderName = pl.Payload.After.AccountHolderName
	merchantModel.BankName = pl.Payload.After.BankName
	merchantModel.AccountNumber = pl.Payload.After.AccountNumber
	merchantModel.IsPKP = pl.Payload.After.IsPKP
	merchantModel.Npwp = pl.Payload.After.Npwp
	merchantModel.NpwpHolderName = pl.Payload.After.NpwpHolderName
	merchantModel.RichContent = pl.Payload.After.RichContent
	merchantModel.NotificationPreferences = pl.Payload.After.NotificationPreferences
	merchantModel.MerchantRank = pl.Payload.After.MerchantRank
	merchantModel.Acquisitor = pl.Payload.After.Acquisitor
	merchantModel.AccountManager = pl.Payload.After.AccountManager
	merchantModel.LaunchDev = pl.Payload.After.LaunchDev
	merchantModel.SkuLive = pl.Payload.After.SkuLive
	merchantModel.MouDate = pl.Payload.After.MouDate
	merchantModel.Note = pl.Payload.After.Note
	merchantModel.AgreementDate = pl.Payload.After.AgreementDate
	merchantModel.IsActive = pl.Payload.After.IsActive
	merchantModel.CreatorID = pl.Payload.After.CreatorID
	merchantModel.CreatorIP = pl.Payload.After.CreatorIP
	merchantModel.EditorID = pl.Payload.After.EditorID
	merchantModel.EditorIP = pl.Payload.After.EditorIP
	merchantModel.Version = pl.Payload.After.Version
	merchantModel.Created = pl.Payload.After.Created
	merchantModel.LastModified = pl.Payload.After.LastModified
	merchantModel.MerchantVillageID = pl.Payload.After.MerchantVillageID
	merchantModel.MerchantDistrictID = pl.Payload.After.MerchantDistrictID
	merchantModel.MerchantCityID = pl.Payload.After.MerchantCityID
	merchantModel.MerchantProvinceID = pl.Payload.After.MerchantProvinceID
	merchantModel.StoreVillageID = pl.Payload.After.StoreVillageID
	merchantModel.StoreDistrictID = pl.Payload.After.StoreDistrictID
	merchantModel.StoreCityID = pl.Payload.After.StoreCityID
	merchantModel.StoreProvinceID = pl.Payload.After.StoreProvinceID
	merchantModel.BankID = pl.Payload.After.BankID
	merchantModel.IsClosed = pl.Payload.After.IsClosed
	merchantModel.MerchantEmail = pl.Payload.After.MerchantEmail
	merchantModel.BankBranch = pl.Payload.After.BankBranch
	merchantModel.PicKtpFile = pl.Payload.After.PicKtpFile
	merchantModel.NpwpFile = pl.Payload.After.NpwpFile
	merchantModel.Source = pl.Payload.After.Source
	merchantModel.BusinessType = pl.Payload.After.BusinessType
	return merchantModel
}

// RestructMerchantDocument function for restruct from cdc
func RestructMerchantDocument(pl serviceModel.MerchantDocumentPayloadCDC) model.B2CMerchantDocument {
	var merchantDocumentModel model.B2CMerchantDocument
	merchantDocumentModel.ID = pl.Payload.After.ID
	merchantDocumentModel.MerchantID = pl.Payload.After.MerchantID
	merchantDocumentModel.DocumentType = pl.Payload.After.DocumentType
	merchantDocumentModel.DocumentValue = pl.Payload.After.DocumentValue
	merchantDocumentModel.DocumentExpirationDate = pl.Payload.After.DocumentExpirationDate
	merchantDocumentModel.CreatorID = pl.Payload.After.CreatorID
	merchantDocumentModel.CreatorIP = pl.Payload.After.CreatorIP
	merchantDocumentModel.EditorID = pl.Payload.After.EditorID
	merchantDocumentModel.EditorIP = pl.Payload.After.EditorIP
	merchantDocumentModel.Version = pl.Payload.After.Version
	merchantDocumentModel.Created = pl.Payload.After.Created
	merchantDocumentModel.LastModified = pl.Payload.After.LastModified
	return merchantDocumentModel
}

// RestructMerchantBank function for restruct from cdc
func RestructMerchantBank(pl serviceModel.MerchantBankPayloadCDC) model.B2CMerchantBank {
	var merchantBankModel model.B2CMerchantBank
	merchantBankModel.ID = pl.Payload.After.ID
	merchantBankModel.BankCode = pl.Payload.After.BankCode
	merchantBankModel.BankName = pl.Payload.After.BankName
	merchantBankModel.Status = pl.Payload.After.Status
	merchantBankModel.CreatorID = pl.Payload.After.CreatorID
	merchantBankModel.CreatorIP = pl.Payload.After.CreatorIP
	merchantBankModel.Created = pl.Payload.After.Created
	merchantBankModel.EditorID = pl.Payload.After.EditorID
	merchantBankModel.EditorIP = pl.Payload.After.EditorIP
	merchantBankModel.LastModified = pl.Payload.After.LastModified
	return merchantBankModel
}

// RestructMerchantGWS function for restruct from GWS
func RestructMerchantGWS(pl serviceModel.GWSMerchantData) model.B2CMerchant {
	var merchantModel model.B2CMerchant

	businessType := mapBusinessType(pl.BusinessType)
	merchantModel.ID = pl.Code
	merchantModel.UserID = &pl.UserID
	merchantModel.MerchantName = &pl.Name
	merchantModel.VanityURL = &pl.VanityURL
	merchantModel.CompanyName = &pl.CompanyName
	merchantModel.Pic = &pl.PicName
	merchantModel.PicOccupation = &pl.PicOccupation
	merchantModel.DailyOperationalStaff = &pl.DailyOperationalStaff
	merchantModel.StoreClosureDate = &pl.StoreClosureDate
	merchantModel.StoreReopenDate = &pl.StoreReopenDate
	merchantModel.StoreActiveShippingDate = &pl.StoreActiveShippingDate
	merchantModel.MerchantAddress = &pl.Address
	merchantModel.MerchantVillage = &pl.SubDistrict.Name
	merchantModel.MerchantDistrict = &pl.District.Name
	merchantModel.MerchantCity = &pl.City.Name
	merchantModel.MerchantProvince = &pl.Province.Name
	merchantModel.ZipCode = &pl.ZipCode
	merchantModel.StoreAddress = &pl.StoreAddress
	merchantModel.StoreVillage = &pl.StoreSubDistrict.Name
	merchantModel.StoreDistrict = &pl.StoreDistrict.Name
	merchantModel.StoreCity = &pl.StoreCity.Name
	merchantModel.StoreProvince = &pl.StoreProvince.Name
	merchantModel.StoreZipCode = &pl.StoreZipCode
	merchantModel.PhoneNumber = &pl.PhoneNumber
	merchantModel.MobilePhoneNumber = &pl.MobileNumber
	merchantModel.AdditionalEmail = &pl.AdditionalEmail
	merchantModel.MerchantDescription = &pl.Description
	merchantModel.MerchantLogo = &pl.Logo
	merchantModel.AccountHolderName = &pl.Bank.AccountName
	merchantModel.BankName = &pl.Bank.Name
	merchantModel.AccountNumber = &pl.Bank.AccountNo
	merchantModel.BankBranch = &pl.Bank.Branch
	merchantModel.IsPKP = pl.IsPKP
	merchantModel.Npwp = &pl.NpwpNo
	merchantModel.NpwpHolderName = &pl.NpwpName
	merchantModel.RichContent = &pl.RichContent
	merchantModel.MerchantRank = &pl.MerchantRank
	merchantModel.Acquisitor = &pl.Acquisitor
	merchantModel.AccountManager = &pl.AccountManager
	merchantModel.LaunchDev = &pl.LaunchDev
	merchantModel.MouDate = &pl.MouDate
	merchantModel.Note = &pl.Note
	merchantModel.AgreementDate = &pl.AgreementDate
	merchantModel.IsActive = pl.IsActive
	merchantModel.CreatorID = &pl.CreatedBy
	merchantModel.CreatorIP = &pl.CreatedIP
	merchantModel.EditorID = &pl.UpdatedBy
	merchantModel.EditorIP = &pl.UpdatedIP
	merchantModel.Version = &pl.Version
	merchantModel.Created = &pl.CreatedAt
	merchantModel.DeletedAt = &pl.DeletedAt
	merchantModel.MerchantVillageID = &pl.SubDistrict.ID
	merchantModel.MerchantDistrictID = &pl.District.ID
	merchantModel.MerchantCityID = &pl.City.ID
	merchantModel.MerchantProvinceID = &pl.Province.ID
	merchantModel.StoreVillageID = &pl.StoreSubDistrict.ID
	merchantModel.StoreDistrictID = &pl.StoreDistrict.ID
	merchantModel.StoreCityID = &pl.StoreCity.ID
	merchantModel.StoreProvinceID = &pl.StoreProvince.ID
	merchantModel.IsClosed = pl.StoreIsClosed
	merchantModel.MerchantEmail = &pl.UserEmail
	merchantModel.BusinessType = &businessType
	merchantModel.MerchantType = &pl.MerchantTypeString
	merchantModel.GenderPic = &pl.GenderPicString
	merchantModel.MerchantGroup = &pl.MerchantGroup
	merchantModel.UpgradeStatus = &pl.UpgradeStatus
	merchantModel.ProductType = &pl.ProductType

	if pl.UpdatedAt != "" {
		m, err := time.Parse(time.RFC3339, pl.UpdatedAt)
		if err != nil {
			pl.UpdatedAt = ""
		}
		pl.UpdatedAt = m.Format(time.RFC3339)
	}

	merchantModel.LastModified = &pl.UpdatedAt
	var legalEntity, numOfEmployee interface{}
	if pl.LegalEntity.Valid {
		m := int(pl.LegalEntity.Int64)
		legalEntity = m
	}
	if pl.NumberOfEmployee.Valid {
		n := int(pl.NumberOfEmployee.Int64)
		numOfEmployee = n
	}
	le := cast.ToInt(legalEntity)
	noe := cast.ToInt(numOfEmployee)

	merchantModel.LegalEntity = &le
	merchantModel.NumberOfEmployee = &noe

	source := "GWS"
	merchantModel.Source = &source
	bankID := cast.ToInt64(pl.Bank.BankID)
	merchantModel.BankID = &bankID

	return merchantModel
}

// RestructMerchantDocumentGWS function for restruct from GWS
func RestructMerchantDocumentGWS(pl serviceModel.GWSMerchantDocumentData, merchantData serviceModel.GWSMerchantData) model.B2CMerchantDocument {
	var merchantDocumentModel model.B2CMerchantDocument

	documentType := mapDocumentType(pl.Type)
	merchantDocumentModel.ID = pl.Code
	merchantDocumentModel.MerchantID = &merchantData.Code
	merchantDocumentModel.DocumentType = &documentType
	merchantDocumentModel.DocumentValue = &pl.Value
	merchantDocumentModel.CreatorID = &merchantData.CreatedBy
	merchantDocumentModel.CreatorIP = &merchantData.CreatedIP
	merchantDocumentModel.EditorID = &merchantData.UpdatedBy
	merchantDocumentModel.EditorIP = &merchantData.UpdatedIP
	merchantDocumentModel.Created = &merchantData.CreatedAt
	merchantDocumentModel.LastModified = &merchantData.UpdatedAt
	merchantDocumentModel.DocumentExpirationDate = &pl.Expired

	return merchantDocumentModel
}

// RestructMerchantDataGWS function for restruct from GWS
func RestructMerchantDataGWS(pl serviceModel.GWSMerchantData) model.B2CMerchantDataV2 {
	var merchantModel model.B2CMerchantDataV2

	if pl.StoreClosureDate != "" {
		storeClosureDate, _ := time.Parse(time.RFC3339, pl.StoreClosureDate)
		merchantModel.StoreClosureDate = &storeClosureDate
	}

	if pl.StoreReopenDate != "" {
		storeReopenDate, _ := time.Parse(time.RFC3339, pl.StoreReopenDate)
		merchantModel.StoreReopenDate = &storeReopenDate
	}

	if pl.StoreActiveShippingDate != "" {
		storeActiveShippingDate, _ := time.Parse(time.RFC3339, pl.StoreActiveShippingDate)
		merchantModel.StoreActiveShippingDate = &storeActiveShippingDate
	}

	if pl.DeletedAt != "" {
		deletedAt, _ := time.Parse(time.RFC3339, pl.DeletedAt)
		merchantModel.DeletedAt = null.TimeFrom(deletedAt)
	}

	if pl.MouDate != "" {
		mouDate, _ := time.Parse(time.RFC3339, pl.MouDate)
		merchantModel.MouDate = null.TimeFrom(mouDate)
	}

	if pl.AgreementDate != "" {
		agreementDate, _ := time.Parse(time.RFC3339, pl.AgreementDate)
		merchantModel.AgreementDate = null.TimeFrom(agreementDate)
	}

	if pl.CreatedAt != "" {
		createdAt, _ := time.Parse(time.RFC3339, pl.CreatedAt)
		merchantModel.Created = zero.TimeFrom(createdAt)
	}

	if pl.UpdatedAt != "" {
		updatedAt, _ := time.Parse(time.RFC3339, pl.UpdatedAt)
		merchantModel.LastModified = null.TimeFrom(updatedAt)
	}

	Version := cast.ToInt64(pl.Version)

	businessType := mapBusinessType(pl.BusinessType)
	merchantModel.ID = pl.Code
	merchantModel.UserID = pl.UserID
	merchantModel.MerchantName = pl.Name
	merchantModel.VanityURL = zero.StringFrom(pl.VanityURL)
	merchantModel.CompanyName = zero.StringFrom(pl.CompanyName)
	merchantModel.Pic = zero.StringFrom(pl.PicName)
	merchantModel.PicOccupation = zero.StringFrom(pl.PicOccupation)
	merchantModel.DailyOperationalStaff = zero.StringFrom(pl.DailyOperationalStaff)
	merchantModel.MerchantAddress = zero.StringFrom(pl.Address)
	merchantModel.MerchantVillage = zero.StringFrom(pl.SubDistrict.Name)
	merchantModel.MerchantDistrict = zero.StringFrom(pl.District.Name)
	merchantModel.MerchantCity = zero.StringFrom(pl.City.Name)
	merchantModel.MerchantProvince = zero.StringFrom(pl.Province.Name)
	merchantModel.ZipCode = zero.StringFrom(pl.ZipCode)
	merchantModel.StoreAddress = zero.StringFrom(pl.StoreAddress)
	merchantModel.StoreVillage = zero.StringFrom(pl.StoreSubDistrict.Name)
	merchantModel.StoreDistrict = zero.StringFrom(pl.StoreDistrict.Name)
	merchantModel.StoreCity = zero.StringFrom(pl.StoreCity.Name)
	merchantModel.StoreProvince = zero.StringFrom(pl.StoreProvince.Name)
	merchantModel.StoreZipCode = zero.StringFrom(pl.StoreZipCode)
	merchantModel.PhoneNumber = zero.StringFrom(pl.PhoneNumber)
	merchantModel.MobilePhoneNumber = zero.StringFrom(pl.MobileNumber)
	merchantModel.AdditionalEmail = zero.StringFrom(pl.AdditionalEmail)
	merchantModel.MerchantDescription = zero.StringFrom(pl.Description)
	merchantModel.MerchantLogo = zero.StringFrom(pl.Logo)
	merchantModel.AccountHolderName = zero.StringFrom(pl.Bank.AccountName)
	merchantModel.BankCode = zero.StringFrom(pl.Bank.BankCode)
	merchantModel.BankName = zero.StringFrom(pl.Bank.Name)
	merchantModel.AccountNumber = zero.StringFrom(pl.Bank.AccountNo)
	merchantModel.BankBranch = zero.StringFrom(pl.Bank.Branch)
	merchantModel.IsPKP = pl.IsPKP
	merchantModel.Npwp = zero.StringFrom(pl.NpwpNo)
	merchantModel.NpwpHolderName = zero.StringFrom(pl.NpwpName)
	merchantModel.RichContent = zero.StringFrom(pl.RichContent)
	merchantModel.MerchantRank = zero.StringFrom(pl.MerchantRank)
	merchantModel.Acquisitor = zero.StringFrom(pl.Acquisitor)
	merchantModel.AccountManager = zero.StringFrom(pl.AccountManager)
	merchantModel.LaunchDev = zero.StringFrom(pl.LaunchDev)
	merchantModel.Note = zero.StringFrom(pl.Note)
	merchantModel.IsActive = pl.IsActive
	merchantModel.CreatorID = null.StringFrom(pl.CreatedBy)
	merchantModel.CreatorIP = null.StringFrom(pl.CreatedIP)
	merchantModel.EditorID = null.StringFrom(pl.UpdatedBy)
	merchantModel.EditorIP = null.StringFrom(pl.UpdatedIP)
	merchantModel.Version.Int64 = Version
	merchantModel.MerchantVillageID = zero.StringFrom(pl.SubDistrict.ID)
	merchantModel.MerchantDistrictID = zero.StringFrom(pl.District.ID)
	merchantModel.MerchantCityID = zero.StringFrom(pl.City.ID)
	merchantModel.MerchantProvinceID = zero.StringFrom(pl.Province.ID)
	merchantModel.StoreVillageID = zero.StringFrom(pl.StoreSubDistrict.ID)
	merchantModel.StoreDistrictID = zero.StringFrom(pl.StoreDistrict.ID)
	merchantModel.StoreCityID = zero.StringFrom(pl.StoreCity.ID)
	merchantModel.StoreProvinceID = zero.StringFrom(pl.StoreProvince.ID)
	merchantModel.IsClosed = null.BoolFrom(pl.StoreIsClosed)
	merchantModel.MerchantEmail = zero.StringFrom(pl.UserEmail)
	merchantModel.BusinessType = zero.StringFrom(businessType)
	merchantModel.GenderPic = pl.GenderPic
	merchantModel.MerchantGroup = zero.StringFrom(pl.MerchantGroup)
	merchantModel.UpgradeStatus = zero.StringFrom(pl.UpgradeStatus)
	merchantModel.MerchantTypeString = zero.StringFrom(pl.MerchantTypeString)
	merchantModel.GenderPicString = zero.StringFrom(pl.GenderPicString)
	merchantModel.ProductType = zero.StringFrom(pl.ProductType)

	source := "GWS"
	merchantModel.Source = zero.StringFrom(source)
	bankID := cast.ToInt64(pl.Bank.BankID)
	merchantModel.BankID = zero.IntFrom(bankID)

	return merchantModel
}

// RestructMerchantDocumentDataGWS function for restruct from GWS
func RestructMerchantDocumentDataGWS(pl serviceModel.GWSMerchantDocumentData, merchantData serviceModel.GWSMerchantData) model.B2CMerchantDocumentData {
	var merchantDocumentModel model.B2CMerchantDocumentData

	documentType := mapDocumentType(pl.Type)
	merchantDocumentModel.ID = pl.Code
	merchantDocumentModel.MerchantID = merchantData.Code
	merchantDocumentModel.DocumentType = documentType
	merchantDocumentModel.DocumentValue = pl.Value
	merchantDocumentModel.CreatorID = merchantData.CreatedBy
	merchantDocumentModel.CreatorIP = merchantData.CreatedIP
	merchantDocumentModel.EditorID = merchantData.UpdatedBy
	merchantDocumentModel.EditorIP = merchantData.UpdatedIP

	if pl.Expired != "" {
		expiredDate, _ := time.Parse(time.RFC3339, pl.Expired)
		merchantDocumentModel.DocumentExpirationDate = null.TimeFromPtr(&expiredDate)
	}

	if merchantData.CreatedAt != "" {
		createdAt, _ := time.Parse(time.RFC3339, merchantData.CreatedAt)
		merchantDocumentModel.Created = null.TimeFrom(createdAt)
	}

	if merchantData.UpdatedAt != "" {
		updatedAt, _ := time.Parse(time.RFC3339, merchantData.UpdatedAt)
		merchantDocumentModel.LastModified = null.TimeFrom(updatedAt)
	}

	return merchantDocumentModel
}

func mapDocumentType(key string) string {
	var documentType = map[string]string{
		"AKTA_PENDIRIAN":                          "AktaPendirian-file",
		"AKTA_PERUBAHAN_TERAKHIR":                 "AktaPerubahan-file",
		"NOMOR_IZIN_USAHA":                        "NIB-file",
		"SIUP":                                    "SIUP-file",
		"TANDA_DAFTAR_PERUSAHAAN":                 "TDP-file",
		"SURAT_KETERANGAN_TERDAFTAR_PAJAK":        "SKTP-file",
		"SURAT_PENGUKUHAN_PENGUSAHA_KENA_PAJAK":   "SPPKP-file",
		"LETTER_OF_APPOINTMENT_FROM_DISTRIBUTION": "LoA-file",
		"LOCATION_SERVICE_CENTER":                 "ServiceCenter-file",
		"SURAT_IJIN":                              "SuratIjin-file",
		"SERTIFIKAT_MEREK":                        "SertifikatMerek-file",
		"SURAT_KETERANGAN_DOMISILI_PERUSAHAAN":    "SKDP-file",
		"SURAT_KETERANGAN_KEMENKUMHAM":            "SKKEMENKUMHAM-file",
		"SERTIFIKAT_KEAHLIAN":                     "SertifikatKeahlian-file",
		"SURAT_PERNYATAAN_UMK":                    "SUMK-file",
	}
	return documentType[key]
}

func mapBusinessType(key string) string {
	var businessType = map[string]string{
		"PERSONAL":  "perorangan",
		"CORPORATE": "perusahaan",
	}
	return businessType[key]
}

// RestructMasterMerchantBankGWS function for restruct from GWS
func RestructMasterMerchantBankGWS(pl serviceModel.GWSMasterMerchantBank) model.B2CMerchantBank {
	var merchantBankModel model.B2CMerchantBank
	merchantBankModel.ID = pl.BankID
	merchantBankModel.BankCode = &pl.Code
	merchantBankModel.BankName = &pl.Name
	merchantBankModel.Status = pl.IsActive
	merchantBankModel.Created = &pl.CreatedAt
	merchantBankModel.LastModified = &pl.UpdatedAt
	merchantBankModel.DeletedAt = &pl.DeletedAt
	merchantBankModel.CreatorID = &pl.CreatedBy
	merchantBankModel.EditorID = &pl.UpdatedBy
	merchantBankModel.CreatorIP = &pl.CreatedIP
	merchantBankModel.EditorIP = &pl.UpdatedIP
	return merchantBankModel
}

// PublishToKafkaUserMerchant function to publish merchant data
func (m *MerchantService) PublishToKafkaUserMerchant(ctxReq context.Context, data *model.B2CMerchantDataV2, eventType, producer string) error {
	ctx := "MerchantService-PublishToKafkaUserMerchant"
	kafkaUserServiceMerchantTopic := golib.GetEnvOrFail(ctx, helper.TextFindServerKafkaConfig, "KAFKA_USER_SERVICE_MERCHANT_TOPIC")

	KafkaPayload := serviceModel.MerchantPayloadKafka{
		EventOrchestration:     "UpsertMerchant",
		TimestampOrchestration: time.Now().Format(time.RFC3339),
		EventType:              eventType,
		Counter:                0,
		Producer:               producer,
		Payload:                data,
	}

	// prepare to send to kafka
	payloadJSON, _ := json.Marshal(KafkaPayload)
	messageKey := data.ID

	if err := m.Qpublisher.PublishKafka(ctxReq, kafkaUserServiceMerchantTopic, messageKey, payloadJSON); err != nil {
		helper.SendErrorLog(ctxReq, ctx, "publish_payload", err, KafkaPayload)
		return err
	}

	return nil
}

// InsertLogMerchant function to write log activity service for merchant
func (m *MerchantService) InsertLogMerchant(ctxReq context.Context, oldData, newData model.B2CMerchantDataV2, action, module string) error {
	targetID := newData.ID
	if targetID == "" {
		targetID = oldData.ID
	}

	payload := serviceModel.Payload{
		Module:    module,
		Action:    action,
		Target:    targetID,
		CreatorID: newData.CreatorID.String,
		EditorID:  newData.EditorID.String,
	}

	m.ActivityService.InsertLog(ctxReq, oldData, newData, payload)

	return nil
}

func (m *MerchantService) InsertLogMerchantPIC(ctxReq context.Context, oldData, newData model.B2CMerchantDataV2, action, module string, member memberModel.Member) error {
	targetID := newData.ID
	if targetID == "" {
		targetID = oldData.ID
	}
	var role string
	if member.IsAdmin {
		role = "admin"
	} else {
		role = "user"
	}
	User := serviceModel.Users{
		Id:       newData.EditorID.String,
		FullName: member.FirstName + " " + member.LastName,
		Role:     role,
		Email:    member.Email,
	}
	text := []string{"private", "internal"}
	payload := serviceModel.Payload{
		Module:    module,
		Action:    action,
		Target:    targetID,
		CreatorID: newData.CreatorID.String,
		EditorID:  newData.EditorID.String,
		User:      User,
		ViewType:  text,
	}

	m.ActivityService.InsertLog(ctxReq, oldData, newData, payload)

	return nil
}
