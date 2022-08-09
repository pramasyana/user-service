package model

const (
	// MessageSuccess constanta
	MessageSuccess = "Success Get All Phone Area"
)

// PhoneArea data structure
type PhoneArea struct {
	CodeArea     string `json:"extensionCode"`
	AreaName     string `json:"areaName"`
	ProvinceName string `json:"provinceName"`
}

// PhoneAreaResponse struct
type PhoneAreaResponse struct {
	Status  string      `json:"status"`
	Data    []PhoneArea `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

// TotalPhoneArea data structure
type TotalPhoneArea struct {
	TotalData int `json:"totalData"`
}

// PhoneAreaError data structure
type PhoneAreaError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
