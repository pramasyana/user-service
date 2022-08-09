package model

import (
	"strings"
	"time"
)

const (

	// Waiting document status
	Waiting = 100
	// Onprocess document status
	Onprocess = 200
	// Verified document status
	Verified = 500
	// Rejected document status
	Rejected = -500

	// WaitingString const variable
	WaitingString = "WAITING"
	// OnprocessString const variable
	OnprocessString = "ONPROCESS"
	// VerifiedString const variable
	VerifiedString = "VERIFIED"
	// RejectedString const variable
	RejectedString = "REJECTED"
	//Unauthorized access
	Unauthorized = "User Unauthorized"
)

// StringToStatus function for converting string document status to int
func StringToStatus(s string) int {
	switch strings.ToUpper(s) {
	case WaitingString:
		return Waiting
	case OnprocessString:
		return Onprocess
	case VerifiedString:
		return Verified
	case RejectedString:
		return Rejected
	}
	return Waiting
}

// DocumentError data structure
type DocumentError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// SuccessResponse data structure
type SuccessResponse struct {
	ID      string `jsonapi:"primary" json:"id"`
	Message string `jsonapi:"attr,message" json:"message"`
}

// DocumentQueryInput query input data structure
type DocumentQueryInput struct {
	ID       string `json:"id" form:"id" query:"id" fieldName:"id"`
	MemberID string `json:"memberId" form:"memberId" query:"memberId" fieldName:"memberId"`
}

// DocumentParameters data structure
type DocumentParameters struct {
	StrPage  string `json:"strPage" form:"strPage" query:"strPage" validate:"omitempty,numeric" fieldname:"strPage" url:"strPage"`
	Page     int    `json:"page" form:"page" query:"page" validate:"omitempty,numeric" fieldname:"page" url:"page"`
	StrLimit string `json:"strLimit" form:"strLimit" query:"strLimit" validate:"omitempty" fieldname:"strLimit" url:"strLimit"`
	Limit    int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,numeric" fieldname:"limit" url:"limit"`
	Offset   int    `json:"offset" form:"offset" query:"offset" validate:"omitempty,numeric" fieldname:"offset" url:"offset"`
	MemberID string `json:"memberId" form:"memberId" query:"memberId" validate:"omitempty,string" fieldname:"memberId" url:"memberId"`
	Query    string `json:"query" form:"query" query:"query" validate:"omitempty,string" fieldname:"query" url:"query"`
	ID       string `json:"id" form:"id" query:"id" fieldName:"id"`
}

// DocumentTypeParameters data structure
type DocumentTypeParameters struct {
	StrPage      string `json:"strPage" form:"strPage" query:"strPage" validate:"omitempty,numeric" fieldname:"strPage" url:"strPage"`
	Page         int    `json:"page" form:"page" query:"page" validate:"omitempty,numeric" fieldname:"page" url:"page"`
	StrLimit     string `json:"strLimit" form:"strLimit" query:"strLimit" validate:"omitempty" fieldname:"strLimit" url:"strLimit"`
	Limit        int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,numeric" fieldname:"limit" url:"limit"`
	Offset       int    `json:"offset" form:"offset" query:"offset" validate:"omitempty,numeric" fieldname:"offset" url:"offset"`
	IsB2b        string `json:"isB2b" form:"isB2b" query:"isB2b" validate:"omitempty,string" fieldname:"isB2b" url:"isB2b"`
	IsB2c        string `json:"isB2c" form:"isB2c" query:"isB2c" validate:"omitempty,string" fieldname:"isB2c" url:"isB2c"`
	ID           string `json:"id" form:"id" query:"id" fieldName:"id"`
	DocumentType string `json:"documentType" form:"documentType" query:"documentType" validate:"omitempty,gte=3,lte=250" fieldName:"documentType"`
	IsActive     string `json:"isActive" form:"isActive" query:"isActive" fieldName:"isActive"`
}

// DocumentData data structure
type DocumentData struct {
	ID           string    `json:"id" db:"id" form:"id"`
	MemberID     string    `json:"memberId" db:"memberId" form:"memberId"`
	DocumentType string    `json:"documentType" db:"documentType" form:"documentType"`
	DocumentFile string    `json:"documentFile" db:"documentFile" form:"documentFile"`
	Title        string    `json:"title" db:"title" form:"title"`
	Number       string    `json:"number" db:"number" form:"number"`
	Status       int       `json:"status" db:"status"`
	StatusText   string    `json:"statusText" db:"status"`
	Description  string    `json:"description" db:"description" form:"description"`
	IsDelete     bool      `json:"isDelete" db:"isDelete"`
	Created      time.Time `json:"created" db:"created"`
	LastModified time.Time `json:"lastModified" db:"lastModified"`
	CreatedBy    string    `json:"createdBy" db:"createdBy"`
	ModifiedBy   string    `json:"modifiedBy" db:"modifiedBy"`
}

// ListDocument data structure
type ListDocument struct {
	Document  []*DocumentData `jsonapi:"relation,document" json:"document"`
	TotalData int             `json:"totalData"`
}

// DocumentTypePayload data structure
type DocumentTypePayload struct {
	DocumentType string `json:"documentType" form:"documentType"`
	IsB2c        string `json:"isB2c" form:"isB2c"`
	IsB2b        string `json:"isB2b" form:"isB2b"`
	IsActive     string `json:"isActive" form:"isActive"`
}

// DocumentType data structure
type DocumentType struct {
	ID             string    `json:"id"`
	DocumentType   string    `json:"documentType" form:"documentType"`
	IsB2c          bool      `json:"isB2c" form:"isB2c"`
	IsB2cString    string    `json:"-"`
	IsB2b          bool      `json:"isB2b" form:"isB2b"`
	IsB2bString    string    `json:"-"`
	IsActive       bool      `json:"isActive" form:"isActive"`
	IsActiveString string    `json:"-"`
	Created        time.Time `json:"created"`
	LastModified   time.Time `json:"lastModified"`
	CreatedBy      string    `json:"createdBy"`
	ModifiedBy     string    `json:"modifiedBy"`
}

// ListDocumentType data structure
type ListDocumentType struct {
	DocumentType []*DocumentType `jsonapi:"relation,documentType" json:"documentType"`
	TotalData    int             `json:"totalData"`
}

// RequiredDocuments data structure
type RequiredDocuments struct {
	Merchant      MerchantType `json:"document"`
	TotalDocument int          `json:"totalDocument"`
}

//MerchantType data structure
type MerchantType struct {
	PerseoranganPkp    []DocumentRequire `json:"perseorangan-pkp"`
	PerseoranganNonPkp []DocumentRequire `json:"perseorangan-nonpkp"`
	RegulerPkp         []DocumentRequire `json:"reguler-pkp"`
	RegulerNonPkp      []DocumentRequire `json:"reguler-non-pkp"`
	Upgrade            []DocumentRequire `json:"upgrade"`
}

// DocumentRequire data structure
type DocumentRequire struct {
	Identifier string `json:"identifier"`
	Label      string `json:"label"`
	Extra      string `json:"extra"`
	New        bool   `json:"new"`
	MaxSize    int    `json:"max-size"`
	Mandatory  bool   `json:"mandatory"`
}
