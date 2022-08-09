package model

import (
	"time"

	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
)

const (
	// Module for module name
	Module = "Auth"
	// AuthTypeAnonymous for anonymous authentication
	AuthTypeAnonymous = "anonymous"
	// AuthTypePassword for user authentication with password
	AuthTypePassword = "password"
	// AuthTypeFacebook for user authentication with facebook
	AuthTypeFacebook = "facebook"
	// AuthTypeGoogle for user authentication with google
	AuthTypeGoogle = "google"
	// AuthTypeGoogleOAauth for user authentication with google one tap login
	AuthTypeGoogleOAauth = "google_oauth"
	// AuthTypeAzure for user authentication with azure
	AuthTypeAzure = "azure"
	// AuthTypeApple for user authentication with apple
	AuthTypeApple = "apple"
	// AuthTypeRefreshToken for refreshing token
	AuthTypeRefreshToken = "refreshtoken"
	// AuthTypeLDAP for authentication with LDAP
	AuthTypeLDAP = "ldap"
	// AuthTypeVerifyMFA for authentication with Multi Factor Authentication
	AuthTypeVerifyMFA = "mfaotp"
	// AuthTypeVerifyMFANarwhal specific for narwhal
	AuthTypeVerifyMFANarwhal = "mfaotp-narwhal"
	// AuthTypeGoogleBackend bypass incognito
	AuthTypeGoogleBackend = "google-backend"

	// Bhinneka for issuer
	Bhinneka = "bhinneka.com"
	// ErrorAzureToken error message for azure token
	ErrorAzureToken = "The provided authorization code or refresh token is expired"

	// ErrorAzureInvalidRedirectURL error invalid redirect url from azure
	ErrorAzureInvalidRedirectURL = "does not match the reply address"

	// ErrorInvalidRedirectURL error invalid redirect url
	ErrorInvalidRedirectURL = "invalid redirect url"

	// ErrorAzureTokenBahasa error message for azure in bahasa
	ErrorAzureTokenBahasa = "Sesi login Anda telah berakhir, silakan login kembali"

	// ErrorAccountInActiveBahasa error message for account with inactive status
	ErrorAccountInActiveBahasa = "Akun Anda belum aktif. Periksa email subject Konfirmasi Alamat Email Anda dari Bhinneka untuk aktifasi akun. Atau kirimkan ulang email aktifasi"

	// ErrorAccountNewBahasa error message for new account that not active
	ErrorNewAccountBahasa = "Akun Anda baru terdaftar dan belum aktif. Periksa email subject Konfirmasi Alamat Email Anda dari Bhinneka untuk aktifasi akun. Atau kirimkan ulang email aktifasi"

	// ErrorAccountBlockedBahasa error message for blocked account
	ErrorAccountBlockedBahasa = "Akun Anda telah diblokir. Silakan coba kembali setelah 5 menit"

	// ErrorAccountDeactiveBahasa error message for deactive account
	ErrorAccountDeactiveBahasa = "Akun Anda telah dinonaktifkan"

	// ErrorInvalidUsernameOrPasswordBahasa error message for invalid username or password
	ErrorInvalidUsernameOrPasswordBahasa = "Login error! Kata sandi salah"

	//ErrorEmailUnregisteredBahasa error message for email not found
	ErrorEmailUnregisteredBahasa = "Maaf, Email anda belum terdaftar, silahkan registrasi terlebih dahulu"

	//ErrorEmailAlreadyRegisteredUsingSocialLogin error message for email already registered using social login
	ErrorEmailAlreadyRegisteredUsingSocialLogin = "Alamat email yang Anda masukan sudah terdaftar untuk %s login, silakan login dengan akun %s Anda untuk mengakses Bhinneka.com."

	//ErrorRefreshToken error message
	ErrorRefreshToken = "refresh token is invalid"

	//ErrorGetToken error message
	ErrorGetToken = "failed to get token"

	//ErrorOldToken error message
	ErrorOldToken = "invalid old token"

	// UserTypeCorporate label for user corporate
	UserTypeCorporate = "corporate"

	// UserTypeCorporate label for user corporate
	UserTypeMicrositeBela = "MICROSITE_BELA"

	// UserTypePersonal label for user personal
	UserTypePersonal = "personal"

	// UserTypeMerchant label for user personal
	UserTypeMerchant = "seller"

	// UserTypeMicrosite label for user microsite
	UserTypeMicrosite = "microsite"

	// MFATokenKeyRedis redis key
	MFATokenKeyRedis = "mfa-otp"

	// DefaultSubject subject claim
	DefaultSubject = "bhinneka-microservices-b13714-5312115"

	// DefaultDeviceLogin device login claim
	DefaultDeviceLogin = "WEB"

	// DefaultDeviceID device id claim
	DefaultDeviceID = "user-service"

	// UserTypeClientMicrosite used for bela, lkpp
	UserTypeClientMicrosite = "clientMicrosite"

	// ErrorUserLKPPBelaNotFoundBahasa error message ketika user belum terdaftar
	ErrorUserLKPPBelaNotFoundBahasa = "Akun anda belum terdaftar, silakan menghubungi pihak BELA LKPP untuk melanjutkan proses pendaftaran akun."
	// NarwhalMFATokenKeyRedis specific for narwhal
	NarwhalMFATokenKeyRedis = "mfa-otp-admin"
	// ErrorIncorrectMemberTypeMicrosite specific for microsite
	ErrorIncorrectMemberTypeMicrosite = "akun tidak terdaftar sebagai pengguna %s"

	LoginTypeShopcart = "shopcart"
)

// RequestToken data structure
type RequestToken struct {
	GrantType         string `json:"grantType,omitempty" form:"grantType"`
	Audience          string
	UserID            string
	Email             string `json:"email,omitempty" form:"email"`
	Username          string `json:"username,omitempty" form:"username"`
	FirstName         string `json:"firstName,omitempty" form:"firstName"`
	LastName          string `json:"lastName,omitempty" form:"lastName"`
	FullName          string `json:"fullName,omitempty"`
	Password          string `json:"password,omitempty" form:"password"`
	Code              string `json:"code,omitempty" form:"code"`
	IP                string `json:"ip,omitempty"`
	UserAgent         string `json:"userAgent,omitempty"`
	DeviceID          string `json:"deviceId,omitempty" form:"deviceId"`
	DeviceLogin       string `json:"deviceLogin,omitempty" form:"deviceLogin"`
	NewMember         bool
	HasPassword       bool
	Token             string
	OldToken          string `json:"oldToken,omitempty" form:"oldToken"`
	RefreshToken      string `json:"refreshToken,omitempty" form:"refreshToken"`
	RedirectURI       string `json:"redirectUri,omitempty" form:"redirectUri"`
	ClientID          string `json:"clientId,omitempty" form:"clientId"`
	ClientSecret      string
	ExpiredAt         time.Time
	MemberType        string `json:"memberType,omitempty" form:"memberType"`
	Department        string `json:"department,omitempty"`
	JobTitle          string `json:"jobTitle,omitempty"`
	MFAEnabled        bool   `json:"mfaEnabled"`
	MFAToken          string `json:"mfaToken,omitempty" form:"mfaToken"`
	OTP               string `json:"otp,omitempty" form:"otp"`
	Mobile            string `json:"mobile,omitempty"`
	Mode              string `json:"mode,omitempty" form:"mode"`
	RequestFrom       string `json:"requestFrom,omitempty" form:"requestFrom"`
	Version           string
	NarwhalMFAEnabled bool
	TransactionType   string `json:"transactionType"`
	AccountID         string `json:"accountId,omitempty"`
	LpseID            string `json:"lpseId"`
	TokenBela         string `json:"tokenBela,omitempty"`
}

type Logout struct {
	Token string `json:"token"`
}

// AccessTokenResponse data structure for json api
type AccessTokenResponse struct {
	ID           string `jsonapi:"primary,authType" json:"authType"`
	UserID       string `jsonapi:"attr,userId,omitempty" json:"userId,omitempty"`
	FirstName    string `jsonapi:"attr,firstName,omitempty" json:"firstName,omitempty"`
	LastName     string `jsonapi:"attr,lastName,omitempty" json:"lastName"`
	FullName     string `json:"fullName,omitempty"`
	Email        string `jsonapi:"attr,email,omitempty" json:"email,omitempty"`
	Token        string `jsonapi:"attr,token" json:"token"`
	NewMember    bool   `jsonapi:"attr,newMember" json:"newMember"`
	HasPassword  bool   `jsonapi:"attr,hasPassword" json:"hasPassword"`
	RefreshToken string `jsonapi:"attr,refreshToken,omitempty" json:"refreshToken,omitempty"`
	ExpiredTime  string `jsonapi:"attr,expiredTime" json:"expiredTime"`
	MemberType   string `jsonapi:"attr,memberType" json:"memberType"`
	Department   string `jsonapi:"attr,department,omitempty" json:"department,omitempty"`
	JobTitle     string `jsonapi:"attr,jobTitle,omitempty" json:"jobTitle,omitempty"`
	Mobile       string `jsonapi:"attr,mobile,omitempty" json:"mobile,omitempty"`
	AccountID    string `json:"accountId,omitempty"`
	CustomToken  string `jsonapi:"attr,customToken" json:"customToken,omitempty"`
}

// RefreshToken data structure
type RefreshToken struct {
	ID              string
	Token           string
	RefreshTokenAge time.Duration
}

// LoginAttempt data structure
type LoginAttempt struct {
	Key             string
	Attempt         string
	LoginAttemptAge time.Duration
}

// LoginSessionRedis Redis Token data struct
type LoginSessionRedis struct {
	Key         string
	Token       string
	ExpiredTime time.Duration
}

// ValidateSocmedRequest data structure for response validate request token socmed
type ValidateSocmedRequest struct {
	Data       *memberModel.Member
	NewMember  bool
	HTTPStatus int
	Error      error
}

// Match function for checking refresh token
func (rt *RefreshToken) Match(refreshToken string) bool {
	return rt.Token == refreshToken
}

// NewRefreshToken function for generating refresh token
func NewRefreshToken(id, token string, age time.Duration) *RefreshToken {
	return &RefreshToken{
		ID:              id,
		Token:           token,
		RefreshTokenAge: age,
	}
}

type CheckEmail struct {
	Email string `json:"email"`
}

// Users data structure
type Users struct {
	UserType    string `json:"userType"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	HasPassword bool   `json:"hasPassword"`
	AccountType string `json:"accountType"`
}

// CheckEmailResponse data structure
type CheckEmailResponse struct {
	Email string  `json:"email"`
	Users []Users `json:"users"`
}

// VerifyResponse data structure
type VerifyResponse struct {
	Adm          bool      `json:"adm"`
	Aud          string    `json:"aud"`
	Authorised   bool      `json:"authorised"`
	Did          string    `json:"did"`
	Dli          string    `json:"dli"`
	Email        string    `json:"email"`
	Exp          float64   `json:"exp"`
	Iat          float64   `json:"iat"`
	Iss          string    `json:"iss"`
	Jti          string    `json:"jti"`
	MemberType   string    `json:"memberType"`
	Staff        bool      `json:"staff"`
	Sub          string    `json:"sub"`
	UserID       string    `json:"userId"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	ExpiredTime  time.Time `json:"expiredTime"`
	Mobile       string    `json:"mobile"`
	IsMerchant   bool      `json:"isMerchant"`
	MerchantID   string    `json:"merchantId"`
	MerchantIDs  []string  `json:"merchantIds"`
	HasPassword  bool      `json:"hasPassword"`
	SignUpFrom   string    `json:"signUpFrom"`
	AccountID    string    `json:"accountId"`
	CustomToken  string    `json:"customToken"`
}

type GoogleCaptcha struct {
	Secret   string `json:"secret" form:"secret"`
	Response string `json:"response" form:"response"`
	RemoteIP string `json:"remoteIp" form:"remoteIp"`
}

// GoogleCaptchaResponse data structure to collect token from google captcha
type GoogleCaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float32   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTs time.Time `json:"challenge_ts"`
	HostName    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// GoogleCaptchaResponseResult data structure to collect token from google captcha
type GoogleCaptchaResponseResult struct {
	Success     bool      `json:"success"`
	Score       float32   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTs time.Time `json:"challengeTs"`
	HostName    string    `json:"hostName"`
	ErrorCodes  []string  `json:"errorCodes"`
}

// InactiveResponse data structure for response inactive status
type InactiveResponse struct {
	Status string `json:"status"`
}

// MFAResponse data structure for response inactive status
type MFAResponse struct {
	MFARequired bool   `jsonapi:"attr,mfaRequired" json:"mfaRequired"`
	MFAToken    string `jsonapi:"attr,mfaToken" json:"mfaToken"`
}

// ClientResponse response for client login
type ClientResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	Email        string    `json:"email"`
	ExpiredAt    time.Time `json:"expiredAt"`
}

// AuthV3TokenResponse specific for v3
type AuthV3TokenResponse struct {
	IsRegistered bool   `json:"isRegistered"`
	GoogleID     string `json:"googleId"`
	FacebookID   string `json:"facebookId"`
	AppleID      string `json:"appleId"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

// CheckEmail specific for v3
type CheckEmailPayload struct {
	Email    string `json:"email" form:"email"`
	UserType string `json:"userType" form:"userType"`
}
