package model

// CorporateError data structure
type CorporateError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// SuccessResponse data structure
type SuccessResponse struct {
	ID      string `jsonapi:"primary" json:"id"`
	Message string `jsonapi:"attr,message" json:"message"`
}

const (
	LoginTypeCorporate = "corporate"
)

type ImportFile struct {
	File string `json:"file"`
}
