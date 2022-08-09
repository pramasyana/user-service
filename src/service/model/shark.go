package model

// B2BAccountTemporary data structure
type B2BAccountTemporary struct {
	ID               int     `json:"id"`
	IndustryID       *string `json:"industry_id"`
	NumberEmployee   *string `json:"number_employee"`
	OfficeEmployee   *int64  `json:"office_employee"`
	BusinessSize     *string `json:"business_size"`
	OrganizationName *string `json:"organization_name"`
	LegalEntity      *string `json:"legal_entity"`
	Name             *string `json:"name"`
	Salutation       *string `json:"salutation"`
	Email            *string `json:"email"`
	Password         *string `json:"password"`
	Token            *string `json:"token"`
	IsDelete         bool    `json:"is_delete"`
	IsDisabled       bool    `json:"is_disabled"`
	Phone            *string `json:"phone"`
	ParentID         *string `json:"parent_id"`
	NpwpNumber       *string `json:"npwp_number"`
}

// AccountTemporaryPayloadData data structure
type AccountTemporaryPayloadData struct {
	After  B2BAccountTemporary `json:"after"`
	Before B2BAccountTemporary `json:"before"`
	Op     string              `json:"op"`
}

// AccountTemporaryPayloadCDC data structure
type AccountTemporaryPayloadCDC struct {
	Payload AccountTemporaryPayloadData `json:"payload"`
}

// B2BContact data structure
type B2BContact struct {
	ID               int         `json:"id"`
	ReferenceID      *string     `json:"reference_id"`
	FirstName        *string     `json:"first_name"`
	LastName         *string     `json:"last_name"`
	Salutation       *string     `json:"salutation"`
	JobTitle         *string     `json:"job_title"`
	Email            *string     `json:"email"`
	NavContactID     *string     `json:"nav_contact_id"`
	IsPrimary        bool        `json:"is_primary"`
	BirthDate        interface{} `json:"birth_date"`
	Note             *string     `json:"note"`
	CreatedAt        *string     `json:"created_at"`
	ModifiedAt       *string     `json:"modified_at"`
	CreatedBy        *int64      `json:"created_by"`
	ModifiedBy       *int64      `json:"modified_by"`
	AccountID        *string     `json:"account_id"`
	IsDisabled       bool        `json:"is_disabled"`
	Password         *string     `json:"password"`
	Avatar           *string     `json:"avatar"`
	Status           *string     `json:"status"`
	Token            *string     `json:"token"`
	IsNew            bool        `json:"is_new"`
	PhoneNumber      *string     `json:"phone_number"`
	OtherPhoneNumber *string     `json:"other_phone_number"`
	Gender           *string     `json:"gender"`
	TransactionType  interface{} `json:"transaction_type"`
}

// TransactionType for multiple account
type TransactionType struct {
	Microsite string `json:"microsite"`
	Type      string `json:"type"`
}

// ContactPayloadData data structure
type ContactPayloadData struct {
	After  B2BContact `json:"after"`
	Before B2BContact `json:"before"`
	Op     string     `json:"op"`
}

// ContactPayloadCDC data structure
type ContactPayloadCDC struct {
	Payload ContactPayloadData `json:"payload"`
}

// B2BAddress data structure
type B2BAddress struct {
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

// B2BContactNpwp data structure
type B2BContactNpwp struct {
	ID         int     `json:"id"`
	Number     *string `json:"number"`
	Name       *string `json:"name"`
	Address    *string `json:"address"`
	Address2   *string `json:"address2"`
	ContactID  *int    `json:"contact_id"`
	File       *string `json:"file"`
	CreatedAt  *string `json:"created_at"`
	ModifiedAt *string `json:"modified_at"`
	CreatedBy  *int64  `json:"created_by"`
	ModifiedBy *int64  `json:"modified_by"`
	IsDelete   bool    `json:"is_delete"`
}

// ContactNpwpPayloadData data structure
type ContactNpwpPayloadData struct {
	After  B2BContactNpwp `json:"after"`
	Before B2BContactNpwp `json:"before"`
	Op     string         `json:"op"`
}

// ContactNpwpPayloadCDC data structure
type ContactNpwpPayloadCDC struct {
	Payload ContactNpwpPayloadData `json:"payload"`
}

// B2BContactTemp data structure
type B2BContactTemp struct {
	ID           int         `json:"id"`
	FirstName    *string     `json:"first_name"`
	LastName     *string     `json:"last_name"`
	Salutation   *string     `json:"salutation"`
	JobTitle     *string     `json:"job_title"`
	Email        *string     `json:"email"`
	NavContactID *string     `json:"nav_contact_id"`
	IsPrimary    bool        `json:"is_primary"`
	BirthDate    interface{} `json:"birth_date"`
	Note         *string     `json:"note"`
	CreatedAt    *string     `json:"created_at"`
	ModifiedAt   *string     `json:"modified_at"`
	CreatedBy    *int64      `json:"created_by"`
	ModifiedBy   *int64      `json:"modified_by"`
	AccountID    *string     `json:"account_id"`
	ReferenceID  *string     `json:"reference_id"`
}

// ContactTempPayloadData data structure
type ContactTempPayloadData struct {
	After  B2BContactTemp `json:"after"`
	Before B2BContactTemp `json:"before"`
	Op     string         `json:"op"`
}

// ContactTempPayloadCDC data structure
type ContactTempPayloadCDC struct {
	Payload ContactTempPayloadData `json:"payload"`
}
