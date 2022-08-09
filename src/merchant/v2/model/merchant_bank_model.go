package model

import "gopkg.in/guregu/null.v4/zero"

// B2CMerchantBank data structure
type B2CMerchantBank struct {
	ID           int     `json:"id"`
	BankCode     *string `json:"BankCode"`
	BankName     *string `json:"BankName"`
	Status       bool    `json:"Status"`
	CreatorID    *string `json:"CreatorID"`
	CreatorIP    *string `json:"creatorIP"`
	Created      *string `json:"created"`
	EditorID     *string `json:"editorID"`
	EditorIP     *string `json:"editorIP"`
	LastModified *string `json:"lastModified"`
	DeletedAt    *string `json:"deletedAt"`
}

// B2CMerchantBankData data structure
type B2CMerchantBankData struct {
	ID           int         `json:"id"`
	BankCode     string      `json:"BankCode"`
	BankName     zero.String `json:"BankName"`
	Status       bool        `json:"Status"`
	CreatorID    string      `json:"CreatorID"`
	CreatorIP    string      `json:"creatorIP"`
	Created      string      `json:"created"`
	EditorID     string      `json:"editorID"`
	EditorIP     string      `json:"editorIP"`
	LastModified string      `json:"lastModified"`
}

// ListMerchantBank data structure
type ListMerchantBank struct {
	MerchantBank []*B2CMerchantBankData `jsonapi:"relation,merchantBank" json:"merchantBank"`
	TotalData    int                    `json:"totalData"`
}

// ParametersMerchantBank data structure
type ParametersMerchantBank struct {
	StrPage  string `json:"strPage" form:"strPage" query:"strPage" validate:"omitempty,numeric" fieldname:"strPage" url:"strPage"`
	Page     int    `json:"page" form:"page" query:"page" validate:"omitempty,numeric" fieldname:"page" url:"page"`
	StrLimit string `json:"strLimit" form:"strLimit" query:"strLimit" validate:"omitempty" fieldname:"strLimit" url:"strLimit"`
	Limit    int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,numeric" fieldname:"limit" url:"limit"`
	Offset   int    `json:"offset" form:"offset" query:"offset" validate:"omitempty,numeric" fieldname:"offset" url:"offset"`
	Sort     string `json:"sort" form:"sort" query:"sort" validate:"omitempty,oneof=asc desc" fieldname:"sort" url:"sort"`
	OrderBy  string `json:"orderBy" form:"orderBy" query:"orderBy" validate:"omitempty,oneof=id bankCode bankName" fieldname:"orderBy" url:"orderBy"`
	Status   string `json:"status" query:"status" validate:"omitempty,is-bool" fieldname:"status bank" url:"status"`
}

// AllowedSortFieldsMerchantBank is allowed field name for sorting
var AllowedSortFieldsMerchantBank = []string{
	"id",
	"bankCode",
	"bankName",
}
