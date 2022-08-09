package model

import (
	"time"

	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// DolphinLog data structure
type DolphinLog struct {
	ID        int
	UserID    string
	EventType string
	LogData   *serviceModel.DolphinPayloadNSQ
	Created   time.Time
}
