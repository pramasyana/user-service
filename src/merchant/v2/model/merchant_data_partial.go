package model

import (
	"time"

	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

type B2CMerchantDataPatial struct {
	ID                       string                    `json:"id" db:"id"`
	UserID                   string                    `json:"userId"`
	MerchantEmail            zero.String               `json:"merchantEmail"`
	MerchantName             string                    `json:"merchantName"`
	VanityURL                zero.String               `json:"vanityURL"`
	MerchantCategory         zero.String               `json:"merchantCategory"`
	CompanyName              zero.String               `json:"companyName"`
	Pic                      zero.String               `json:"pic"`
	PicOccupation            zero.String               `json:"picOccupation"`
	PicKtpFile               zero.String               `json:"picKtpFile"`
	DailyOperationalStaff    zero.String               `json:"dailyOperationalStaff"`
	StoreClosureDate         *time.Time                `json:"storeClosureDate"`
	StoreReopenDate          *time.Time                `json:"storeReopenDate"`
	StoreActiveShippingDate  *time.Time                `json:"storeActiveShippingDate"`
	MerchantAddress          zero.String               `json:"merchantAddress"`
	MerchantVillage          zero.String               `json:"merchantVillage"`
	MerchantVillageID        zero.String               `json:"merchantVillageId"`
	MerchantDistrict         zero.String               `json:"merchantDistrict"`
	MerchantDistrictID       zero.String               `json:"merchantDistrictId"`
	MerchantCity             zero.String               `json:"merchantCity"`
	MerchantCityID           zero.String               `json:"merchantCityId"`
	MerchantProvince         zero.String               `json:"merchantProvince"`
	MerchantProvinceID       zero.String               `json:"merchantProvinceId"`
	ZipCode                  zero.String               `json:"zipCode"`
	StoreAddress             zero.String               `json:"storeAddress"`
	StoreVillage             zero.String               `json:"storeVillage"`
	StoreVillageID           zero.String               `json:"storeVillageId"`
	StoreDistrict            zero.String               `json:"storeDistrict"`
	StoreDistrictID          zero.String               `json:"storeDistrictId"`
	StoreCity                zero.String               `json:"storeCity"`
	StoreCityID              zero.String               `json:"storeCityId"`
	StoreProvince            zero.String               `json:"storeProvince"`
	StoreProvinceID          zero.String               `json:"storeProvinceId"`
	StoreZipCode             zero.String               `json:"storeZipCode"`
	PhoneNumber              zero.String               `json:"phoneNumber"`
	MobilePhoneNumber        zero.String               `json:"mobilePhoneNumber"`
	AdditionalEmail          zero.String               `json:"additionalEmail"`
	MerchantDescription      zero.String               `json:"merchantDescription"`
	MerchantLogo             zero.String               `json:"merchantLogo"`
	AccountHolderName        zero.String               `json:"accountHolderName"`
	BankID                   zero.Int                  `json:"bankId"`
	BankCode                 zero.String               `json:"-"`
	BankName                 zero.String               `json:"bankName"`
	BankBranch               zero.String               `json:"bankBranch"`
	AccountNumber            zero.String               `json:"accountNumber"`
	IsPKP                    bool                      `json:"isPKP"`
	Npwp                     zero.String               `json:"npwp"`
	NpwpHolderName           zero.String               `json:"npwpHolderName"`
	NpwpFile                 zero.String               `json:"npwpFile"`
	RichContent              zero.String               `json:"richContent"`
	NotificationPreferences  zero.Int                  `json:"notificationPreferences"`
	MerchantRank             zero.String               `json:"merchantRank"`
	Acquisitor               zero.String               `json:"acquisitor"`
	AccountManager           zero.String               `json:"accountManager"`
	LaunchDev                zero.String               `json:"launchDev"`
	SkuLive                  null.Time                 `json:"skuLive"`
	MouDate                  null.Time                 `json:"mouDate"`
	Note                     zero.String               `json:"-"`
	BusinessType             zero.String               `json:"businessType"`
	Source                   zero.String               `json:"source"`
	AgreementDate            null.Time                 `json:"-"`
	DefaultBanner            zero.String               `json:"-"`
	IsActive                 bool                      `json:"isActive"`
	Status                   string                    `json:"status"`
	Version                  zero.Int                  `json:"version,omitempty"`
	IsClosed                 null.Bool                 `json:"isClosed"`
	CreatorID                null.String               `json:"creatorId"`
	CreatorIP                null.String               `json:"creatorIp,omitempty"`
	EditorID                 null.String               `json:"editorId"`
	EditorIP                 null.String               `json:"editorIp,omitempty"`
	MerchantType             MerchantType              `json:"-"`
	MerchantTypeString       zero.String               `json:"merchantType"`
	GenderPic                GenderPic                 `json:"-"`
	GenderPicString          zero.String               `json:"genderPic"`
	MerchantGroup            zero.String               `json:"merchantGroup"`
	UpgradeStatus            zero.String               `json:"upgradeStatus"`
	Created                  zero.Time                 `json:"created"`
	LastModified             null.Time                 `json:"lastModified"`
	DeletedAt                null.Time                 `json:"deletedAt"`
	Documents                []B2CMerchantDocumentData `json:"documents"`
	ProductType              zero.String               `json:"productType"`
	LegalEntity              zero.Int                  `json:"legalEntity"`
	LegalEntityName          zero.String               `json:"legalEntityName,omitempty"`
	NumberOfEmployee         zero.Int                  `json:"numberOfEmployee"`
	NumberOfEmployeeName     zero.String               `json:"numberOfEmployeeName,omitempty"`
	Maps                     Maps                      `json:"maps"`
	IsMapAvailable           bool                      `json:"isMapAvailable"`
	CountUpdateNameAvailable int                       `json:"countUpdateNameAvailable"`
	SellerOfficerName        zero.String               `json:"-"`
	SellerOfficerEmail       zero.String               `json:"-"`
	Reason                   zero.String               `json:"-"`
}
