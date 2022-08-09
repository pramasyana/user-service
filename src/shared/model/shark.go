package model

import (
	"time"
)

const (

	// StatusNeedReviewAccount for user need review
	StatusNeedReviewAccount = "NEED_REVIEW"

	// StatusDeactiveAccount for user deactive
	StatusDeactiveAccount = "DEACTIVATED"
)

// specific for Kafka Payload
type (

	// B2BAccountCDC data structure
	B2BAccountCDC struct {
		ID                   string  `json:"id"`
		ParentID             *string `json:"parent_id"`
		NavID                *string `json:"nav_id"`
		LegalEntity          *string `json:"legal_entity"`
		Name                 *string `json:"name"`
		IndustryID           *string `json:"industry_id"`
		NumberEmployee       *string `json:"number_employee"`
		OfficeEmployee       *int    `json:"office_employee"`
		BusinessSize         *string `json:"business_size"`
		BusinessGroup        *string `json:"business_group"`
		EstablishedYear      *int64  `json:"established_year"`
		IsDelete             bool    `json:"is_delete"`
		TermOfPayment        *string `json:"term_of_payment"`
		CustomerCategory     *string `json:"customer_category"`
		UserID               *int    `json:"user_id"`
		CreatedAt            *string `json:"created_at"`
		ModifiedAt           *string `json:"modified_at"`
		CreatedBy            *int64  `json:"created_by"`
		ModifiedBy           *int64  `json:"modified_by"`
		AccountGroupID       *int64  `json:"accountgroup_id"`
		IsDisabled           bool    `json:"is_disabled"`
		Status               *string `json:"status"`
		PaymentMethodID      *string `json:"payment_method_id"`
		PaymentMethodType    *string `json:"payment_method_type"`
		SubPaymentMethodName *string `json:"sub_payment_method_name"`
		IsCf                 bool    `json:"is_cf"`
		Logo                 *string `json:"logo"`
		IsParent             bool    `json:"is_parent"`
		IsMicrosite          bool    `json:"is_microsite"`
		Membertype           *string `json:"member_type"`
		ErpID                *string `json:"erp_id"`
	}

	// B2BContact data structure
	B2BContactCDC struct {
		ID                   int         `json:"id"`
		ReferenceID          *string     `json:"reference_id"`
		FirstName            *string     `json:"first_name"`
		LastName             *string     `json:"last_name"`
		Salutation           *string     `json:"salutation"`
		JobTitle             *string     `json:"job_title"`
		Email                *string     `json:"email"`
		NavContactID         *string     `json:"nav_contact_id"`
		IsPrimary            bool        `json:"is_primary"`
		BirthDate            interface{} `json:"birth_date"`
		Note                 *string     `json:"note"`
		CreatedAt            *string     `json:"created_at"`
		ModifiedAt           *string     `json:"modified_at"`
		CreatedBy            *int64      `json:"created_by"`
		ModifiedBy           *int64      `json:"modified_by"`
		AccountID            *string     `json:"account_id"`
		IsDisabled           bool        `json:"is_disabled"`
		Password             *string     `json:"password"`
		Avatar               *string     `json:"avatar"`
		Status               *string     `json:"status"`
		Token                *string     `json:"token"`
		IsNew                bool        `json:"is_new"`
		PhoneNumber          *string     `json:"phone_number"`
		OtherPhoneNumber     *string     `json:"other_phone_number"`
		Gender               *string     `json:"gender"`
		TransactionType      interface{} `json:"transaction_type"`
		ErpID                *string     `json:"erp_id"`
		IsSync               bool        `json:"is_sync"`
		Salt                 *string     `json:"salt"`
		LastPasswordModified *time.Time  `json:"last_password_modified"`
	}

	// B2BContactData data structure
	B2BContactData struct {
		ID                   int               `json:"id"`
		ReferenceID          string            `json:"reference_id"`
		FirstName            string            `json:"first_name"`
		LastName             string            `json:"last_name"`
		Salutation           string            `json:"salutation"`
		JobTitle             string            `json:"job_title"`
		Email                string            `json:"email"`
		NavContactID         string            `json:"nav_contact_id"`
		IsPrimary            bool              `json:"is_primary"`
		BirthDate            time.Time         `json:"birth_date"`
		Note                 string            `json:"note"`
		CreatedAt            time.Time         `json:"created_at"`
		ModifiedAt           time.Time         `json:"modified_at"`
		CreatedBy            int64             `json:"created_by"`
		ModifiedBy           int64             `json:"modified_by"`
		AccountID            string            `json:"account_id"`
		IsDisabled           bool              `json:"is_disabled"`
		Password             string            `json:"password,omitempty"`
		Token                string            `json:"token,omitempty"`
		Avatar               string            `json:"avatar"`
		Status               string            `json:"status"`
		IsNew                bool              `json:"is_new"`
		PhoneNumber          string            `json:"phone_number"`
		MemberType           string            `json:"member_type,omitempty"`
		IsMicrosite          bool              `json:"is_microsite"`
		TransactionType      []TransactionType `json:"transaction_type"`
		LoginType            string            `json:"login_type,omitempty"`
		IsSync               bool              `json:"is_sync"`
		Salt                 string            `json:"salt"`
		LastPasswordModified *time.Time        `json:"last_password_modified"`
	}

	Account struct {
		ID                   string    `json:"id"`
		AccountgroupID       int       `json:"accountgroup_id"`
		ParentID             string    `json:"parent_id"`
		LegalEntity          string    `json:"legal_entity"`
		Name                 string    `json:"name"`
		IndustryID           string    `json:"industry_id"`
		NumberEmployee       string    `json:"number_employee"`
		OfficeEmployee       int       `json:"office_employee"`
		BusinessSize         string    `json:"business_size"`
		TermOfPayment        string    `json:"term_of_payment"`
		EstablishedYear      int       `json:"established_year"`
		NavID                string    `json:"nav_id"`
		UserID               int       `json:"user_id"`
		IsDelete             bool      `json:"is_delete"`
		IsDisabled           bool      `json:"is_disabled"`
		PaymentMethodID      string    `json:"payment_method_id"`
		PaymentMethodType    string    `json:"payment_method_type"`
		SubPaymentMethodName string    `json:"sub_payment_method_name"`
		Status               string    `json:"status"`
		IsCf                 bool      `json:"is_cf"`
		IsParent             bool      `json:"is_parent"`
		MemberType           string    `json:"member_type"`
		Logo                 string    `json:"logo"`
		CreatedAt            time.Time `json:"created_at"`
	}

	AccountContact struct {
		ID         int         `json:"id"`
		Status     string      `json:"status"`
		IsDelete   bool        `json:"is_delete"`
		IsDisabled bool        `json:"is_disabled"`
		CreatedAt  time.Time   `json:"created_at"`
		ModifiedAt time.Time   `json:"modified_at"`
		CreatedBy  interface{} `json:"created_by"`
		ModifiedBy int         `json:"modified_by"`
		IsAdmin    bool        `json:"is_admin"`
		Account    string      `json:"account"`
		Contact    int         `json:"contact"`
		Department interface{} `json:"department"`
	}

	// B2BAccountContact data structure
	B2BAccountContactCDC struct {
		ID           int     `json:"id"`
		Status       *string `json:"status"`
		IsDelete     bool    `json:"is_delete"`
		IsDisabled   bool    `json:"is_disabled"`
		AccountID    *string `json:"account_id"`
		ContactID    *int    `json:"contact_id"`
		CreatedAt    *string `json:"created_at"`
		ModifiedAt   *string `json:"modified_at"`
		CreatedBy    *int64  `json:"created_by"`
		ModifiedBy   *int64  `json:"modified_by"`
		IsAdmin      bool    `json:"is_admin"`
		DepartmentID *int64  `json:"department_id"`
	}

	// B2BDocumentCDC data structure
	B2BDocumentCDC struct {
		ID                  int     `json:"id"`
		DocumentType        *string `json:"document_type"`
		DocumentFile        *string `json:"document_file"`
		DocumentTitle       *string `json:"document_title"`
		DocumentDescription *string `json:"document_description"`
		NpwpNumber          *string `json:"npwp_number"`
		NpwpName            *string `json:"npwp_name"`
		NpwpAddress         *string `json:"npwp_address"`
		SiupNumber          *string `json:"siup_number"`
		SiupCompanyName     *string `json:"siup_company_name"`
		SiupType            *string `json:"siup_type"`
		IsDelete            bool    `json:"is_delete"`
		CreatedAt           *string `json:"created_at"`
		ModifiedAt          *string `json:"modified_at"`
		CreatedBy           *int64  `json:"created_by"`
		ModifiedBy          *int64  `json:"modified_by"`
		AccountID           *string `json:"account_id"`
		IsDisabled          bool    `json:"is_disabled"`
		ProvinceID          *string `json:"province_id"`
		ProvinceName        *string `json:"province_name"`
		CityID              *string `json:"city_id"`
		CityName            *string `json:"city_name"`
		DistrictID          *string `json:"district_id"`
		DistrictName        *string `json:"district_name"`
		SubdistrictID       *string `json:"subdistrict_id"`
		SubdistrictName     *string `json:"subdistrict_name"`
		PostalCode          *string `json:"postal_code"`
		Street              *string `json:"street"`
	}

	// B2BContactAddressCDC data structure
	B2BContactAddressCDC struct {
		ID              int     `json:"id"`
		Name            *string `json:"name"`
		PicName         *string `json:"pic_name"`
		Phone           *string `json:"phone"`
		Street          *string `json:"street"`
		ProvinceID      *string `json:"province_id"`
		ProvinceName    *string `json:"province_name"`
		CityID          *string `json:"city_id"`
		CityName        *string `json:"city_name"`
		DistrictID      *string `json:"district_id"`
		DistrictName    *string `json:"district_name"`
		SubDistrictID   *string `json:"subdistrict_id"`
		SubDistrictName *string `json:"subdistrict_name"`
		IsBilling       bool    `json:"is_billing"`
		IsShipping      bool    `json:"is_shipping"`
		PostalCode      *string `json:"postal_code"`
		CreatedAt       *string `json:"created_at"`
		ModifiedAt      *string `json:"modified_at"`
		CreatedBy       *int64  `json:"created_by"`
		ModifiedBy      *int64  `json:"modified_by"`
		IsDelete        bool    `json:"is_delete"`
		ContactID       *int    `json:"contact_id"`
	}

	// B2BLeadsCDC data structure
	B2BLeadsCDC struct {
		ID           int     `json:"id"`
		Name         *string `json:"name"`
		Email        *string `json:"email"`
		Phone        *string `json:"phone"`
		SourceID     *string `json:"source_id"`
		Source       *string `json:"source"`
		SourceStatus *string `json:"source_status"`
		Status       *string `json:"status"`
		Notes        *string `json:"notes"`
		CreatedAt    *string `json:"created_at"`
		ModifiedAt   *string `json:"modified_at"`
		CreatedBy    *int64  `json:"created_by"`
		ModifiedBy   *int64  `json:"modified_by"`
	}
	// B2BPhoneCDC data structure
	B2BPhoneCDC struct {
		ID           int     `json:"id"`
		RelationID   *int    `json:"relation_id"`
		RelationType *string `json:"relation_type"`
		TypePhone    *string `json:"type_phone"`
		Label        *string `json:"label"`
		Number       *string `json:"number"`
		Area         *string `json:"area"`
		Ext          *string `json:"ext"`
		IsPrimary    bool    `json:"is_primary"`
		IsDelete     bool    `json:"is_delete"`
		CreatedAt    *string `json:"created_at"`
		ModifiedAt   *string `json:"modified_at"`
		CreatedBy    *int64  `json:"created_by"`
		ModifiedBy   *int64  `json:"modified_by"`
	}

	// B2BAddress data structure
	B2BAddressCDC struct {
		ID              int     `json:"id"`
		TypeAddress     *string `json:"type_address"`
		Name            *string `json:"name"`
		Street          *string `json:"street"`
		ProvinceID      *string `json:"province_id"`
		ProvinceName    *string `json:"province_name"`
		CityID          *string `json:"city_id"`
		CityName        *string `json:"city_name"`
		DistrictID      *string `json:"district_id"`
		DistrictName    *string `json:"district_name"`
		SubDistrictID   *string `json:"subdistrict_id"`
		SubDistrictName *string `json:"subdistrict_name"`
		IsPrimary       bool    `json:"is_primary"`
		PostalCode      *string `json:"postal_code"`
		CreatedAt       *string `json:"created_at"`
		ModifiedAt      *string `json:"modified_at"`
		CreatedBy       *int64  `json:"created_by"`
		ModifiedBy      *int64  `json:"modified_by"`
		AccountID       *string `json:"account_id"`
		IsDisabled      bool    `json:"is_disabled"`
		IsDelete        bool    `json:"is_delete"`
	}
)

// TransactionType for multiple account
type TransactionType struct {
	Microsite string `json:"microsite"`
	Type      string `json:"type"`
}

// Shark Non-CDC Kafka
type (
	PayloadAccount struct {
		EventType string        `json:"eventType"`
		Payload   B2BAccountCDC `json:"payload"`
	}

	PayloadAccountContact struct {
		EventType string               `json:"eventType"`
		Payload   B2BAccountContactCDC `json:"payload"`
	}
	PayloadContact struct {
		EventType string        `json:"eventType"`
		Payload   B2BContactCDC `json:"payload"`
	}
	PayloadAddress struct {
		EventType string        `json:"eventType"`
		Payload   B2BAddressCDC `json:"payload"`
	}
	PayloadPhone struct {
		EventType string      `json:"eventType"`
		Payload   B2BPhoneCDC `json:"payload"`
	}
	PayloadDocument struct {
		EventType string         `json:"eventType"`
		Payload   B2BDocumentCDC `json:"payload"`
	}

	PayloadContactAddress struct {
		EventType string               `json:"eventType"`
		Payload   B2BContactAddressCDC `json:"payload"`
	}

	PayloadLeads struct {
		EventType string      `json:"eventType"`
		Payload   B2BLeadsCDC `json:"payload"`
	}

	PayloadContactDocument struct {
		EventType string             `json:"eventType"`
		Payload   B2BContactDocument `json:"payload"`
	}
)

// ListContact data structure
type ListContact struct {
	Contact   []*B2BContactData `jsonapi:"relation,contact" json:"contact"`
	TotalData int               `json:"totalData"`
}

// ContactPayloadData data structure

// for CDC sync
type (
	// LeadsPayloadData data structure
	LeadsPayloadData struct {
		After  B2BLeadsCDC `json:"after"`
		Before B2BLeadsCDC `json:"before"`
		Op     string      `json:"op"`
	}

	// LeadsPayloadCDC data structure
	LeadsPayloadCDC struct {
		Payload LeadsPayloadData `json:"payload"`
	}

	// ContactAddressPayloadData data structure
	ContactAddressPayloadData struct {
		After  B2BContactAddressCDC `json:"after"`
		Before B2BContactAddressCDC `json:"before"`
		Op     string               `json:"op"`
	}

	// ContactAddressPayloadCDC data structure
	ContactAddressPayloadCDC struct {
		Payload ContactAddressPayloadData `json:"payload"`
	}
	// DocumentPayloadData data structure
	DocumentPayloadData struct {
		After  B2BDocumentCDC `json:"after"`
		Before B2BDocumentCDC `json:"before"`
		Op     string         `json:"op"`
	}

	// DocumentPayloadCDC data structure
	DocumentPayloadCDC struct {
		Payload DocumentPayloadData `json:"payload"`
	}

	ContactPayloadData struct {
		After  B2BContactCDC `json:"after"`
		Before B2BContactCDC `json:"before"`
		Op     string        `json:"op"`
	}

	// ContactPayloadCDC data structure
	ContactPayloadCDC struct {
		Payload ContactPayloadData `json:"payload"`
	}

	// PhonePayloadCDC data structure
	PhonePayloadCDC struct {
		Payload PhonePayloadData `json:"payload"`
	}

	// PhonePayloadData data structure
	PhonePayloadData struct {
		After  B2BPhoneCDC `json:"after"`
		Before B2BPhoneCDC `json:"before"`
		Op     string      `json:"op"`
	}

	// AddressPayloadData data structure
	AddressPayloadData struct {
		After  B2BAddressCDC `json:"after"`
		Before B2BAddressCDC `json:"before"`
		Op     string        `json:"op"`
	}

	// AddressPayloadCDC data structure
	AddressPayloadCDC struct {
		Payload AddressPayloadData `json:"payload"`
	}

	// AccountPayloadData data structure
	AccountPayloadData struct {
		After  B2BAccountCDC `json:"after"`
		Before B2BAccountCDC `json:"before"`
		Op     string        `json:"op"`
	}

	// AccountContactPayloadData data structure
	AccountContactPayloadData struct {
		After  B2BAccountContactCDC `json:"after"`
		Before B2BAccountContactCDC `json:"before"`
		Op     string               `json:"op"`
	}

	// AccountPayloadCDC data structure
	AccountPayloadCDC struct {
		Payload AccountPayloadData `json:"payload"`
	}

	// AccountContactPayloadCDC data structure
	AccountContactPayloadCDC struct {
		Payload AccountContactPayloadData `json:"payload"`
	}
)
