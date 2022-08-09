package model

import "time"

// B2BAccount data structure
type B2BAccount struct {
	ID                   string  `json:"id"`
	ParentID             *string `json:"parentId"`
	NavID                *string `json:"navId"`
	LegalEntity          *string `json:"legalEntity"`
	Name                 *string `json:"name"`
	IndustryID           *string `json:"industryId"`
	NumberEmployee       *string `json:"numberEmployee"`
	OfficeEmployee       *int    `json:"officeEmployee"`
	BusinessSize         *string `json:"businessSize"`
	BusinessGroup        *string `json:"businessGroup"`
	EstablishedYear      *int64  `json:"establishedYear"`
	IsDelete             bool    `json:"isDelete"`
	TermOfPayment        *string `json:"termOfPayment"`
	CustomerCategory     *string `json:"customerCategory"`
	UserID               *int    `json:"userId"`
	CreatedAt            *string `json:"createdAt"`
	ModifiedAt           *string `json:"modifiedAt"`
	CreatedBy            *int64  `json:"createdBy"`
	ModifiedBy           *int64  `json:"modifiedBy"`
	AccountGroupID       *int64  `json:"accountGroupId"`
	IsDisabled           bool    `json:"isDisabled"`
	Status               *string `json:"status"`
	PaymentMethodID      *string `json:"paymentMethodId"`
	PaymentMethodType    *string `json:"paymentMethodType"`
	SubPaymentMethodName *string `json:"subPaymentMethodName"`
	IsCf                 bool    `json:"isCf"`
	Logo                 *string `json:"logo"`
	IsParent             bool    `json:"isParent"`
	IsMicrosite          bool    `json:"isMicrosite"`
	MemberType           *string `json:"memberType"`
	ErpID                *string `json:"erpId"`
}

type B2BAccountContact struct {
	ID           int     `json:"id"`
	Status       *string `json:"status"`
	IsDelete     bool    `json:"isDelete"`
	IsDisabled   bool    `json:"isDisabled"`
	AccountID    *string `json:"accountId"`
	ContactID    *int    `json:"contactId"`
	CreatedAt    *string `json:"createdAt"`
	ModifiedAt   *string `json:"modifiedAt"`
	CreatedBy    *int64  `json:"createdBy"`
	ModifiedBy   *int64  `json:"modifiedBy"`
	IsAdmin      bool    `json:"isAdmin"`
	DepartmentID *int64  `json:"departmentId"`
}

// B2BAddress data structure
type B2BAddress struct {
	ID              int     `json:"id"`
	TypeAddress     *string `json:"typeAddress"`
	Name            *string `json:"name"`
	Street          *string `json:"street"`
	ProvinceID      *string `json:"provinceId"`
	ProvinceName    *string `json:"provinceName"`
	CityID          *string `json:"cityId"`
	CityName        *string `json:"cityName"`
	DistrictID      *string `json:"districtId"`
	DistrictName    *string `json:"districtName"`
	SubDistrictID   *string `json:"subdistrictId"`
	SubDistrictName *string `json:"subdistrictName"`
	IsPrimary       bool    `json:"isPrimary"`
	PostalCode      *string `json:"postalCode"`
	CreatedAt       *string `json:"createdAt"`
	ModifiedAt      *string `json:"modifiedAt"`
	CreatedBy       *int64  `json:"createdBy"`
	ModifiedBy      *int64  `json:"modifiedBy"`
	AccountID       *string `json:"accountId"`
	IsDisabled      bool    `json:"isDisabled"`
	IsDelete        bool    `json:"isDelete"`
}

// B2BContact data structure
type B2BContact struct {
	ID                   int               `json:"id"`
	ReferenceID          *string           `json:"referenceId"`
	FirstName            *string           `json:"firstName"`
	LastName             *string           `json:"lastName"`
	Salutation           *string           `json:"salutation"`
	JobTitle             *string           `json:"jobTitle"`
	Email                *string           `json:"email"`
	NavContactID         *string           `json:"navContactId"`
	IsPrimary            bool              `json:"isPrimary"`
	BirthDate            time.Time         `json:"birthDate"`
	Note                 *string           `json:"note"`
	CreatedAt            *string           `json:"createdAt"`
	ModifiedAt           *string           `json:"modifiedAt"`
	CreatedBy            *int64            `json:"createdBy"`
	ModifiedBy           *int64            `json:"modifiedBy"`
	AccountID            *string           `json:"accountId"`
	IsDisabled           bool              `json:"isDisabled"`
	Password             *string           `json:"password"`
	Avatar               *string           `json:"avatar"`
	Status               *string           `json:"status"`
	Token                *string           `json:"token"`
	IsNew                bool              `json:"isNew"`
	PhoneNumber          *string           `json:"phone_number"`
	OtherPhoneNumber     *string           `json:"other_phone_number"`
	Gender               *string           `json:"gender"`
	TransactionType      []TransactionType `json:"transaction_type"`
	ErpID                *string           `json:"erp_id"`
	IsSync               bool              `json:"is_sync"`
	Salt                 *string           `json:"salt"`
	LastPasswordModified *time.Time        `json:"last_password_modified"`
}

// B2BPhone data structure
type B2BPhone struct {
	ID           int     `json:"id"`
	RelationID   *int    `json:"relationId"`
	RelationType *string `json:"relationType"`
	TypePhone    *string `json:"typePhone"`
	Label        *string `json:"label"`
	Number       *string `json:"number"`
	Area         *string `json:"area"`
	Ext          *string `json:"ext"`
	IsPrimary    bool    `json:"isPrimary"`
	IsDelete     bool    `json:"isDelete"`
	CreatedAt    *string `json:"createdAt"`
	ModifiedAt   *string `json:"modifiedAt"`
	CreatedBy    *int64  `json:"createdBy"`
	ModifiedBy   *int64  `json:"modifiedBy"`
}

// B2BDocument data structure
type B2BDocument struct {
	ID                  int     `json:"id"`
	DocumentType        *string `json:"documentType"`
	DocumentFile        *string `json:"documentFile"`
	DocumentTitle       *string `json:"documentTitle"`
	DocumentDescription *string `json:"documentDescription"`
	NpwpNumber          *string `json:"npwpNumber"`
	NpwpName            *string `json:"npwpName"`
	NpwpAddress         *string `json:"npwpAddress"`
	SiupNumber          *string `json:"siupNumber"`
	SiupCompanyName     *string `json:"siupCompanyName"`
	SiupType            *string `json:"siupType"`
	IsDelete            bool    `json:"isDelete"`
	CreatedAt           *string `json:"createdAt"`
	ModifiedAt          *string `json:"modifiedAt"`
	CreatedBy           *int64  `json:"createdBy"`
	ModifiedBy          *int64  `json:"modifiedBy"`
	AccountID           *string `json:"accountID"`
	IsDisabled          bool    `json:"isDisabled"`
	ProvinceID          *string `json:"provinceID"`
	ProvinceName        *string `json:"provinceName"`
	CityID              *string `json:"cityID"`
	CityName            *string `json:"cityName"`
	DistrictID          *string `json:"districtID"`
	DistrictName        *string `json:"districtName"`
	SubdistrictID       *string `json:"subdistrictID"`
	SubdistrictName     *string `json:"subdistrictName"`
	PostalCode          *string `json:"postalCode"`
	Street              *string `json:"street"`
}

// B2BContactAddress data structure
type B2BContactAddress struct {
	ID              int     `json:"id"`
	Name            *string `json:"name"`
	PicName         *string `json:"picName"`
	Phone           *string `json:"phone"`
	Street          *string `json:"street"`
	ProvinceID      *string `json:"provinceId"`
	ProvinceName    *string `json:"provinceName"`
	CityID          *string `json:"cityId"`
	CityName        *string `json:"cityName"`
	DistrictID      *string `json:"districtId"`
	DistrictName    *string `json:"districtName"`
	SubDistrictID   *string `json:"subdistrictId"`
	SubDistrictName *string `json:"subdistrictName"`
	IsBilling       bool    `json:"isBilling"`
	IsShipping      bool    `json:"isShipping"`
	PostalCode      *string `json:"postalCode"`
	CreatedAt       *string `json:"createdAt"`
	ModifiedAt      *string `json:"modifiedAt"`
	CreatedBy       *int64  `json:"createdBy"`
	ModifiedBy      *int64  `json:"modifiedBy"`
	IsDelete        bool    `json:"isDelete"`
	ContactID       *int    `json:"contactId"`
}

// B2BLeads data structure
type B2BLeads struct {
	ID           int     `json:"id"`
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	SourceID     *string `json:"sourceId"`
	Source       *string `json:"source"`
	SourceStatus *string `json:"sourceStatus"`
	Status       *string `json:"status"`
	Notes        *string `json:"notes"`
	CreatedAt    *string `json:"createdAt"`
	ModifiedAt   *string `json:"modifiedAt"`
	CreatedBy    *int64  `json:"createdBy"`
	ModifiedBy   *int64  `json:"modifiedBy"`
}

type B2BContactDocument struct {
	ID                  int       `json:"id"`
	DocumentFile        string    `json:"document_file"`
	DocumentType        string    `json:"document_type"`
	DocumentTitle       string    `json:"document_title"`
	DocumentDescription string    `json:"document_description"`
	NpwpNumber          string    `json:"npwp_number"`
	NpwpName            string    `json:"npwp_name"`
	NpwpAddress         string    `json:"npwp_address"`
	NpwpAddress2        string    `json:"npwp_address2"`
	SiupNumber          string    `json:"siup_number"`
	SiupCompanyName     string    `json:"siup_company_name"`
	SiupType            string    `json:"siup_type"`
	CreatedAt           time.Time `json:"created_at"`
	ModifiedAt          time.Time `json:"modified_at"`
	CreatedBy           int       `json:"created_by"`
	ModifiedBy          int64     `json:"modified_by"`
	IsDelete            bool      `json:"is_delete"`
	ContactID           int       `json:"contact_id"`
}
