package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/user-service/helper"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

// MerchantStatus is type int
type MerchantStatus int

// MerchantType is type int
type MerchantType int

// GenderPic is type int
type GenderPic int

const (
	// PeroranganType const variable
	PeroranganType = "perorangan"
	// PerusahaanType const variable
	PerusahaanType = "perusahaan"
	// MerchantRegistrationSubject const variable for subject email
	MerchantRegistrationSubject = "Pendaftaran Toko Berhasil"
	// MerchantUpgradeSubject const variable for subject email
	MerchantUpgradeSubject = "Pengajuan Upgrade Toko Berhasil"
	// MerchantUpgradeApproved approved
	MerchantUpgradeApproved = "Pengajuan Upgrade Toko Telah Disetujui"

	MerchantActivated = "Pendaftaran Toko Disetujui"

	MerchantRegistrationRejectSubject = "Pendaftaran Toko Ditolak"
	MerchantUpgradeRejectSubject      = "Pengajuan Upgrade Toko Ditolak"

	// FailedGetMerchant const variable
	FailedGetMerchant = "Failed get merchant"

	// Regular const variable
	Regular MerchantType = iota
	// Manage const variable
	Manage
	// Associate const variable
	Associate

	// InActive status merchant
	InActive MerchantStatus = iota
	// Active status merchant
	Active
	// Deleted status merchant
	Deleted
	// New status merchant
	New

	// Inactive const variable
	InactiveString = "INACTIVE"
	// Deleted const variable
	DeletedString = "DELETED"
	// New const variable
	NewString = "NEW"

	// RegularString const variable
	RegularString = "REGULAR"
	// ManageString const variable
	ManageString = "MANAGE"
	// AssociateString const variable
	AssociateString = "ASSOCIATE"

	// Male const variable
	Male GenderPic = iota
	// Female const variable
	Female

	Secret
	// MaleString const variable
	MaleString = "MALE"
	// FemaleString const variable
	FemaleString = "FEMALE"

	SecretString = "SECRET"

	// MicroString const variable
	MicroString = "MICRO"
	// SmallString const variable
	SmallString = "SMALL"
	// MediumString const variable
	MediumString = "MEDIUM"
	// EnterpriseString const variable
	EnterpriseString = "COMPANY"

	// PendingManageString const variable
	PendingManageString = "PENDING_MANAGE"
	// PendingAssociateString const variable
	PendingAssociateString = "PENDING_ASSOCIATE"
	// RejectManageString const variable
	RejectManageString = "REJECT_MANAGE"
	// RejectAssociateString const variable
	RejectAssociateString = "REJECT_ASSOCIATE"
	// ActiveString const variable
	ActiveString = "ACTIVE"

	// ProductTypePhysic merchant only sell pyhcical goods
	ProductTypePhysic = "PHYSIC"
	// ProductTypeNonPhysic merchant only sell non-physical goods such as online ticket
	ProductTypeNonPhysic = "NON_PHYSIC"
	// ProductTypeCombined merchant sell all kind of goods
	ProductTypeCombined = "PHYSIC_NON_PHYSIC"
	// DefaultInputFormat default input
	DefaultInputFormat = time.RFC3339
	// DefaultDateFormat default date onlu
	DefaultDateFormat = "2006-01-02"
)

// StringToMerchantType function
func StringToMerchantType(s string) MerchantType {
	switch strings.ToUpper(s) {
	case RegularString:
		return Regular
	case ManageString:
		return Manage
	case AssociateString:
		return Associate
	}
	return 0
}

// String function for converting MerchantType
func (g MerchantType) String() string {
	switch g {
	case Regular:
		return RegularString
	case Manage:
		return ManageString
	case Associate:
		return AssociateString
	}
	return ""
}

// StringToGenderPic function
func StringToGenderPic(s string) GenderPic {
	switch strings.ToUpper(s) {
	case MaleString:
		return Male
	case FemaleString:
		return Female
	case SecretString:
		return Secret
	}
	return 0
}

// String function for converting GenderPic
func (g GenderPic) String() string {
	switch g {
	case Male:
		return MaleString
	case Female:
		return FemaleString
	case Secret:
		return SecretString
	}
	return ""
}

// String function for converting user status
func (ms MerchantStatus) String() string {
	switch ms {
	case InActive:
		return InactiveString
	case Active:
		return ActiveString
	case Deleted:
		return DeletedString
	case New:
		return NewString
	}
	return InactiveString
}

func StringToStatus(s string) MerchantStatus {
	switch strings.ToUpper(s) {
	case InactiveString:
		return InActive
	case ActiveString:
		return Active
	case DeletedString:
		return Deleted
	case NewString:
		return New
	}
	return InActive
}

// B2CMerchant data structure
type B2CMerchant struct {
	ID                      string  `json:"id"`
	UserID                  *string `json:"userId"`
	MerchantName            *string `json:"merchantName"`
	VanityURL               *string `json:"vanityURL"`
	MerchantCategory        *string `json:"merchantCategory"`
	CompanyName             *string `json:"companyName"`
	Pic                     *string `json:"pic"`
	PicOccupation           *string `json:"picOccupation"`
	DailyOperationalStaff   *string `json:"dailyOperationalStaff"`
	StoreClosureDate        *string `json:"storeClosureDate"`
	StoreReopenDate         *string `json:"storeReopenDate"`
	StoreActiveShippingDate *string `json:"storeActiveShippingDate"`
	MerchantAddress         *string `json:"merchantAddress"`
	MerchantVillage         *string `json:"merchantVillage"`
	MerchantDistrict        *string `json:"merchantDistrict"`
	MerchantCity            *string `json:"merchantCity"`
	MerchantProvince        *string `json:"merchantProvince"`
	ZipCode                 *string `json:"zipCode"`
	StoreAddress            *string `json:"storeAddress"`
	StoreVillage            *string `json:"storeVillage"`
	StoreDistrict           *string `json:"storeDistrict"`
	StoreCity               *string `json:"storeCity"`
	StoreProvince           *string `json:"storeProvince"`
	StoreZipCode            *string `json:"storeZipCode"`
	PhoneNumber             *string `json:"phoneNumber"`
	MobilePhoneNumber       *string `json:"mobilePhoneNumber"`
	AdditionalEmail         *string `json:"additionalEmail"`
	MerchantDescription     *string `json:"merchantDescription"`
	MerchantLogo            *string `json:"merchantLogo"`
	AccountHolderName       *string `json:"accountHolderName"`
	BankName                *string `json:"bankName"`
	AccountNumber           *string `json:"accountNumber"`
	IsPKP                   bool    `json:"isPKP"`
	Npwp                    *string `json:"npwp"`
	NpwpHolderName          *string `json:"npwpHolderName"`
	RichContent             *string `json:"richContent"`
	NotificationPreferences *int    `json:"notificationPreferences"`
	MerchantRank            *string `json:"merchantRank"`
	Acquisitor              *string `json:"acquisitor"`
	AccountManager          *string `json:"accountManager"`
	LaunchDev               *string `json:"launchDev"`
	SkuLive                 *string `json:"skuLive"`
	MouDate                 *string `json:"mouDate"`
	Note                    *string `json:"note"`
	AgreementDate           *string `json:"agreementDate"`
	IsActive                bool    `json:"isActive"`
	Status                  *string `json:"status"`
	CreatorID               *string `json:"creatorId"`
	CreatorIP               *string `json:"creatorIp"`
	EditorID                *string `json:"editorId"`
	EditorIP                *string `json:"editorIp"`
	Version                 *int64  `json:"version"`
	Created                 *string `json:"created"`
	LastModified            *string `json:"lastModified"`
	DeletedAt               *string `json:"deletedAt"`
	MerchantVillageID       *string `json:"merchantVillageId"`
	MerchantDistrictID      *string `json:"merchantDistrictId"`
	MerchantCityID          *string `json:"merchantCityId"`
	MerchantProvinceID      *string `json:"merchantProvinceId"`
	StoreVillageID          *string `json:"storeVillageId"`
	StoreDistrictID         *string `json:"storeDistrictId"`
	StoreCityID             *string `json:"storeCityId"`
	StoreProvinceID         *string `json:"storeProvinceId"`
	BankID                  *int64  `json:"bankId"`
	IsClosed                bool    `json:"isClosed"`
	MerchantEmail           *string `json:"merchantEmail"`
	BankBranch              *string `json:"bankBranch"`
	PicKtpFile              *string `json:"picKtpFile"`
	NpwpFile                *string `json:"npwpFile"`
	Source                  *string `json:"source"`
	BusinessType            *string `json:"businessType"`
	MerchantType            *string `json:"merchantType"`
	GenderPic               *string `json:"genderPic"`
	MerchantGroup           *string `json:"merchantGroup"`
	UpgradeStatus           *string `json:"upgradeStatus"`
	ProductType             *string `json:"productType"`
	LegalEntity             *int    `json:"legalEntity"`
	NumberOfEmployee        *int    `json:"numberOfEmployee"`
}

// B2CMerchantCreateInput data structure for input
type B2CMerchantCreateInput struct {
	ID                      string                     `json:"id" query:"id" validate:"required,lte=150" fieldname:"id"`
	UserID                  string                     `json:"userId" validate:"required,lte=25" fieldname:"user id"`
	MerchantEmail           string                     `json:"merchantEmail" validate:"required,email,lte=200" fieldname:"email"`
	MerchantName            string                     `json:"merchantName" validate:"required,gte=3,lte=200,is-alphanum-dot-comma-space-ampersand-dash-parenthesis" fieldname:"nama"`
	VanityURL               string                     `json:"vanityURL" validate:"omitempty,is-alphanum-dash" fieldname:"url"`
	CompanyName             string                     `json:"companyName" validate:"omitempty,gte=3,lte=200,is-alphanum-dot-comma-space-ampersand-dash-parenthesis" fieldname:"nama perusahaan"`
	IsPKP                   bool                       `json:"isPKP" validate:"omitempty,is-bool" fieldname:"pengusaha kena pajak"`
	MerchantAddress         string                     `json:"merchantAddress" validate:"omitempty,gte=3,lte=250,excludesall=<>{{}}" fieldname:"alamat perusahaan"`
	Pic                     string                     `json:"pic" validate:"omitempty,gte=3,lte=200,is-alphanum-dot-comma-space" fieldname:"nama pic"`
	PicOccupation           string                     `json:"picOccupation" validate:"omitempty,gte=3,lte=100,is-alphanum-dot-comma-space" fieldname:"pekerjaan pic"`
	PicKtpFile              string                     `json:"picKtpFile" validate:"omitempty,is-s3-url" fieldname:"pic KTP file"`
	PhoneNumber             string                     `json:"phoneNumber" validate:"omitempty,numeric,gte=5,lte=15" fieldname:"nomor telepon"`
	MobilePhoneNumber       string                     `json:"mobilePhoneNumber" validate:"omitempty,gte=5,lte=13,numeric,check-first-two-mobile,check-first-three-mobile" fieldname:"nomor hanphone"`
	MerchantDescription     string                     `json:"merchantDescription" validate:"omitempty,lte=255,excludesall=<>" fieldname:"deskripsi merchant"`
	AccountHolderName       string                     `json:"accountHolderName" form:"accountHolderName" query:"accountHolderName" validate:"omitempty,gte=3,lte=200" fieldname:"nama pemilik rekening"`
	BankID                  int32                      `json:"bankId" form:"bankId" query:"bankId" validate:"omitempty,numeric" fieldname:"bank id"`
	BankBranch              string                     `json:"bankBranch" form:"bankBranch" query:"bankBranch" validate:"omitempty" fieldname:"bank branch"`
	AccountNumber           string                     `json:"accountNumber" form:"accountNumber" query:"accountNumber" validate:"omitempty,numeric" fieldname:"nomor rekening"`
	Npwp                    string                     `json:"npwp" form:"npwp" query:"npwp" validate:"omitempty,len=15,numeric" fieldname:"nomor NPWP"`
	NpwpHolderName          string                     `json:"npwpHolderName" form:"npwpHolderName" query:"npwpHolderName" validate:"omitempty,gte=3,lte=200,is-alphanum-dot-comma-space-ampersand-dash-parenthesis" fieldname:"nama pemegang NPWP"`
	NpwpFile                string                     `json:"npwpFile" form:"npwpFile" query:"npwpFile" validate:"omitempty,is-s3-url" fieldname:"NPWP File"`
	DailyOperationalStaff   string                     `json:"dailyOperationalStaff" form:"dailyOperationalStaff" query:"dailyOperationalStaff" validate:"omitempty" fieldname:"nama operational staff"`
	StoreClosureDate        string                     `json:"storeClosureDate" form:"storeClosureDate" query:"storeClosureDate" validate:"omitempty,is-store-date-format" fieldname:"tanggal tutup"`
	StoreReopenDate         string                     `json:"storeReopenDate" form:"storeReopenDate" query:"storeReopenDate" validate:"omitempty,is-store-date-format" fieldname:"tanggal buka"`
	StoreActiveShippingDate string                     `json:"storeActiveShippingDate" form:"storeActiveShippingDate" query:"storeActiveShippingDate" validate:"omitempty,is-store-date-format" fieldname:"tanggal shipping aktif"`
	MerchantCategory        string                     `json:"merchantCategory" query:"merchantCategory" validate:"omitempty" fieldname:"kategori"`
	MerchantVillageID       string                     `json:"merchantVillageId" query:"merchantVillageId" validate:"omitempty,numeric" fieldname:"kelurahan perusahaan"`
	MerchantVillage         string                     `json:"merchantVillage" query:"merchantVillageId" validate:"omitempty,numeric" fieldname:"kelurahan perusahaan"`
	MerchantDistrictID      string                     `json:"merchantDistrictId" query:"merchantDistrictId" validate:"omitempty,numeric" fieldname:"kecamatan perusahaan"`
	MerchantDistrict        string                     `json:"merchantDistrict" query:"merchantDistrictId" validate:"omitempty,numeric" fieldname:"kecamatan perusahaan"`
	MerchantCityID          string                     `json:"merchantCityId" query:"merchantCityId" validate:"omitempty,numeric" fieldname:"kota perusahaan"`
	MerchantCity            string                     `json:"merchantCity" query:"merchantCityId" validate:"omitempty,numeric" fieldname:"kota perusahaan"`
	MerchantProvinceID      string                     `json:"merchantProvinceId" query:"merchantProvinceId" validate:"omitempty,numeric" fieldname:"provinsi perusahaan"`
	MerchantProvince        string                     `json:"merchantProvince" query:"merchantProvinceId" validate:"omitempty,numeric" fieldname:"provinsi perusahaan"`
	ZipCode                 string                     `json:"zipCode" query:"zipCode" validate:"omitempty,numeric" fieldname:"zipcode perusahaan"`
	StoreAddress            string                     `json:"storeAddress" query:"storeAddress" validate:"omitempty,gte=3,lte=250" fieldname:"alamat toko"`
	StoreVillageID          string                     `json:"storeVillageId" query:"storeVillageId" validate:"omitempty,numeric" fieldname:"kelurahan toko"`
	StoreVillage            string                     `json:"storeVillage" query:"storeVillage" validate:"omitempty,numeric" fieldname:"kelurahan toko"`
	StoreDistrictID         string                     `json:"storeDistrictId" query:"storeDistrictId" validate:"omitempty,numeric" fieldname:"kecamatan toko"`
	StoreDistrict           string                     `json:"storeDistrict" query:"storeDistrict" validate:"omitempty,numeric" fieldname:"kecamatan toko"`
	StoreCityID             string                     `json:"storeCityId" query:"storeCityId" validate:"omitempty,numeric" fieldname:"kota toko"`
	StoreCity               string                     `json:"storeCity" query:"storeCity" validate:"omitempty,numeric" fieldname:"kota toko"`
	StoreProvinceID         string                     `json:"storeProvinceId" query:"storeProvinceId" validate:"omitempty,numeric" fieldname:"provinsi toko"`
	StoreProvince           string                     `json:"storeProvince" query:"storeProvince" validate:"omitempty,numeric" fieldname:"provinsi toko"`
	StoreZipCode            string                     `json:"storeZipCode" query:"storeZipCode" validate:"omitempty,numeric" fieldname:"zipcode toko"`
	AdditionalEmail         string                     `json:"additionalEmail" query:"additionalEmail" validate:"omitempty,lte=250,email" fieldname:"additional email"`
	MerchantLogo            string                     `json:"merchantLogo" query:"merchantLogo" validate:"omitempty,is-s3-url" fieldname:"logo merchant"`
	RichContent             string                     `json:"richContent" query:"richContent" validate:"omitempty,gte=15,lte=15,is-alphanum-dot-comma-space-ampersand-dash" fieldname:"rich content"`
	NotificationPreferences int                        `json:"notificationPreferences" query:"notificationPreferences" validate:"omitempty,numeric" fieldname:"notification preference"`
	MerchantRank            string                     `json:"merchantRank" query:"merchantRank" validate:"omitempty" fieldname:"rank merchant"`
	Acquisitor              string                     `json:"acquisitor" query:"acquisitor" validate:"omitempty,lte=50" fieldname:"acquisitor"`
	AccountManager          string                     `json:"accountManager" query:"accountManager" validate:"omitempty,lte=50" fieldname:"manajer akun"`
	LaunchDev               string                     `json:"launchDev" query:"launchDev" validate:"omitempty,lte=50" fieldname:"launch dev"`
	SkuLive                 string                     `json:"skuLive" query:"skuLive" validate:"omitempty,is-store-date-format" fieldname:"SKU live"`
	MouDate                 string                     `json:"mouDate" validate:"omitempty,is-store-date-format" fieldname:"tanggal MOU"`
	AgreementDate           string                     `json:"agreementDate" validate:"omitempty,is-store-date-format" fieldname:"tanggal aggreement"`
	Note                    string                     `json:"note" validate:"omitempty" fieldname:"catatan"`
	BusinessType            string                     `json:"businessType" validate:"omitempty" fieldname:"Jenis Usaha"`
	Source                  string                     `json:"source" validate:"omitempty" fieldname:"Source"`
	IsActive                bool                       `json:"isActive" validate:"omitempty,is-bool" fieldname:"toko aktif"`
	Status                  string                     `json:"status" validate:"omitempty"`
	IsClosed                bool                       `json:"isClosed" validate:"omitempty,is-bool" fieldname:"toko tutup"`
	Documents               []B2CMerchantDocumentInput `json:"documents" validate:"omitempty" fieldname:"documents"`
	MerchantType            MerchantType               `json:"-"`
	MerchantTypeString      string                     `json:"merchantType" validate:"omitempty"`
	GenderPic               GenderPic                  `json:"-"`
	GenderPicString         string                     `json:"genderPic" validate:"omitempty"`
	MerchantGroup           string                     `json:"merchantGroup,omitempty"`
	UpgradeStatus           string                     `json:"upgradeStatus" validate:"omitempty"`
	ProductType             string                     `json:"productType" validate:"omitempty" query:"productType" fieldname:"productType"`
	LegalEntity             int                        `json:"legalEntity" validate:"omitempty"`
	NumberOfEmployee        int                        `json:"numberOfEmployee" validate:"omitempty"`
	Maps                    Maps                       `json:"maps"`
	SellerOfficerName       string                     `json:"sellerOfficerName,omitempty"`
	SellerOfficerEmail      string                     `json:"sellerOfficerEmail,omitempty"`
	Token                   string
}

// B2CMerchantData data structure
type B2CMerchantData struct {
	ID                      string                    `json:"id" db:"id"`
	UserID                  string                    `json:"userId"`
	MerchantEmail           string                    `json:"merchantEmail"`
	MerchantName            string                    `json:"merchantName"`
	VanityURL               string                    `json:"vanityURL"`
	MerchantCategory        string                    `json:"merchantCategory"`
	CompanyName             string                    `json:"companyName"`
	Pic                     string                    `json:"pic"`
	PicOccupation           string                    `json:"picOccupation"`
	PicKtpFile              string                    `json:"picKtpFile"`
	DailyOperationalStaff   string                    `json:"dailyOperationalStaff"`
	StoreClosureDate        *time.Time                `json:"storeClosureDate"`
	StoreReopenDate         *time.Time                `json:"storeReopenDate"`
	StoreActiveShippingDate *time.Time                `json:"storeActiveShippingDate"`
	MerchantAddress         string                    `json:"merchantAddress"`
	MerchantVillage         string                    `json:"merchantVillage"`
	MerchantVillageID       string                    `json:"merchantVillageId"`
	MerchantDistrict        string                    `json:"merchantDistrict"`
	MerchantDistrictID      string                    `json:"merchantDistrictId"`
	MerchantCity            string                    `json:"merchantCity"`
	MerchantCityID          string                    `json:"merchantCityId"`
	MerchantProvince        string                    `json:"merchantProvince"`
	MerchantProvinceID      string                    `json:"merchantProvinceId"`
	ZipCode                 string                    `json:"zipCode"`
	StoreAddress            string                    `json:"storeAddress"`
	StoreVillage            string                    `json:"storeVillage"`
	StoreVillageID          string                    `json:"storeVillageId"`
	StoreDistrict           string                    `json:"storeDistrict"`
	StoreDistrictID         string                    `json:"storeDistrictId"`
	StoreCity               string                    `json:"storeCity"`
	StoreCityID             string                    `json:"storeCityId"`
	StoreProvince           string                    `json:"storeProvince"`
	StoreProvinceID         string                    `json:"storeProvinceId"`
	StoreZipCode            string                    `json:"storeZipCode"`
	PhoneNumber             string                    `json:"phoneNumber"`
	MobilePhoneNumber       string                    `json:"mobilePhoneNumber"`
	AdditionalEmail         string                    `json:"additionalEmail"`
	MerchantDescription     string                    `json:"merchantDescription"`
	MerchantLogo            string                    `json:"merchantLogo"`
	AccountHolderName       string                    `json:"accountHolderName"`
	BankID                  int32                     `json:"bankId"`
	BankCode                string                    `json:"bankCode,omitempty"`
	BankName                string                    `json:"bankName"`
	BankBranch              string                    `json:"bankBranch"`
	AccountNumber           string                    `json:"accountNumber"`
	IsPKP                   bool                      `json:"isPKP"`
	Npwp                    string                    `json:"npwp"`
	NpwpHolderName          string                    `json:"npwpHolderName"`
	NpwpFile                string                    `json:"npwpFile"`
	RichContent             string                    `json:"richContent"`
	NotificationPreferences int                       `json:"notificationPreferences"`
	MerchantRank            string                    `json:"merchantRank"`
	Acquisitor              string                    `json:"acquisitor"`
	AccountManager          string                    `json:"accountManager"`
	LaunchDev               string                    `json:"launchDev"`
	SkuLive                 *time.Time                `json:"skuLive"`
	MouDate                 *time.Time                `json:"mouDate"`
	Note                    string                    `json:"note"`
	BusinessType            string                    `json:"businessType"`
	Source                  string                    `json:"source"`
	AgreementDate           *time.Time                `json:"agreementDate"`
	DefaultBanner           string                    `json:"defaultBanner"`
	IsActive                bool                      `json:"isActive"`
	Status                  string                    `json:"status"`
	Version                 int                       `json:"version,omitempty"`
	IsClosed                bool                      `json:"isClosed"`
	CreatorID               string                    `json:"creatorId"`
	CreatorIP               string                    `json:"creatorIp,omitempty"`
	EditorID                string                    `json:"editorId"`
	EditorIP                string                    `json:"editorIp,omitempty"`
	MerchantType            MerchantType              `json:"-"`
	MerchantTypeString      string                    `json:"merchantType"`
	GenderPic               GenderPic                 `json:"-"`
	GenderPicString         string                    `json:"genderPic"`
	MerchantGroup           string                    `json:"merchantGroup"`
	UpgradeStatus           string                    `json:"upgradeStatus"`
	Created                 time.Time                 `json:"created"`
	LastModified            time.Time                 `json:"lastModified"`
	DeletedAt               *time.Time                `json:"deletedAt"`
	Documents               []B2CMerchantDocumentData `json:"documents"`
	ProductType             string                    `json:"productType"`
}

// MerchantUserAttribute data structure
type MerchantUserAttribute struct {
	UserID string
	UserIP string
}

// CheckMerchantName data
type CheckMerchantName struct {
	MerchantName string `json:"merchantName"`
}

// Maps data
type Maps struct {
	ID           string  `json:"id,omitempty"`
	RelationID   string  `json:"relationId"`
	RelationName string  `json:"relationName"`
	Label        string  `json:"label"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

// ResponseMerchantName data structure
type ResponseMerchantName struct {
	MerchantName string `json:"merchantName"`
	Slug         string `json:"slug"`
	URL          string `json:"url"`
}

// ResponseAvailable data structure
type ResponseAvailable struct {
	MerchantServiceAvailable bool              `json:"merchantServiceAvailable"`
	MerchantData             B2CMerchantDataV2 `json:"merchantData"`
}

// ValidateMerchantType function for validating merchantType
func ValidateMerchantType(s string) (MerchantType, bool) {
	var g MerchantType
	if strings.ToUpper(s) == RegularString || strings.ToUpper(s) == ManageString ||
		strings.ToUpper(s) == AssociateString {
		g = StringToMerchantType(s)
		return g, true
	}
	return g, false
}

// ValidateGenderPic function for validating genderPic
func ValidateGenderPic(s string) (GenderPic, bool) {
	var g GenderPic
	if strings.ToUpper(s) == MaleString || strings.ToUpper(s) == FemaleString || strings.ToUpper(s) == SecretString {
		g = StringToGenderPic(s)
		return g, true
	}
	return g, false
}

// ValidateMerchantGroup function for validating genderPic
func ValidateMerchantGroup(s string) (string, bool) {
	if strings.ToUpper(s) == MicroString || strings.ToUpper(s) == SmallString ||
		strings.ToUpper(s) == MediumString || strings.ToUpper(s) == EnterpriseString {
		return strings.ToUpper(s), true
	}
	return strings.ToUpper(s), false
}

// ValidateUpgradeStatus function for validating upgradeStatus
func ValidateUpgradeStatus(s string) (string, bool) {
	if strings.ToUpper(s) == PendingAssociateString || strings.ToUpper(s) == PendingManageString || strings.ToUpper(s) == ActiveString {
		return strings.ToUpper(s), true
	}
	return strings.ToUpper(s), false
}

// ValidateProductType function for validating upgradeStatus
func ValidateProductType(input string) bool {
	return golib.StringInSlice(strings.ToUpper(input), []string{ProductTypePhysic, ProductTypeNonPhysic, ProductTypeCombined})
}

// ValidateStatus function for validating merchantStatus
func ValidateMerchantStatus(s string) bool {
	return golib.StringInSlice(strings.ToUpper(s), []string{ActiveString, InactiveString, DeletedString, NewString})
}

type Param struct {
	IsAttachment string `json:"isAttachment,omitempty"`
}

// QueryParameters for search
type QueryParameters struct {
	Page          int
	Limit         int
	Offset        int
	StrPage       string `json:"page" query:"page"`
	StrLimit      string `json:"limit" query:"limit"`
	OrderBy       string `json:"orderBy" query:"orderBy"`
	SortBy        string `json:"sortBy" query:"sortBy"`
	Search        string `json:"search" query:"search"`
	Status        string `json:"status" query:"status"`
	IsPKP         string `json:"isPkp" query:"isPkp"`
	BusinessType  string `json:"businessType" query:"businessType"`
	MerchantType  string `json:"merchantType" query:"merchantType"`   // support comma delimited
	UpgradeStatus string `json:"upgradeStatus" query:"upgradeStatus"` // support comma delimited
	MerchantIDS   string `json:"merchantIds" query:"merchantids"`
	Name          string `json:"Name" query:"Name"`
	Email         string `json:"Email" query:"Email"`
	Officer       string `json:"officer" query:"officer"` //support comma
}

type QueryParametersPublic struct {
	Page          int
	Limit         int
	Offset        int
	StrPage       string `json:"page" query:"page"`
	StrLimit      string `json:"limit" query:"limit"`
	OrderBy       string `json:"orderBy" query:"orderBy"`
	SortBy        string `json:"sortBy" query:"sortBy"`
	Search        string `json:"search" query:"search"`
	Status        string `json:"-" query:"-"`
	IsPKP         string `json:"-" query:"-"`
	BusinessType  string `json:"-" query:"-"`
	MerchantType  string `json:"-" query:"-"` // support comma delimited
	UpgradeStatus string `json:"-" query:"-"` // support comma delimited
	MerchantIDS   string `json:"merchantIds" query:"merchantids"`
	Name          string `json:"-" query:"-"`
	Email         string `json:"-" query:"-"`
	Officer       string `json:"-" query:"-"` //support comma
}

var (
	allowedOrder       = []string{"created", "companyName", "status", ""}
	allowedOrderPublic = []string{"created", "LastModified", "merchantEmail", "picKtpFile", "phoneNumber", "mobilePhoneNumber", "additionalEmail",
		"accountHolderName", "bankId", "bankCode", "bankName", "bankBranch", "accountNumber", "npwp", "npwpHolderName", "npwpFile", "source", "agreementDate", "version",
		"creatorId", "creatorIp", "editorId", "editorIp", "lastModified", "deletedAt", "documents", "legalEntity", "legalEntityName", "numberOfEmployee", "numberOfEmployeeName",
		"maps", "isMapAvailable", "countUpdateNameAvailable", "sellerOfficerName", "sellerOfficerEmail"}
	allowedOrderMerchantEmployee = []string{"id", "createdAt", "modifiedAt", "status", ""}
	allowedSort                  = []string{"asc", "desc", ""}
	allowedBool                  = []string{"true", "false", ""}
	allowedBusinessType          = []string{"perorangan", "perusahaan", ""}
	allowedUpgradeStatus         = []string{"PENDING_MANAGE", "PENDING_ASSOCIATE", "ACTIVE", "REJECT_MANAGE", "REJECT_ASSOCIATE", ""}
	allowedMerchantType          = []string{"regular", "associate", "manage", ""}
	allowedStatusEmpolyee        = []string{helper.TextInvited, helper.TextActive, helper.TextInactive, helper.TextRevoked, ""}
)

const delimiter = ", "

// Validate allowed input
func (q *QueryParameters) Validate() error {
	if !golib.StringInSlice(q.OrderBy, allowedOrder, false) {
		return fmt.Errorf("orderBy must be one of %s", strings.Join(allowedOrder[:2], delimiter))
	}
	if !golib.StringInSlice(q.SortBy, allowedSort, false) {
		return fmt.Errorf("sortBy must be one of %s", strings.Join(allowedSort[:2], delimiter))
	}
	if !golib.StringInSlice(q.Status, allowedBool, false) {
		return fmt.Errorf("status must be one of %s", strings.Join(allowedBool[:2], delimiter))
	}
	if !golib.StringInSlice(q.IsPKP, allowedBool, false) {
		return fmt.Errorf("isPkp must be one of %s", strings.Join(allowedBool[:2], delimiter))
	}
	if !golib.StringInSlice(q.BusinessType, allowedBusinessType, false) {
		return fmt.Errorf("businessType must be one of %s", strings.Join(allowedBusinessType[:2], delimiter))
	}
	if err := q.validateStatus(); err != nil {
		return err
	}

	if err := q.validateMerchantType(); err != nil {
		return err
	}

	return nil
}

func (q *QueryParametersPublic) ValidatePublic() error {
	if golib.StringInSlice(q.OrderBy, allowedOrderPublic, false) {
		return fmt.Errorf("orderBy is invalid for public use")
	}
	if !golib.StringInSlice(q.SortBy, allowedSort, false) {
		return fmt.Errorf("sortBy must be one of %s", strings.Join(allowedSort[:2], delimiter))
	}
	return nil
}
func (q *QueryParameters) validateStatus() error {
	if q.UpgradeStatus != "" {
		statuses := strings.Split(q.UpgradeStatus, ",")
		for _, s := range statuses {
			ss := helper.TrimSpace(s)
			if !golib.StringInSlice(ss, allowedUpgradeStatus, false) {
				return fmt.Errorf("upgradeStatus must be any of of %s", strings.Join(allowedUpgradeStatus[:3], delimiter))
			}
		}
	}
	return nil
}

func (q *QueryParameters) validateMerchantType() error {
	if q.MerchantType != "" {
		merchantTypes := strings.Split(q.MerchantType, ",")
		for _, types := range merchantTypes {
			nTypes := helper.TrimSpace(types)
			if !golib.StringInSlice(nTypes, allowedMerchantType, false) {
				return fmt.Errorf("merchantType must be any of of %s", strings.Join(allowedMerchantType[:3], delimiter))
			}
		}
	}
	return nil
}

// B2CMerchantDataV2 data structure
type B2CMerchantDataV2 struct {
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
	BankCode                 zero.String               `json:"bankCode,omitempty"`
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
	Note                     zero.String               `json:"note"`
	BusinessType             zero.String               `json:"businessType"`
	Source                   zero.String               `json:"source"`
	AgreementDate            null.Time                 `json:"agreementDate"`
	DefaultBanner            zero.String               `json:"defaultBanner"`
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
	SellerOfficerName        zero.String               `json:"sellerOfficerName"`
	SellerOfficerEmail       zero.String               `json:"sellerOfficerEmail"`
	Reason                   zero.String               `json:"reason,omitempty"`
}

// LegalEntity legal entity
type LegalEntity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CompanySize company size entity
type CompanySize struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MerchantPayloadEmail struct {
	MemberName   string            `json:"memberName"`
	Data         B2CMerchantDataV2 `json:"merchant"`
	ReasonReject string            `json:"reasonReject"`
	AdminCMS     string            `json:"adminCMS"`
}

type MerchantLog struct {
	Before B2CMerchantDataV2 `json:"before"`
	After  B2CMerchantDataV2 `json:"after"`
}

type B2CMerchantDataPublic struct {
	ID                       string                    `json:"id" db:"id"`
	UserID                   string                    `json:"userId"`
	MerchantEmail            zero.String               `json:"-"`
	MerchantName             string                    `json:"merchantName"`
	VanityURL                zero.String               `json:"vanityURL"`
	MerchantCategory         zero.String               `json:"merchantCategory"`
	CompanyName              zero.String               `json:"companyName"`
	Pic                      zero.String               `json:"pic"`
	PicOccupation            zero.String               `json:"picOccupation"`
	PicKtpFile               zero.String               `json:"-"`
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
	PhoneNumber              zero.String               `json:"-"`
	MobilePhoneNumber        zero.String               `json:"-"`
	AdditionalEmail          zero.String               `json:"-"`
	MerchantDescription      zero.String               `json:"merchantDescription"`
	MerchantLogo             zero.String               `json:"merchantLogo"`
	AccountHolderName        zero.String               `json:"-"`
	BankID                   zero.Int                  `json:"-"`
	BankCode                 zero.String               `json:"-"`
	BankName                 zero.String               `json:"-"`
	BankBranch               zero.String               `json:"-"`
	AccountNumber            zero.String               `json:"-"`
	IsPKP                    bool                      `json:"isPKP"`
	Npwp                     zero.String               `json:"-"`
	NpwpHolderName           zero.String               `json:"-"`
	NpwpFile                 zero.String               `json:"-"`
	RichContent              zero.String               `json:"richContent"`
	NotificationPreferences  zero.Int                  `json:"notificationPreferences"`
	MerchantRank             zero.String               `json:"merchantRank"`
	Acquisitor               zero.String               `json:"acquisitor"`
	AccountManager           zero.String               `json:"accountManager"`
	LaunchDev                zero.String               `json:"launchDev"`
	SkuLive                  null.Time                 `json:"skuLive"`
	MouDate                  null.Time                 `json:"mouDate"`
	Note                     zero.String               `json:"note"`
	BusinessType             zero.String               `json:"businessType"`
	Source                   zero.String               `json:"-"`
	AgreementDate            null.Time                 `json:"-"`
	DefaultBanner            zero.String               `json:"defaultBanner"`
	IsActive                 bool                      `json:"isActive"`
	Status                   string                    `json:"status"`
	Version                  zero.Int                  `json:"-"`
	IsClosed                 null.Bool                 `json:"isClosed"`
	CreatorID                null.String               `json:"-"`
	CreatorIP                null.String               `json:"-"`
	EditorID                 null.String               `json:"-"`
	EditorIP                 null.String               `json:"-"`
	MerchantType             MerchantType              `json:"-"`
	MerchantTypeString       zero.String               `json:"merchantType"`
	GenderPic                GenderPic                 `json:"-"`
	GenderPicString          zero.String               `json:"genderPic"`
	MerchantGroup            zero.String               `json:"merchantGroup"`
	UpgradeStatus            zero.String               `json:"-"`
	Created                  zero.Time                 `json:"-"`
	LastModified             null.Time                 `json:"-"`
	DeletedAt                null.Time                 `json:"-"`
	Documents                []B2CMerchantDocumentData `json:"-"`
	ProductType              zero.String               `json:"productType"`
	LegalEntity              zero.Int                  `json:"-"`
	LegalEntityName          zero.String               `json:"-"`
	NumberOfEmployee         zero.Int                  `json:"-"`
	NumberOfEmployeeName     zero.String               `json:"-"`
	Maps                     Maps                      `json:"-"`
	IsMapAvailable           bool                      `json:"-"`
	CountUpdateNameAvailable int                       `json:"-"`
	SellerOfficerName        zero.String               `json:"-"`
	SellerOfficerEmail       zero.String               `json:"-"`
	Reason                   zero.String               `json:"-"`
}

func (input B2CMerchantDataV2) RestructForPublic() (res B2CMerchantDataPublic) {
	res = B2CMerchantDataPublic(input)
	return res
}

func (input QueryParametersPublic) RestructParamPublic() (res QueryParameters) {
	res = QueryParameters(input)
	return res
}

func (mr B2CMerchantCreateInput) IsBhinnekaEmail() bool {
	bhinnekaDomain := "@bhinneka.com"
	return strings.Contains(mr.SellerOfficerEmail, bhinnekaDomain)
}
