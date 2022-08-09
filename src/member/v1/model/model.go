package model

import (
	"database/sql"
	"strings"
	"time"

	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/lib/pq"
)

const (
	Module = "member"
	// InActive status user
	InActive FgStatus = iota
	// Active status user
	Active
	// Blocked status user
	Blocked
	// New status user
	New

	// Male const variable
	Male Gender = iota
	// Female const variable
	Female
	// Secret const variable
	Secret

	// MaleString const variable
	MaleString = "MALE"
	// FemaleString const variable
	FemaleString = "FEMALE"
	// SecretString const variable
	SecretString = "SECRET"

	// ActiveString const variable
	ActiveString = "ACTIVE"
	// InactiveString const variable
	InactiveString = "INACTIVE"
	// BlockedString const variable
	BlockedString = "BLOCKED"
	// NewString const variable
	NewString = "NEW"

	// DefaultDateTime const default RFC3339 time
	DefaultDateTime = "0001-01-01T00:00:00Z"
	// DefaultDate const default date
	DefaultDate = "0001-01-01"

	// ForgotPasswordKeyRedis redis key
	ForgotPasswordKeyRedis = "forgotpassword"

	// Dolphin pack name
	Dolphin = "dolphin"
	// Starfish pack name
	Starfish = "starfish"

	//Success MFA Activation
	SuccessMFAActivation = "Success activate MFA"

	// ErrorInvalidEmailBahasa error message for invalid email
	ErrorInvalidEmailBahasa = "Alamat email tidak bisa digunakan"

	// ErrorMaxAttemptResendActivation error message for max attempt
	ErrorMaxAttemptResendActivation = "Anda telah mencapai limit harian, silakan mencoba kembali besok"

	// ErrorResendActivation error message for resend activation
	ErrorResendActivation = "Gagal mengirimkan email aktifasi, silahkan mencoba kembali"

	// ErrorResendActivationTime error message for resend activation times
	ErrorResendActivationTime = "Masih terdapat permintaan kirim ulang email aktifasi yang sebelumnya"

	// SuccessResendActivation success message for resend activation
	SuccessResendActivation = "Periksa email Anda, link konfirmasi telah dikirimkan ke email"

	// SubjectConfirmEmail for send subject confirm email
	SubjectConfirmEmail = "Informasi Akun – Konfirmasi Email Pendaftaran Bhinneka.Com"

	// SubjectWelcomeEmail for send subject welcome email
	SubjectWelcomeEmail = "Informasi Akun – Selamat Bergabung di Bhinneka.Com"

	// SubjectForgotPassword for send subject forgot password email
	SubjectForgotPassword = "Informasi Akun - Lupa Password Bhinneka.Com"

	// SubjectConfirmPassword for send subject confirm password email
	SubjectConfirmPassword = "Informasi Akun - Konfirmasi Password Baru"

	SubjectAddMember = "Informasi Akun – Pembuatan Password Bhinneka.Com"

	// Sturgeon flag for all new
	Sturgeon = "sturgeon"

	// ErrorMFAPassword error message
	ErrorMFAPassword = "Wrong Password"

	// ErrorMFAOTP error message for otp
	ErrorMFAOTP = "Verification Code Invalid"

	// ErrorEnabledMFA error message for enabled MFA
	ErrorEnabledMFA = "Failed Enable Multi Factor Authentication"

	// ErrorDisabledMFA error message for disabled MFA
	ErrorDisabledMFA = "Failed Disable Multi Factor Authentication"

	// ErrorGenerateMFA error message for generate MFA
	ErrorGenerateMFA = "MFA Generated failed"

	// ErrorTokenMFA error message for token MFA
	ErrorTokenMFA = "this token has been expired"

	// ErrorMFARequired error message for mfa required
	ErrorMFARequired = "Multi-factor authentication required"

	// StaticSharedMfaKeyForDev value of shared key static development/staging
	StaticSharedMfaKeyForDev = "ST4T1CSH4R3DK3Y1"

	// StaticOTPMfaForDev value of shared key static development/staging
	StaticOTPMfaForDev = "098765"

	// ErrorRevokeAllAccess error message for revoke login access
	ErrorRevokeAllAccess = "Failed Revoke All Login Access"

	// ErrorGetLoginActivity error message for get login activity
	ErrorGetLoginActivity = "Failed Get Login Activity"

	//ErrorOldPasswordInvalid error message for update password
	ErrorOldPasswordInvalid = "your old password is invalid"

	FieldEmail           = "email"
	FieldFirstName       = "firstName"
	FieldLastName        = "lastName"
	FieldMobile          = "mobile"
	FieldGender          = "gender"
	FieldStreet1         = "street1"
	FieldStreet2         = "street2"
	FieldPostalCode      = "postalCode"
	FieldDOB             = "dob"
	FieldSubDistrictID   = "subDistrictId"
	FieldSubDistrictName = "subDistrictName"
	FieldDistrictID      = "districtId"
	FieldDistrictName    = "districtName"
	FieldCityID          = "cityId"
	FieldCityName        = "cityName"
	FieldProvinceID      = "provinceId"
	FieldProvinceName    = "provinceName"
	FieldPhone           = "phone"
	FieldExt             = "ext"
	FieldUpdateType      = "update"
	FieldOldPassword     = "oldPassword"
	FieldNewPassword     = "newPassword"
	FieldRePassword      = "rePassword"
	FieldPassword        = "password"
	FieldStatus          = "status"
)

var (
	// SignUpFrom variable to explain register source
	SignUpFrom = []string{
		Dolphin,
		Starfish,
	}
)

// AllowedSortFields is allowed field name for sorting
var AllowedSortFields = []string{
	"id",
	"firstName",
	"lastName",
	"email",
}

// Gender is type int
type Gender int

// FgStatus initialize status data type
type FgStatus int

// Parameters data structure
type Parameters struct {
	Query    string
	StrPage  string
	Page     int
	StrLimit string
	Limit    int
	Offset   int
	Status   string
	Sort     string
	OrderBy  string
	Email    string
	IsStaff  string
	IsAdmin  string
	UserID   string
	IsActive string
}

type GetMemberResult struct {
	Result     Member
	Error      error
	HTTPStatus int
	Scope      string
}

// PlainResponse data structure
type PlainResponse struct {
	FirstName   string
	LastName    string
	HasPassword bool
}

// SuccessResponse data structure
type SuccessResponse struct {
	ID          string `jsonapi:"primary,memberType" json:"memberType,omitempty"`
	Message     string `jsonapi:"attr,message" json:"message,omitempty"`
	Token       string `jsonapi:"attr,token,omitempty" json:"token,omitempty"`
	HasPassword bool   `jsonapi:"attr,hasPassword" json:"hasPassword"`
	FirstName   string `jsonapi:"attr,firstName,omitempty" json:"firstName,omitempty"`
	LastName    string `jsonapi:"attr,lastName,omitempty" json:"lastName,omitempty"`
	Email       string `jsonapi:"attr,email,omitempty" json:"email,omitempty"`
}

// PlainSuccessResponse data structure
type PlainSuccessResponse struct {
	ID           string `jsonapi:"primary,memberType" json:"memberType,omitempty"`
	Message      string `jsonapi:"attr,message" json:"message,omitempty"`
	FirstName    string `jsonapi:"attr,firstName,omitempty" json:"firstName,omitempty"`
	LastName     string `jsonapi:"attr,lastName,omitempty" json:"lastName,omitempty"`
	Email        string `jsonapi:"attr,email,omitempty" json:"email,omitempty"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// Member data structure
type Member struct {
	ID                         string      `jsonapi:"primary,memberType" json:"id" query:"id" fieldname:"id"`
	FirstName                  string      `jsonapi:"attr,firstName" json:"firstName"  fieldname:"id" form:"firstName"`
	LastName                   string      `jsonapi:"attr,lastName" json:"lastName"  fieldname:"lastName" form:"lastName"`
	Email                      string      `jsonapi:"attr,email" json:"email"  fieldname:"email" form:"email"`
	Gender                     Gender      `json:"-"`
	GenderString               string      `jsonapi:"attr,gender" json:"gender"  fieldname:"gender" form:"gender"`
	Mobile                     string      `jsonapi:"attr,mobile" json:"mobile" query:"mobile" fieldname:"mobile" form:"mobile"`
	Phone                      string      `jsonapi:"attr,phone" json:"phone"  fieldname:"phone" form:"phone"`
	Ext                        string      `jsonapi:"attr,ext" json:"ext"  fieldname:"ext" form:"ext"`
	BirthDate                  time.Time   `json:"-"`
	BirthDateString            string      `jsonapi:"attr,birthDate" json:"birthDate"  fieldname:"id" form:"dob"`
	Password                   string      `jsonapi:"attr,password,omitempty" json:"password,omitempty" fieldname:"id" form:"password"`
	OldPassword                string      `json:"-"`
	NewPassword                string      `json:"-"`
	RePassword                 string      `json:"rePassword,omitempty" form:"rePassword"`
	Salt                       string      `jsonapi:"attr,salt,omitempty" json:"salt,omitempty"  fieldname:"salt"`
	Address                    Address     `jsonapi:"attr,address" json:"address"  fieldname:"address" form:"address"`
	JobTitle                   string      `jsonapi:"attr,jobTitle" json:"jobTitle"  fieldname:"jobTitle" form:"jobTitle"`
	Department                 string      `jsonapi:"attr,department" json:"department"  fieldname:"department" form:"department"`
	Status                     FgStatus    `json:"-"`
	StatusString               string      `jsonapi:"attr,status" json:"status"  fieldname:"status"`
	IsActive                   bool        `jsonapi:"attr,isActive" json:"isActive" fieldname:"isActive"`
	IsActiveString             string      `json:"-"`
	SocialMedia                SocialMedia `jsonapi:"attr,socialMedia" json:"socialMedia"  fieldname:"socialMedia" form:"socialMedia"`
	IsAdmin                    bool        `jsonapi:"attr,isAdmin" json:"isAdmin"  fieldname:"isAdmin"`
	IsAdminString              string      `json:"-"`
	IsStaff                    bool        `jsonapi:"attr,isStaff" json:"isStaff"  fieldname:"isStaff"`
	IsStaffString              string      `json:"-"`
	SignUpFrom                 string      `jsonapi:"attr,signUpFrom" json:"signUpFrom"  fieldname:"signUpFrom" form:"signUpFrom"`
	HasPassword                bool        `jsonapi:"attr,hasPassword" json:"hasPassword"  fieldname:"hasPassword"`
	Token                      string      `json:"token"  fieldname:"token"`
	LastLogin                  time.Time   `json:"-"`
	LastLoginString            string      `jsonapi:"attr,lastLogin" json:"lastLogin"  fieldname:"lastLogin"`
	Version                    int         `jsonapi:"attr,version,omitempty" json:"version"  fieldname:"version"`
	Created                    time.Time   `json:"-"`
	CreatedString              string      `jsonapi:"attr,created" json:"created"  fieldname:"created"`
	LastModified               time.Time   `json:"-"`
	LastModifiedString         string      `jsonapi:"attr,lastModified" json:"lastModified"  fieldname:"lastModified"`
	LastBlocked                time.Time   `json:"-"`
	LastBlockedString          string      `jsonapi:"attr,lastBlocked" json:"lastBlocked"  fieldname:"lastBlocked"`
	ProfilePicture             string      `jsonapi:"attr,profilePicture" json:"profilePicture"  fieldname:"profilePicture"`
	Type                       string      `json:"-"`
	UpdateFrom                 string      `json:"-"`
	LastTokenAttempt           time.Time   `json:"-"`
	LastTokenAttemptString     string      `jsonapi:"attr,lastTokenAttempt,omitempty" json:"lastTokenAttempt,omitempty"  fieldname:"lastTokenAttempt"`
	MFAEnabled                 bool        `jsonapi:"attr,mfaEnabled" json:"mfaEnabled"  fieldname:"mfaEnabled"`
	MFAKey                     string      `jsonapi:"attr,mfaKey" json:"mfaKey"  fieldname:"mfaKey"`
	LastPasswordModified       time.Time   `json:"-"`
	LastPasswordModifiedString string      `jsonapi:"attr,lastPasswordModified" json:"lastPasswordModified"  fieldname:"lastPasswordModified"`
	RequestFrom                string      `json:"-"`
	CreatedBy                  string      `json:"-"`
	ModifiedBy                 string      `json:"-"`
	RegisterType               string      `json:"-,omitempty" form:"registerType"`
	APIVersion                 string      `json:"-"`
	NewMember                  bool        `json:"-"`
	AdminMFAEnabled            bool        `jsonapi:"attr,mfaAdminEnabled" json:"mfaAdminEnabled"  fieldname:"mfaAdminEnabled"`
	MFAAdminKey                string      `jsonapi:"attr,mfaAdminKey" json:"mfaAdminKey"  fieldname:"mfaAdminKey"`
	IsSync                     bool        `jsonapi:"attr,isSync" json:"isSync"  fieldname:"isSync"`
}

// SocialMedia data structure
type SocialMedia struct {
	FacebookID            string    `jsonapi:"attr,facebookId" json:"facebookId"`
	FacebookConnect       time.Time `json:"-"`
	FacebookConnectString string    `jsonapi:"attr,facebookConnect" json:"facebookConnect"  fieldname:"facebookConnect"`
	GoogleID              string    `jsonapi:"attr,googleId" json:"googleId"`
	GoogleConnect         time.Time `json:"-"`
	GoogleConnectString   string    `jsonapi:"attr,googleConnect" json:"googleConnect"  fieldname:"googleConnect"`
	AppleID               string    `jsonapi:"attr,appleId" json:"appleId"`
	AppleConnect          time.Time `json:"-"`
	AppleConnectString    string    `jsonapi:"attr,appleConnect" json:"appleConnect"  fieldname:"appleConnect"`
	AzureID               string    `jsonapi:"attr,azureId" json:"azureId"`
	LDAPID                string    `jsonapi:"attr,ldapId" json:"ldapId"`
}

// Address data structure
type Address struct {
	Province      string `jsonapi:"attr,province" json:"province"`
	ProvinceID    string `jsonapi:"attr,provinceId" json:"provinceId"`
	City          string `jsonapi:"attr,city" json:"city"`
	CityID        string `jsonapi:"attr,cityId" json:"cityId"`
	District      string `jsonapi:"attr,district" json:"district"`
	DistrictID    string `jsonapi:"attr,districtId" json:"districtId"`
	SubDistrict   string `jsonapi:"attr,subDistrict" json:"subDistrict"`
	SubDistrictID string `jsonapi:"attr,subDistrictId" json:"subDistrictId"`
	ZipCode       string `jsonapi:"attr,zipCode" json:"zipCode"`
	Street1       string `jsonapi:"attr,street1" json:"street1"`
	Street2       string `jsonapi:"attr,street2" json:"street2"`
	Address       string `json:"-"`
}

// ListMembers data structure
type ListMembers struct {
	ID        string    `jsonapi:"primary,members" json:"members"`
	Name      string    `jsonapi:"attr,name" json:"name"`
	Members   []*Member `jsonapi:"relation,member" json:"member"`
	TotalData int       `json:"totalData"`
}

// Members data structure
type Members struct {
	Data []Member `json:"data"`
}

// MemberRedis data structure
type MemberRedis struct {
	ID    string        `json:"id"`
	Token string        `json:"token"`
	Count int           `json:"count"`
	TTL   time.Duration `json:"-"`
}

// MemberError data structure
type MemberError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// PayloadUpdate request structure
type PayloadUpdate struct {
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
	RePassword  string `json:"rePassword,omitempty"`
}

// User data structure
type User struct {
	Email       string        `json:"email"`
	HasPassword bool          `json:"hasPassword"`
	Users       []interface{} `json:"users"`
	Sellers     []interface{} `json:"sellers"`
	Microsites  []interface{} `json:"microsites"`
}

// ListClient data structure
type ListClient struct {
	UserType        string  `json:"userType"`
	FirstName       *string `json:"firstName"`
	LastName        *string `json:"lastName"`
	Logo            string  `json:"logo"`
	IsSync          bool    `json:"isSync"`
	TransactionType string  `json:"transactionType,omitempty"`
}

// ListClient data structure
type ListClientMerchant struct {
	SellerId                 string `json:"sellerID,omitempty"`
	UserType                 string `json:"userType"`
	FirstName                string `json:"firstName,omitempty"`
	LastName                 string `json:"lastName,omitempty"`
	Logo                     string `json:"logo"`
	IsSync                   bool   `json:"isSync"`
	MerchantServiceAvailable bool   `json:"merchantServiceAvailable"`
	VanityURL                string `json:"vanityURL"`
	IsActive                 bool   `json:"isActive"`
	IsPKP                    bool   `json:"isPKP"`
	UpgradeStatus            string `json:"upgradeStatus"`
	MerchantType             string `json:"merchantType"`
	MerchantName             string `json:"merchantName"`
}

// String function for converting gender
func (g Gender) String() string {
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

// GetDolpinGender for dolphin gender
func (g Gender) GetDolpinGender() string {
	switch g {
	case Male:
		return "M"
	case Female:
		return "F"
	case Secret:
		return "S"
	}
	return ""
}

// String function for converting user status
func (fg FgStatus) String() string {
	switch fg {
	case InActive:
		return InactiveString
	case Active:
		return ActiveString
	case Blocked:
		return BlockedString
	case New:
		return NewString
	}
	return InactiveString
}

// StringToStatus function for converting string user status to int
func StringToStatus(s string) FgStatus {
	switch strings.ToUpper(s) {
	case InactiveString:
		return InActive
	case ActiveString:
		return Active
	case BlockedString:
		return Blocked
	case NewString:
		return New
	}
	return InActive
}

// StringToGender function
func StringToGender(s string) Gender {
	switch strings.ToUpper(s) {
	case "M", MaleString:
		return Male
	case "F", FemaleString:
		return Female
	case "S", SecretString:
		return Secret
	}
	return 0
}

// IsBhinnekaEmail function for check valid bhinneka's email
func (m Member) IsBhinnekaEmail() bool {
	bhinnekaDomain := "@bhinneka.com"
	return strings.Contains(m.Email, bhinnekaDomain)
}

// ValidateGender function for validating gender
func ValidateGender(s string) (Gender, bool) {
	var g Gender
	if strings.ToUpper(s) == MaleString || strings.ToUpper(s) == "M" || strings.ToUpper(s) == FemaleString || strings.ToUpper(s) == "F" ||
		strings.ToUpper(s) == SecretString || strings.ToUpper(s) == "S" || strings.ToUpper(s) == "O" {
		g = StringToGender(s)
		return g, true
	}
	return g, false
}

// ValidateStatus function for validating status
func ValidateStatus(s string) (FgStatus, bool) {
	var st FgStatus
	if strings.ToUpper(s) == ActiveString || strings.ToUpper(s) == InactiveString || strings.ToUpper(s) == BlockedString || strings.ToUpper(s) == NewString {
		st = StringToStatus(s)
		return st, true
	}
	return st, false
}

// ProfilePicture data structure
type ProfilePicture struct {
	ID             string `json:"id"`
	ProfilePicture string `json:"profilePicture"`
}

// ProfileName data structure
type ProfileName struct {
	ID          string `json:"id"`
	ProfileName string `json:"profileName"`
}

// ResendActivationAttempt data structure
type ResendActivationAttempt struct {
	Key                        string
	Attempt                    string
	ResendActivationAttemptAge time.Duration
}

// MFASettings data structure
type MFASettings struct {
	MfaEnabled           bool      `json:"mfaEnabled"`
	LastMfaEnabled       time.Time `json:"-"`
	LastMfaEnabledString string    `json:"lastMfaEnabled"  fieldname:"lastMfaEnabled"`
}

// MFAAdminSettings specific for narwhal
type MFAAdminSettings struct {
	MfaAdminEnabled           bool      `json:"mfaAdminEnabled"`
	LastMfaAdminEnabled       time.Time `json:"-"`
	LastMfaAdminEnabledString string    `json:"lastMfaAdminEnabled"  fieldname:"lastMfaAdminEnabled"`
}

// MFAGenerateSettings data structure
type MFAGenerateSettings struct {
	SharedKeyQRCode string `json:"sharedKeyQRCode"`
	SharedKeyText   string `json:"sharedKeyText"`
}

// MFAActivateSettings data structure
type MFAActivateSettings struct {
	MemberID      string `json:"memberId"`
	SharedKeyText string `json:"sharedKeyText"`
	Otp           string `json:"otp"`
	Password      string `json:"password,omitempty"`
	RequestFrom   string
}

// SessionInfoDetail data structure
type SessionInfoDetail struct {
	ID         string     `json:"id"`
	DeviceType string     `json:"deviceType"`
	IP         string     `json:"ip"`
	UserAgent  string     `json:"userAgent"`
	LastLogin  *time.Time `json:"lastLogin"`
	ActiveNow  bool       `json:"activeNow"`
	GrantType  string     `json:"grantType"`
	IsMobile   bool       `json:"isMobile"`
	IsApp      bool       `json:"isApp"`
}

// SessionInfo data structure
type SessionInfo struct {
	ActiveSession  []SessionInfoDetail `json:"activeSession"`
	HistorySession []SessionInfoDetail `json:"historySession"`
}

// SessionInfoList data structure
type SessionInfoList struct {
	Data      SessionInfo `jsonapi:"relation,data" json:"data"`
	TotalData int         `json:"totalData"`
}

// SessionHistoryInfoList data structure
type SessionHistoryInfoList struct {
	Data      []SessionInfoDetail `jsonapi:"relation,data" json:"data"`
	TotalData int                 `json:"totalData"`
}

// ProfileComplete data structure
type ProfileComplete struct {
	Field      []ProfileField `json:"field"`
	Percentage string         `json:"percentage"`
}

// ProfileField data structure
type ProfileField struct {
	Step  int    `json:"step"`
	Key   string `json:"key"`
	Value bool   `json:"value"`
	Label string `json:"label"`
}

// ParametersLoginActivity data structure
type ParametersLoginActivity struct {
	StrPage   string `json:"strPage" form:"strPage" query:"strPage" validate:"omitempty,numeric" fieldname:"strPage" url:"strPage"`
	Page      int    `json:"page" form:"page" query:"page" validate:"omitempty,numeric" fieldname:"page" url:"page"`
	StrLimit  string `json:"strLimit" form:"strLimit" query:"strLimit" validate:"omitempty" fieldname:"strLimit" url:"strLimit"`
	Limit     int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,numeric" fieldname:"limit" url:"limit"`
	Offset    int    `json:"offset" form:"offset" query:"offset" validate:"omitempty,numeric" fieldname:"offset" url:"offset"`
	MemberID  string `json:"memberId"`
	Token     string `json:"token"`
	ExcludeID string `json:"excludeID"`
}

// MemberAdditionalInfo data structure
type MemberAdditionalInfo struct {
	ID           int
	MemberID     string
	AuthType     string
	Data         interface{}
	Created      time.Time
	LastModified time.Time
}

// DataMemberAdditionalInfo data structure
type DataMemberAdditionalInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GenericResponse generic response
type GenericResponse struct {
	Message     string `json:"message"`
	HasPassword bool   `json:"hasPassword,omitempty"`
}

// ForgotPasswordInput support json and form
type ForgotPasswordInput struct {
	Email string `json:"email" form:"email"`
}

// IsSocialMediaExist check if member have social media
func (m Member) IsSocialMediaExist() bool {
	var exist bool
	if m.SocialMedia.FacebookID != "" {
		exist = true
	} else if m.SocialMedia.AppleID != "" {
		exist = true
	} else if m.SocialMedia.GoogleID != "" {
		exist = true
	}
	return exist
}

func (m Member) SetHasPassword(input sql.NullString) {
	if input.Valid {
		m.HasPassword = true
	}
}
func (m Member) SetGender(input sql.NullString) {
	if input.Valid {
		m.Gender = StringToGender(input.String)
		m.GenderString = m.Gender.String()
	}
}

func (m Member) SetBirthDate(input pq.NullTime) {
	if input.Valid {
		m.BirthDate = input.Time
	}
}

type MemberLog struct {
	Before *Member `json:"before"`
	After  *Member `json:"after"`
}

type MemberEmailQueue struct {
	Member *Member         `json:"member"`
	Data   SuccessResponse `json:"response"`
}

type MemberPayloadEmail struct {
	Merchant *merchantModel.B2CMerchantDataV2 `json:"merchant"`
	Member   *Member                          `json:"member"`
}
