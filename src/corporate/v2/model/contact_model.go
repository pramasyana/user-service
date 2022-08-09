package model

import (
	"time"
)

const (
	// StatusActivated for user active
	StatusActivated = "ACTIVATED"
)

// ParametersContact data structure
type ParametersContact struct {
	StrPage  string `json:"strPage" form:"strPage" query:"strPage" validate:"omitempty,numeric" fieldname:"strPage" url:"strPage"`
	Page     int    `json:"page" form:"page" query:"page" validate:"omitempty,numeric" fieldname:"page" url:"page"`
	StrLimit string `json:"strLimit" form:"strLimit" query:"strLimit" validate:"omitempty" fieldname:"strLimit" url:"strLimit"`
	Limit    int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,numeric" fieldname:"limit" url:"limit"`
	Offset   int    `json:"offset" form:"offset" query:"offset" validate:"omitempty,numeric" fieldname:"offset" url:"offset"`
	Query    string `json:"query" form:"query" query:"query" validate:"omitempty,string" fieldname:"query" url:"query"`
	Status   string `json:"status" form:"status" query:"status" validate:"omitempty,string" fieldname:"status" url:"status"`
	IsNew    string `json:"isNew" form:"isNew" query:"isNew" validate:"omitempty,string" fieldname:"isNew" url:"isNew"`
}

type ContactPayload struct {
	ID              int       `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	CreatedAt       time.Time `json:"createdAt"`
	AccountID       string    `json:"accountId"`
	PhoneNumber     string    `json:"phoneNumber"`
	TransactionType string    `json:"transactionType"`
	LpseID          string    `json:"lpseId"`
	Password        string    `json:"password"`
	Salt            string    `json:"salt"`
}
