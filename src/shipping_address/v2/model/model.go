package model

import (
	"time"

	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

const (
	// Module name
	Module = "shippingAddress"
	// MaximumShippingAddress maximum data
	MaximumShippingAddress = 20
	//Unauthorized message
	Unauthorized = "User Unauthorized"
)

// ShippingAddressError data structure
type ShippingAddressError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// SuccessResponse data structure
type SuccessResponse struct {
	ID      string `jsonapi:"primary" json:"id"`
	Message string `jsonapi:"attr,message" json:"message"`
}

// ShippingAddress data structure
type ShippingAddress struct {
	ID              string  `json:"id"`
	MemberID        *string `json:"memberId"`
	Name            *string `json:"name"`
	Mobile          *string `json:"mobile"`
	Phone           *string `json:"phone"`
	ProvinceID      *string `json:"provinceId"`
	ProvinceName    *string `json:"provinceName"`
	CityID          *string `json:"cityId"`
	CityName        *string `json:"cityName"`
	DistrictID      *string `json:"districtId"`
	DistrictName    *string `json:"districtName"`
	SubDistrictID   *string `json:"subdistrictId"`
	SubDistrictName *string `json:"subdistrictName"`
	PostalCode      *string `json:"postalCode"`
	Street1         *string `json:"street1"`
	Street2         *string `json:"street2"`
	Version         *int    `json:"version"`
	Created         *string `json:"created"`
	LastModified    *string `json:"lastModified"`
	Ext             *string `json:"ext"`
	Label           *string `json:"label"`
	IsPrimary       bool    `json:"isPrimary"`
}

// ShippingAddressData data structure
type ShippingAddressData struct {
	ID              string    `json:"id"`
	MemberID        string    `json:"memberId"`
	Name            string    `json:"name"`
	Mobile          string    `json:"mobile"`
	Phone           string    `json:"phone"`
	ProvinceID      string    `json:"provinceId"`
	ProvinceName    string    `json:"provinceName"`
	CityID          string    `json:"cityId"`
	CityName        string    `json:"cityName"`
	DistrictID      string    `json:"districtId"`
	DistrictName    string    `json:"districtName"`
	SubDistrictID   string    `json:"subdistrictId"`
	SubDistrictName string    `json:"subdistrictName"`
	PostalCode      string    `json:"postalCode"`
	Street1         string    `json:"street1"`
	Street2         string    `json:"street2"`
	Ext             string    `json:"ext"`
	Version         int       `json:"version"`
	IsPrimary       bool      `json:"isPrimary"`
	Created         time.Time `json:"created"`
	LastModified    time.Time `json:"lastModified"`
	CreatedBy       string    `json:"createdBy"`
	ModifiedBy      string    `json:"modifiedBy"`
	IsMapAvailable  bool      `json:"isMapAvailable"`
	MapsID          string    `json:"mapsId,omitempty"`
	RelationID      string    `json:"relationId,omitempty"`
	RelationName    string    `json:"relationName,omitempty"`
	Label           string    `json:"label,omitempty"`
	Latitude        float64   `json:"latitude,omitempty"`
	Longitude       float64   `json:"longitude,omitempty"`
}

// ParametersShippingAddress data structure
type ParametersShippingAddress struct {
	StrPage  string `json:"strPage" form:"strPage" query:"strPage" validate:"omitempty,numeric" fieldname:"strPage" url:"strPage"`
	Page     int    `json:"page" form:"page" query:"page" validate:"omitempty,numeric" fieldname:"page" url:"page"`
	StrLimit string `json:"strLimit" form:"strLimit" query:"strLimit" validate:"omitempty" fieldname:"strLimit" url:"strLimit"`
	Limit    int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,numeric" fieldname:"limit" url:"limit"`
	Offset   int    `json:"offset" form:"offset" query:"offset" validate:"omitempty,numeric" fieldname:"offset" url:"offset"`
	MemberID string `json:"memberId" form:"memberId" query:"memberId" validate:"omitempty,string" fieldname:"memberId" url:"memberId"`
	Query    string `json:"query" form:"query" query:"query" validate:"omitempty,string" fieldname:"query" url:"query"`
}

// ListShippingAddress data structure
type ListShippingAddress struct {
	ShippingAddress []*ShippingAddressData `jsonapi:"relation,shippingAddress" json:"shippingAddress"`
	TotalData       int                    `json:"totalData"`
}

// ParamaterPrimaryShippingAddress data structure
type ParamaterPrimaryShippingAddress struct {
	MemberID   string `json:"memberId"`
	ShippingID string `json:"shippingId"`
	UserID     string `json:"userId"`
}

// RestructShippingAddress function for restruct from cdc
func RestructShippingAddress(pl serviceModel.ShippingAddressPayloadCDC) ShippingAddress {
	var shippingAddress ShippingAddress
	shippingAddress.ID = pl.Payload.After.ID
	shippingAddress.MemberID = pl.Payload.After.MemberID
	shippingAddress.Name = pl.Payload.After.Name
	shippingAddress.Mobile = pl.Payload.After.Mobile
	shippingAddress.Phone = pl.Payload.After.Phone
	shippingAddress.ProvinceID = pl.Payload.After.ProvinceID
	shippingAddress.ProvinceName = pl.Payload.After.ProvinceName
	shippingAddress.CityID = pl.Payload.After.CityID
	shippingAddress.CityName = pl.Payload.After.CityName
	shippingAddress.DistrictID = pl.Payload.After.DistrictID
	shippingAddress.DistrictName = pl.Payload.After.DistrictName
	shippingAddress.SubDistrictID = pl.Payload.After.SubDistrictID
	shippingAddress.SubDistrictName = pl.Payload.After.SubDistrictName
	shippingAddress.PostalCode = pl.Payload.After.PostalCode
	shippingAddress.Street1 = pl.Payload.After.Street1
	shippingAddress.Street2 = pl.Payload.After.Street2
	shippingAddress.Version = pl.Payload.After.Version
	shippingAddress.Created = pl.Payload.After.Created
	shippingAddress.LastModified = pl.Payload.After.LastModified
	shippingAddress.Ext = pl.Payload.After.Ext
	shippingAddress.Label = pl.Payload.After.Label
	shippingAddress.IsPrimary = pl.Payload.After.IsPrimary
	return shippingAddress
}

type ShippingAddressLog struct {
	Before *ShippingAddressData `json:"before"`
	After  *ShippingAddressData `json:"after"`
}
