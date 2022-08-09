package model

import (
	"strings"
	"time"
)

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

// Status initialize status data type
type Status int

const (
	// InActive status user
	InActive Status = iota
	// Active status user
	Active
	// Blocked status user
	Blocked

	active   = "ACTIVE"
	inactive = "INACTIVE"
	blocked  = "BLOCKED"
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
	}
	return Active
}

// Authenticate function for authenticating
func (c *ClientApp) Authenticate(secret string) bool {
	return c.ClientSecret == secret
}

// IsActive check if client status is active
func (c *ClientApp) IsActive() bool {
	return c.Status == Active
}
