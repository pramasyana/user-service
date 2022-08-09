package model

import (
	"strings"
	"time"

	"github.com/Bhinneka/golib"
)

// Status initialize status data type
type Status int

const (
	// InActive status user
	InActive Status = iota
	// Active status user
	Active
	// Blocked status user
	Blocked
	//New status user
	New

	active   = "ACTIVE"
	inactive = "INACTIVE"
	new      = "NEW"
	blocked  = "BLOCKED"
	revoked  = "REVOKED"
	invited  = "INVITED"
)

// String function for converting user status
func (s Status) String() string {
	switch s {
	case InActive:
		return inactive
	case Active:
		return active
	case Blocked:
		return blocked
	case New:
		return new

	}
	return active
}

// StringToStatus function for converting string user status to int
func StringToStatus(s string) Status {
	switch strings.ToUpper(s) {
	case active:
		return Active
	case inactive:
		return InActive
	case blocked:
		return Blocked
	case new:
		return New
	}
	return Active
}

// ClientApp data structure
type ClientApp struct {
	ID           string    `json:"id"`
	ClientID     string    `json:"clientId"`
	ClientSecret string    `json:"clientSecret"`
	Name         string    `json:"clientName"`
	Status       Status    `json:"status"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Version      int       `json:"version"`
	Secret       string    `json:"Secret"`
}

// NewClientApp function
func NewClientApp(name string) *ClientApp {
	clientSecret := golib.RandomString(32)
	return &ClientApp{
		ClientID:     name,
		ClientSecret: clientSecret,
		Name:         name,
		Status:       Active,
		Created:      time.Now(),
		LastModified: time.Now(),
		Version:      0,
	}
}

// Authenticate function for authenticating
func (c *ClientApp) Authenticate(secret string) bool {
	return c.ClientSecret == secret
}

// IsActive check if client status is active
func (c *ClientApp) IsActive() bool {
	return c.Status == Active
}
