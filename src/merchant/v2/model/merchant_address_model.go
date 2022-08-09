package model

import (
	"time"
)

// MerchantAddressType is type int
type MerchantAddressType int

const (
	// WarehouseString const variable
	WarehouseString = "WAREHOUSE"
	// StoreString const variable
	StoreString = "STORE"
	// MainString const variable
	MainString = "MAIN"

	//MobileString const variable for phone type
	MobileString = "MOBILE"
	//PhoneString const variable for phone type
	PhoneString = "PHONE"

	//MerchantString const variable for Relation Name
	MerchantString = "MERCHANT"

	//AddressString const variable for Relation Name
	AddressString = "ADDRESS"
)

// WarehouseData data structure
type WarehouseData struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Phone  string `json:"phone"`
	WarehousePrimary
	DistrictID         string    `json:"districtId"`
	DistrictName       string    `json:"districtName"`
	SubDistrictID      string    `json:"subdistrictId"`
	SubDistrictName    string    `json:"subdistrictName"`
	PostalCode         string    `json:"postalCode"`
	Address            string    `json:"address"`
	Created            time.Time `json:"-"`
	CreatedString      string    `json:"created"`
	LastModified       time.Time `json:"-"`
	LastModifiedString string    `json:"lastModified"`
	CreatedBy          string    `json:"createdBy"`
	ModifiedBy         string    `json:"modifiedBy"`
	Type               string    `json:"type"`
	MerchantID         string    `json:"merchantId"`
	IsPrimary          bool      `json:"isPrimary"`
	Version            int       `json:"version"`
	Status             string    `json:"status"`
	Maps               Maps      `json:"maps"`
	IsMapAvailable     bool      `json:"isMapAvailable"`
}

type WarehousePrimary struct {
	CityID       string `json:"cityId"`
	CityName     string `json:"cityName"`
	ProvinceID   string `json:"provinceId"`
	ProvinceName string `json:"provinceName"`
}

type ListPublicWarehouse struct {
	WarehouseData []*WarehousePrimary `jsonapi:"relation,warehouseData" json:"warehouseData"`
	TotalData     int                 `json:"totalData"`
}

// AddressData data structure
type AddressData struct {
	ID                 string    `json:"id"`
	RelationID         string    `json:"relationId"`
	RelationName       string    `json:"relationName"`
	Label              string    `json:"label"`
	ProvinceID         string    `json:"provinceId"`
	ProvinceName       string    `json:"provinceName"`
	CityID             string    `json:"cityId"`
	CityName           string    `json:"cityName"`
	DistrictID         string    `json:"districtId"`
	DistrictName       string    `json:"districtName"`
	SubDistrictID      string    `json:"subdistrictId"`
	SubDistrictName    string    `json:"subdistrictName"`
	PostalCode         string    `json:"postalCode"`
	Address            string    `json:"address"`
	Created            time.Time `json:"-"`
	CreatedString      string    `json:"created"`
	LastModified       time.Time `json:"-"`
	LastModifiedString string    `json:"lastModified"`
	CreatedBy          string    `json:"createdBy"`
	ModifiedBy         string    `json:"modifiedBy"`
	Type               string    `json:"type"`
	IsPrimary          bool      `json:"isPrimary"`
	Version            int       `json:"version"`
	Status             string    `json:"status"`
}

// PhoneData data structure
type PhoneData struct {
	ID                 string    `json:"id"`
	RelationID         string    `json:"relationId"`
	RelationName       string    `json:"relationName"`
	Label              string    `json:"label"`
	Number             string    `json:"number"`
	Type               string    `json:"type"`
	IsPrimary          bool      `json:"isPrimary"`
	Version            int       `json:"version"`
	Created            time.Time `json:"-"`
	CreatedString      string    `json:"created"`
	LastModified       time.Time `json:"-"`
	LastModifiedString string    `json:"lastModified"`
	CreatedBy          string    `json:"createdBy"`
	ModifiedBy         string    `json:"modifiedBy"`
}

// ParameterPrimaryWarehouse data structure
type ParameterPrimaryWarehouse struct {
	MerchantID string `json:"merchantId"`
	AddressID  string `json:"addressId"`
	MemberID   string `json:"memberId"`
}

// RestructToAddress data structure
func RestructToAddress(data WarehouseData) AddressData {
	return AddressData{
		ID:                 data.ID,
		RelationID:         data.MerchantID,
		RelationName:       MerchantString,
		Label:              data.Label,
		ProvinceID:         data.ProvinceID,
		ProvinceName:       data.ProvinceName,
		CityID:             data.CityID,
		CityName:           data.CityName,
		DistrictID:         data.DistrictID,
		DistrictName:       data.DistrictName,
		SubDistrictID:      data.SubDistrictID,
		SubDistrictName:    data.SubDistrictName,
		PostalCode:         data.PostalCode,
		Address:            data.Address,
		Created:            data.Created,
		CreatedString:      data.CreatedString,
		LastModified:       data.LastModified,
		LastModifiedString: data.LastModifiedString,
		CreatedBy:          data.CreatedBy,
		ModifiedBy:         data.ModifiedBy,
		Type:               data.Type,
		IsPrimary:          data.IsPrimary,
		Version:            data.Version,
		Status:             data.Status,
	}
}

// ParameterWarehouse data structure
type ParameterWarehouse struct {
	StrPage    string `json:"strPage" form:"strPage" query:"strPage" validate:"omitempty,numeric" fieldname:"strPage" url:"strPage"`
	Page       int    `json:"page" form:"page" query:"page" validate:"omitempty,numeric" fieldname:"page" url:"page"`
	StrLimit   string `json:"strLimit" form:"strLimit" query:"strLimit" validate:"omitempty" fieldname:"strLimit" url:"strLimit"`
	Limit      int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,numeric" fieldname:"limit" url:"limit"`
	Offset     int    `json:"offset" form:"offset" query:"offset" validate:"omitempty,numeric" fieldname:"offset" url:"offset"`
	MemberID   string `json:"memberId" form:"memberId" query:"memberId" validate:"omitempty,string" fieldname:"memberId" url:"memberId"`
	MerchantID string `json:"merchantId" form:"merchantId" query:"merchantId" validate:"omitempty,string" fieldname:"merchantId" url:"merchantId"`
	Query      string `json:"query" form:"query" query:"query" validate:"omitempty,string" fieldname:"query" url:"query"`
	OrderBy    string `json:"order" query:"order"`
	Sort       string `json:"ort" query:"sort"`
	ShowAll    string `json:"showAll" query:"showAll"`
}

// ListWarehouse data structure
type ListWarehouse struct {
	WarehouseData []*WarehouseData `jsonapi:"relation,warehouseData" json:"warehouseData"`
	TotalData     int              `json:"totalData"`
}

// ListPhone data structure
type ListPhone struct {
	PhoneData []*PhoneData `jsonapi:"relation,phoneData" json:"phoneData"`
	TotalData int          `json:"totalData"`
}

func (input *ListWarehouse) RestructToPublic() {
	for _, wh := range input.WarehouseData {
		wh.ID = ""
		wh.Label = ""
		wh.Name = ""
		wh.Mobile = ""
		wh.Phone = ""
		wh.DistrictID = ""
		wh.DistrictName = ""
		wh.SubDistrictID = ""
		wh.SubDistrictName = ""
		wh.PostalCode = ""
		wh.Address = ""
		wh.CreatedString = ""
		wh.LastModifiedString = ""
		wh.CreatedBy = ""
		wh.ModifiedBy = ""
		wh.Type = ""
		wh.MerchantID = ""
		wh.IsPrimary = false
		wh.Version = 0
		wh.Status = ""
	}
}
