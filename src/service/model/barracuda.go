package model

// ResponseZipCode area data structure
type ResponseZipCode struct {
	Status  bool
	Message string
	Data    []ZipCode
}

// ZipCodeQueryParameter area data structure
type ZipCodeQueryParameter struct {
	ProvinceID      string
	ProvinceName    string
	CityID          string
	CityName        string
	DistrictID      string
	DistrictName    string
	SubDistrictID   string
	SubDistrictName string
	ZipCode         string
}

//ZipCode data structure
type ZipCode struct {
	Province    Province    `json:"province"`
	City        City        `json:"city"`
	District    District    `json:"district"`
	SubDistrict SubDistrict `json:"subDistrict"`
	ZipCode     int         `json:"zipCode"`
}

//Province data structure
type Province struct {
	ProvinceID   string `json:"id"`
	ProvinceName string `json:"name"`
}

//City data structure
type City struct {
	ProvinceID string `json:"province,omitempty"`
	CityID     string `json:"id"`
	CityName   string `json:"name"`
}

//District data structure
type District struct {
	CityID       string `json:"cityId,omitempty"`
	DistrictID   string `json:"id"`
	DistrictName string `json:"name"`
}

//SubDistrict data structure
type SubDistrict struct {
	DistrictID      string `json:"districtId,omitempty"`
	SubDistrictID   string `json:"id"`
	SubDistrictName string `json:"name"`
}
