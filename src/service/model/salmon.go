package model

import (
	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
	"gopkg.in/guregu/null.v4"
)

// MerchantPayloadData data structure
type MerchantPayloadData struct {
	After  merchantModel.B2CMerchant `json:"after"`
	Before merchantModel.B2CMerchant `json:"before"`
	Op     string                    `json:"op"`
}

// MerchantPayloadCDC data structure
type MerchantPayloadCDC struct {
	Payload MerchantPayloadData `json:"payload"`
}

// B2CMerchantDocument data structure
type B2CMerchantDocument struct {
	ID                     string  `json:"id"`
	MerchantID             *string `json:"merchantId"`
	DocumentType           *string `json:"documentType"`
	DocumentValue          *string `json:"documentValue"`
	DocumentExpirationDate *string `json:"documentExpirationDate"`
	CreatorID              *string `json:"creatorId"`
	CreatorIP              *string `json:"creatorIp"`
	EditorID               *string `json:"editorId"`
	EditorIP               *string `json:"editorIp"`
	Version                *int    `json:"version"`
	Created                *string `json:"created"`
	LastModified           *string `json:"lastModified"`
}

// MerchantDocumentPayloadData data structure
type MerchantDocumentPayloadData struct {
	After  B2CMerchantDocument `json:"after"`
	Before B2CMerchantDocument `json:"before"`
	Op     string              `json:"op"`
}

// MerchantDocumentPayloadCDC data structure
type MerchantDocumentPayloadCDC struct {
	Payload MerchantDocumentPayloadData `json:"payload"`
}

// MerchantPayloadKafka data structure for pushing to kafka
type MerchantPayloadKafka struct {
	EventOrchestration     string                           `json:"eventOrchestration,omitempty"`
	TimestampOrchestration string                           `json:"timestampOrchestration,omitempty"`
	EventType              string                           `json:"eventType"`
	Counter                int                              `json:"counter"`
	Producer               string                           `json:"producer"`
	Payload                *merchantModel.B2CMerchantDataV2 `json:"payload"`
}

// MerchantBankPayloadData data structure
type MerchantBankPayloadData struct {
	After  merchantModel.B2CMerchantBank `json:"after"`
	Before merchantModel.B2CMerchantBank `json:"before"`
	Op     string                        `json:"op"`
}

// MerchantBankPayloadCDC data structure
type MerchantBankPayloadCDC struct {
	Payload MerchantBankPayloadData `json:"payload"`
}

// ResponseMerchantService data structure
type ResponseMerchantService struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    merchantModel.B2CMerchant `json:"data"`
}

// ResponseGWSMerchant data structure
type ResponseGWSMerchant struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    MerchantDetailGWS `json:"data"`
}

// MerchantDetailGWS data structure
type MerchantDetailGWS struct {
	GetMerchantDetail Merchantdetail `json:"getMerchantDetail"`
}

// Merchantdetail data structure
type Merchantdetail struct {
	Code    int  `json:"code"`
	Success bool `json:"success"`
}

// GWSMerchantPayloadMessage data structure
type GWSMerchantPayloadMessage struct {
	EventType string          `json:"eventType"`
	Data      GWSMerchantData `json:"data"`
	Producer  string          `json:"producer"`
}

// GWSMasterMerchantBankPayloadMessage data structure
type GWSMasterMerchantBankPayloadMessage struct {
	EventType string                `json:"eventType"`
	Data      GWSMasterMerchantBank `json:"data"`
}

// GWSMerchantData data structure
type GWSMerchantData struct {
	ID                      string                     `json:"id"`
	Code                    string                     `json:"code"`
	UserID                  string                     `json:"userId"`
	UserEmail               string                     `json:"userEmail"`
	Name                    string                     `json:"name"`
	VanityURL               string                     `json:"vanityURL"`
	BusinessType            string                     `json:"businessType"`
	Description             string                     `json:"Description"`
	CompanyName             string                     `json:"companyName"`
	StoreAddress            string                     `json:"storeAddress"`
	StoreProvince           GWSMerchantAddressData     `json:"storeProvince"`
	StoreCity               GWSMerchantAddressData     `json:"storeCity"`
	StoreDistrict           GWSMerchantAddressData     `json:"storeDistrict"`
	StoreSubDistrict        GWSMerchantAddressData     `json:"storeSubDistrict"`
	StoreZipCode            string                     `json:"storeZipCode"`
	StoreIsClosed           bool                       `json:"storeIsClosed"`
	StoreClosureDate        string                     `json:"storeClosureDate"`
	StoreReopenDate         string                     `json:"storeReopenDate"`
	StoreActiveShippingDate string                     `json:"storeActiveShippingDate"`
	Logo                    string                     `json:"logo"`
	Address                 string                     `json:"address"`
	Province                GWSMerchantAddressData     `json:"province"`
	City                    GWSMerchantAddressData     `json:"city"`
	District                GWSMerchantAddressData     `json:"district"`
	SubDistrict             GWSMerchantAddressData     `json:"subDistrict"`
	ZipCode                 string                     `json:"zipCode"`
	IsPKP                   bool                       `json:"isPKP"`
	MerchantRank            string                     `json:"merchantRank"`
	LaunchDev               string                     `json:"launchDev"`
	MouDate                 string                     `json:"mouDate"`
	Note                    string                     `json:"note"`
	AgreementDate           string                     `json:"agreementDate"`
	PicName                 string                     `json:"picName"`
	PicOccupation           string                     `json:"picOccupation"`
	AccountManager          string                     `json:"accountManager"`
	Acquisitor              string                     `json:"acquisitor"`
	DailyOperationalStaff   string                     `json:"dailyOperationalStaff"`
	PhoneNumber             string                     `json:"phoneNumber"`
	MobileNumber            string                     `json:"mobileNumber"`
	AdditionalEmail         string                     `json:"additionalEmail"`
	RichContent             string                     `json:"richContent"`
	NpwpNo                  string                     `json:"npwpNo"`
	NpwpName                string                     `json:"npwpName"`
	Bank                    GWSMerchantBank            `json:"bank"`
	Documents               []GWSMerchantDocumentData  `json:"documents"`
	IsActive                bool                       `json:"isActive"`
	IsApproved              bool                       `json:"isApproved"`
	CreatedBy               string                     `json:"createdBy"`
	UpdatedBy               string                     `json:"updatedBy"`
	CreatedIP               string                     `json:"createdIp"`
	UpdatedIP               string                     `json:"updatedIp"`
	Version                 int64                      `json:"version"`
	CreatedAt               string                     `json:"createdAt"`
	UpdatedAt               string                     `json:"updatedAt"`
	DeletedAt               string                     `json:"deletedAt"`
	MerchantType            merchantModel.MerchantType `json:"-"`
	MerchantTypeString      string                     `json:"merchantType"`
	GenderPic               merchantModel.GenderPic    `json:"-"`
	GenderPicString         string                     `json:"genderPic"`
	MerchantGroup           string                     `json:"merchantGroup"`
	UpgradeStatus           string                     `json:"upgradeStatus"`
	ProductType             string                     `json:"productType"`
	LegalEntity             null.Int                   `json:"legalEntity"`
	LegalEntityName         null.String                `json:"legalEntityName"`
	NumberOfEmployee        null.Int                   `json:"numberOfEmployee"`
	NumberOfEmployeeName    null.String                `json:"numberOfEmployeeName"`
}

// GWSMerchantDocumentData data structure
type GWSMerchantDocumentData struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	Value   string `json:"value"`
	Expired string `json:"expired"`
}

// GWSMerchantAddressData data structure
type GWSMerchantAddressData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GWSMerchantBank data structure
type GWSMerchantBank struct {
	ID          string      `json:"id"`
	BankID      interface{} `json:"bankId"`
	BankCode    string      `json:"bankCode"`
	Name        string      `json:"name"`
	Branch      string      `json:"branch"`
	AccountNo   string      `json:"accountNo"`
	AccountName string      `json:"accountName"`
}

// GWSMasterMerchantBank data structure
type GWSMasterMerchantBank struct {
	ID        string `json:"id"`
	BankID    int    `json:"bankId"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"deletedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
	CreatedIP string `json:"createdIP"`
	UpdatedIP string `json:"updatedIP"`
}
