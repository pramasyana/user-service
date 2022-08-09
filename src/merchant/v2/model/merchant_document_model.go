package model

import (
	"gopkg.in/guregu/null.v4"
)

// B2CMerchantDocument data structure
type B2CMerchantDocument struct {
	ID                     string  `json:"id"`
	MerchantID             *string `json:"merchantID"`
	DocumentType           *string `json:"documentType"`
	DocumentValue          *string `json:"documentValue"`
	DocumentOriginal       *string `json:"documentOriginal"`
	DocumentExpirationDate *string `json:"documentExpirationDate"`
	CreatorID              *string `json:"creatorID"`
	CreatorIP              *string `json:"creatorIP"`
	EditorID               *string `json:"editorID"`
	EditorIP               *string `json:"editorIP"`
	Version                *int    `json:"version"`
	Created                *string `json:"created"`
	LastModified           *string `json:"lastModified"`
}

// B2CMerchantDocumentInput data structure
type B2CMerchantDocumentInput struct {
	ID                     string `json:"id,omitempty" query:"id" validate:"omitempty,numeric" fieldname:"merchant document ID"`
	MerchantID             string `json:"merchantId,omitempty" query:"merchantId" validate:"required,lte=25" fieldname:"merchant ID"`
	DocumentType           string `json:"documentType,omitempty" query:"documentType" validate:"required,gte=3,lte=250" fieldname:"tipe dokumen"`
	DocumentValue          string `json:"documentValue,omitempty" query:"documentValue" validate:"omitempty,is-s3-url" fieldname:"data dokumen"`
	DocumentExpirationDate string `json:"documentExpirationDate,omitempty" query:"documentExpirationDate" validate:"omitempty,is-store-date-format" fieldname:"tanggal kadaluarsa dokumen"`
}

// B2CMerchantDocumentQueryInput query input data structure
type B2CMerchantDocumentQueryInput struct {
	OrderBy      string `json:"orderBy" form:"orderBy" query:"orderBy" validate:"omitempty,oneof=id merchantId documentType" fieldname:"order by"`
	DocumentType string `json:"documentType" form:"documentType" query:"documentType" validate:"omitempty,gte=3,lte=250" fieldName:"tipe dokumen"`
	MerchantID   string `json:"merchantId" form:"merchantId" query:"merchantId" validate:"omitempty,lte=25" fieldName:"tipe dokumen"`
}

// B2CMerchantDocumentData data structure
type B2CMerchantDocumentData struct {
	ID                     string    `json:"id"`
	MerchantID             string    `json:"merchantID"`
	DocumentType           string    `json:"documentType"`
	DocumentValue          string    `json:"documentValue"`
	DocumentOriginal       string    `json:"documentOriginal"`
	DocumentExpirationDate null.Time `json:"documentExpirationDate"`
	CreatorID              string    `json:"creatorID"`
	CreatorIP              string    `json:"creatorIP"`
	EditorID               string    `json:"editorID"`
	EditorIP               string    `json:"editorIP"`
	Version                int       `json:"version"`
	Created                null.Time `json:"created"`
	LastModified           null.Time `json:"lastModified"`
}

// ListB2CMerchantDocument data structure
type ListB2CMerchantDocument struct {
	MerchantDocument []B2CMerchantDocumentData `jsonapi:"relation,merchantBank" json:"merchantBank"`
}
