package model

import (
	"time"
)

// TokenActivation data structure
type TokenActivation struct {
	ID    string
	Value string
	TTL   time.Duration
}
