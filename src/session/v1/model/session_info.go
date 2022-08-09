package model

import (
	"time"
)

const (
	// MessageSuccess constanta
	MessageSuccess = "Success Get All Session Info"
)

// SessionInfoRequest data structure
type SessionInfoRequest struct {
	GrantType   string `json:"grantType"`
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	Email       string `json:"email"`
	IP          string `json:"ip"`
	UserAgent   string `json:"userAgent"`
	DeviceID    string `json:"deviceId"`
	DeviceLogin string `json:"deviceLogin"`
	JTI         string `json:"jti"`
}

// SessionInfoResponse data structure
type SessionInfoResponse struct {
	ID         *string    `json:"id"`
	GrantType  *string    `json:"grantType"`
	ClientType *string    `json:"ClientType"`
	UserID     *string    `json:"userId"`
	UserName   *string    `json:"userName"`
	IP         *string    `json:"ip"`
	JTI        *string    `json:"jti"`
	UserAgent  *string    `json:"userAgent"`
	DeviceID   *string    `json:"deviceId"`
	CreatedAt  *time.Time `json:"createdAt"`
}

// SessionInfoList data structure
type SessionInfoList struct {
	Data      []SessionInfoResponse `json:"data"`
	TotalData int
}

// ParamList param structure
type ParamList struct {
	Query      string `json:"query"`
	StrPage    string `json:"page"`
	Page       int    `default:"1"`
	StrLimit   string `json:"limit"`
	Limit      int    `default:"10"`
	Sort       string `json:"sort"`
	OrderBy    string `json:"orderBy"`
	Offset     int    `default:"0"`
	Email      string `json:"email"`
	MemberID   string `json:"userId"`
	ClientType string `json:"clientType"`
	Range      string `json:"range"`
	RangeInt   int
}

// AllowedSortFields is allowed field to sorting data
var AllowedSortFields = []string{
	"id",
	"userName",
	"userId",
	"clientType",
	"grantType",
	"createdAt",
}

// ParametersGetSession param structure
type ParametersGetSession struct {
	DeviceID   string `json:"deviceId"`
	ClientType string `json:"clientType"`
	UserID     string `json:"userId"`
	Jti        string `json:"jti"`
	SessionID  string `json:"id"`
}
