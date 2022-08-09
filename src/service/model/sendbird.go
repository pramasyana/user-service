package model

import "gopkg.in/guregu/null.v4/zero"

// SendbirdRequest struct
type SendbirdRequest struct {
	User
	MemberType string   `json:"member_type"`
	ExpiresAt  int64    `json:"expires_at"`
	Token      string   `json:"token"`
	Metadata   MetaData `json:"metadata"`
}

// SendbirdRequest struct
type SendbirdRequestV4 struct {
	User
	MemberType string     `json:"member_type"`
	ExpiresAt  int64      `json:"expires_at"`
	Token      string     `json:"token"`
	MetadataV4 MetaDataV4 `json:"metadata"`
	Client     string     `json:"client"`
}

// User struct
type User struct {
	UserID     string `json:"user_id"`
	NickName   string `json:"nickname"`
	ProfileURL string `json:"profile_url"`
}

// Merchant struct
type Merchant struct {
	IsMerchant   string `json:"is_merchant"`
	MerchantID   string `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	Reference	 string `json:"reference"`
}

// SessionTokenResponse struct
type SessionTokenResponse struct {
	SendbirdErrorResponse
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type SessionTokenRequest struct {
	SendbirdErrorResponse
	ExpiresAt int64 `json:"expires_at"`
}

// MetaData struct
type MetaData struct {
	Token        SessionTokenResponse `json:"token"`
	Merchant     Merchant             `json:"merchant"`
	MerchantLogo zero.String          `json:"merchant_logo"`
}

// MetaDataV4 struct
type MetaDataV4 struct {
	Token     SessionTokenResponse `json:"token"`
	Reference string               `json:"reference"`
}

// MetadataRequest struct
type MetadataRequest struct {
	Metadata MetaDataResponse `json:"metadata"`
}

// MetadataRequestV1 struct
type MetadataRequestV1 struct {
	Metadata MetaDataResponseV1 `json:"metadata"`
}

// MetadataRequestV4 struct
type MetadataRequestV4 struct {
	MetadataV4 MetaDataResponseV4 `json:"metadata"`
}

// MetaDataResponse struct
type MetaDataResponse struct {
	SendbirdErrorResponse
	Token        string      `json:"token"`
	Merchant     string      `json:"merchant,omitempty"`
	MerchantLogo zero.String `json:"merchant_logo,omitempty"`
	Reference    string		 `json:"reference"`
}

// MetaDataResponse struct
type MetaDataResponseV1 struct {
	SendbirdErrorResponse
	Token string `json:"token"`
}

// MetaDataResponse struct
type MetaDataResponseV4 struct {
	SendbirdErrorResponse
	Token string `json:"token"`
}

// SendbirdResponse struct
type SendbirdResponse struct {
	SendbirdErrorResponse
	User
	Metadata MetaData `json:"metadata"`
}

// SendbirdResponseV4 struct
type SendbirdResponseV4 struct {
	SendbirdErrorResponse
	User
	MetadataV4 MetaDataV4 `json:"metadata"`
}

// SendbirdStringResponse struct
type SendbirdStringResponse struct {
	SendbirdErrorResponse
	User
	Metadata MetaDataResponse `json:"metadata"`
}

// SendbirdStringResponseV4 struct
type SendbirdStringResponseV4 struct {
	SendbirdErrorResponse
	User
	MetadataV4 MetaDataResponseV4 `json:"metadata"`
}

// SendbirdErrorResponse struct
type SendbirdErrorResponse struct {
	Error   bool   `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
