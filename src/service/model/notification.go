package model

import "time"

const (
	// NoReply account email
	NoReply = "noreply@bhinneka.com"
	// NoReplyName account email name
	NoReplyName = "Bhinneka.com"
	// Source client name
	Source = "user-services"
	// RedisKey key of redis auth
	RedisKey     = "email-config-auth"
	DeviceIDAuth = "AuthEmailSTG"

	EmailCare = "care@bhinneka.com"
)

// Attachment data structure
type Attachment struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// Email data structure
type Email struct {
	From        string       `json:"from"`
	FromName    string       `json:"fromName"`
	To          []string     `json:"to"`
	ToName      []string     `json:"toName"`
	CC          []string     `json:"cc"`
	CCName      []string     `json:"ccName"`
	BCC         []string     `json:"bcc"`
	BCCName     []string     `json:"bccName"`
	Subject     string       `json:"subject"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachments"`
}

// PayloadEmail data structure for email payload request
type PayloadEmail struct {
	Data struct {
		Attributes Email `json:"attributes"`
	} `json:"data"`
}

// ErrorMessage data structure for error sending email
type ErrorMessage struct {
	Errors []struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

// SuccessMessage data structure for success sending email
type SuccessMessage struct {
	Data APIResponse `json:"data"`
}

type APIResponse struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
}

type Template struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  string    `json:"createdBy"`
	ModifiedAt time.Time `json:"modifiedAt"`
	ModifiedBy string    `json:"modifiedBy"`
}

// AuthResponse data structure for email services authentication
type AuthResponse struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			ExpiredAt    string `json:"expiredTime"`
			ExpiredUnix  int    `json:"expiredUnix"`
			RefreshToken string `json:"refreshToken"`
			Token        string `json:"token"`
		} `json:"attributes"`
	} `json:"data"`
}

// ResponseZipCode area data structure
type ResponseGetTemplate struct {
	Succcess bool     `json:"success"`
	Code     int      `json:"code"`
	Message  string   `json:"message"`
	Data     Template `json:"data"`
}
