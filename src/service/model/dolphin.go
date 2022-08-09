package model

import (
	"time"
)

// MemberDolphin data structure special for
type MemberDolphin struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Gender          string `json:"gender"`
	DOB             string `json:"dob"`
	Phone           string `json:"phone"`
	Ext             string `json:"ext"`
	Mobile          string `json:"mobile"`
	Street1         string `json:"street1"`
	Street2         string `json:"street2"`
	PostalCode      string `json:"postalCode"`
	SubDistrictID   string `json:"subDistrictId"`
	SubDistrictName string `json:"subDistrictName"`
	DistrictID      string `json:"districtId"`
	DistrictName    string `json:"districtName"`
	CityID          string `json:"cityId"`
	CityName        string `json:"cityName"`
	ProvinceID      string `json:"provinceId"`
	ProvinceName    string `json:"provinceName"`
	Status          string `json:"status"`
	Created         string `json:"created"`
	LastModified    string `json:"lastModified"`
	Destination     string `json:"destination"`
	FacebookID      string `json:"facebookId"`
	GoogleID        string `json:"googleId"`
	AppleID         string `json:"appleId"`
	AzureID         string `json:"azureId"`
	LDAPID          string `json:"ldapId"`
	Message         string `json:"message"`
}

// DolphinShippingRequests data struct
type DolphinShippingRequests struct {
	ID              string    `url:"id"`
	Label           string    `url:"label"`
	UserID          string    `url:"userId"`
	Name            string    `url:"name,omitempty"`
	Mobile          string    `url:"mobile,omitempty"`
	Phone           string    `url:"phone,omitempty"`
	ProvinceID      string    `url:"provinceId"`
	ProvinceName    string    `url:"provinceName"`
	CityID          string    `url:"cityId"`
	CityName        string    `url:"cityName"`
	DistrictID      string    `url:"districtId"`
	DistrictName    string    `url:"districtName"`
	SubDistrictID   string    `url:"subdistrictId"`
	SubDistrictName string    `url:"subdistrictName"`
	PostalCode      string    `url:"postalCode"`
	Street1         string    `url:"street1"`
	Street2         string    `url:"street2,omitempty"`
	Ext             string    `url:"ext"`
	Created         time.Time `url:"created"`
	LastModified    time.Time `url:"lastModified"`
	Destination     string    `url:"destination"`
}

// DolphinPayloadNSQ data structure for pushing to nsq
type DolphinPayloadNSQ struct {
	EventOrchestration     string        `json:"eventOrchestration,omitempty"`
	TimestampOrchestration string        `json:"timestampOrchestration,omitempty"`
	EventType              string        `json:"eventType"`
	Counter                int           `json:"counter"`
	Payload                MemberDolphin `json:"payload"`
}

// DolphinShippingPayloadNSQ data structure for pushing to nsq
type DolphinShippingPayloadNSQ struct {
	EventType string                  `json:"eventType"`
	Counter   int                     `json:"counter"`
	Payload   DolphinShippingRequests `json:"payload"`
}

// Response data structure for dolphin response
type Response struct {
	Data Data `json:"data"`
}

// Data response
type Data struct {
	Type       string     `json:"type"`
	ID         int        `json:"id"`
	Attributes Attributes `json:"attributes"`
}

// Attributes response
type Attributes struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

// MemberResponse data structure for dolphin response
type MemberResponse struct {
	Data DataMember `json:"data"`
}

// DataMember from dolphin
type DataMember struct {
	Type       string        `json:"type"`
	ID         int           `json:"id"`
	Attributes MemberDolphin `json:"attributes"`
}
