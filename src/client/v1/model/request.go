package model

import (
	corporateModel "github.com/Bhinneka/user-service/src/corporate/v2/model"
)

// LKPPUser payload
type LKPPUser struct {
	Payload   Payload `json:"payload"`
	TokenBela string  `json:"token"`
}

// ContactKafka payload
type ContactKafka struct {
	EventType string                        `json:"eventType"`
	Payload   corporateModel.ContactPayload `json:"payload"`
}

// Payload generic payload
type Payload struct {
	UserName  string `json:"userName"`
	RealName  string `json:"realName"`
	Role      string `json:"role"`
	LpseID    string `json:"lpseId"`
	IsLatihan bool   `json:"isLatihan"`
	Time      string `json:"time"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
}

const (
	// DeviceIDBela deviceID assign to Bela LKPP Request
	DeviceIDBela = "belaDID"
	// BelaLKPP signUpFrom flag in DB
	BelaLKPP = "lkpp"
)

// QueryParam generic param, can be used as JSON payload or query string in URL
type QueryParam struct {
	Token       string `json:"token" query:"token"`
	RedirectURL string `json:"returnUrl" query:"returnUrl"`
}

// TokenResponse return verify token response
type TokenResponse struct {
	Email      string  `json:"email"`
	Exp        float64 `json:"exp"`
	IssueAt    float64 `json:"iat"`
	Issuer     string  `json:"iss"`
	Token      string  `json:"token"`
	SignUpFrom string  `json:"signUpFrom"`
}
