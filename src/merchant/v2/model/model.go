package model

const (
	// Module name
	Module = "merchant"
	// ModuleWarehouse name
	ModuleWarehouse = "merchantWarehouse"

	MerchantFailedUpdateError = "failed to update merchant"
)

// MerchantError data structure
type MerchantError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// SuccessResponse data structure
type SuccessResponse struct {
	ID      string `jsonapi:"primary" json:"id"`
	Message string `jsonapi:"attr,message" json:"message"`
}
