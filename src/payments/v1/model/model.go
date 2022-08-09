package model

import "time"

const (
	ErrorTokenInvalid = "Token data is invalid"
	MessageGetData    = "Get Data Payment Token"
	SaveSuccess       = "Save token success"
)

// MemberError data structure
type PaymentsError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	ID        string `jsonapi:"primary" json:"id"`
	Message   string `jsonapi:"attr,message" json:"message,omitempty"`
	Email     string `jsonapi:"attr,email" json:"email,omitempty"`
	Token     string `jsonapi:"attr,token" json:"token,omitempty"`
	Channel   string `jsonapi:"attr,channel" json:"channel,omitempty"`
	Method    string `jsonapi:"attr,method" json:"method,omitempty"`
	ExpiredAt string `jsonapi:"attr,expiredAt" json:"expiredAt,omitempty"`
}

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
}

type Payments struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Channel   string `json:"channel"`
	Method    string `json:"method"`
	Token     string `json:"token"`
	ExpiredAt time.Time
}

type PaymentsInput struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Channel   string `json:"channel"`
	Method    string `json:"method"`
	Token     string `json:"token"`
	ExpiredAt string `json:"expiredAt"`
}

type QueryPaymentParameters struct {
	Page     int
	Limit    int
	Offset   int
	StrPage  string `json:"page" query:"page" form:"page" param:"page"`
	StrLimit string `json:"limit" query:"limit" form:"limit" param:"limit"`
	OrderBy  string `json:"orderBy" query:"orderBy" form:"orderBy" param:"orderBy"`
	SortBy   string `json:"sortBy" query:"sortBy" form:"sortBy" param:"sortBy"`
	Search   string `json:"search" query:"search" form:"search" param:"search"`
	Channel  string `json:"channel" query:"channel" form:"channel" param:"channel"`
	Method   string `json:"method" query:"method" form:"method" param:"method"`
	Email    string `json:"email" query:"email" form:"email" param:"email"`
}
