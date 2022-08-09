package model

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

// ShippingAddressPayloadData data structure
type ShippingAddressPayloadData struct {
	After  ShippingAddress `json:"after"`
	Before ShippingAddress `json:"before"`
	Op     string          `json:"op"`
}

// ShippingAddressPayloadCDC data structure
type ShippingAddressPayloadCDC struct {
	Payload ShippingAddressPayloadData `json:"payload"`
}

// ShippingAddressPayloadNSQ data structure for pushing to nsq
type ShippingAddressPayloadNSQ struct {
	EventOrchestration     string          `json:"eventOrchestration,omitempty"`
	TimestampOrchestration string          `json:"timestampOrchestration,omitempty"`
	EventType              string          `json:"eventType"`
	Counter                int             `json:"counter"`
	Payload                ShippingAddress `json:"payload"`
}

// SturgeonShippingPayloadNSQ data structure for pushing to nsq
type SturgeonShippingPayloadNSQ struct {
	EventType string                  `json:"eventType"`
	Counter   int                     `json:"counter"`
	Payload   DolphinShippingRequests `json:"payload"`
}
