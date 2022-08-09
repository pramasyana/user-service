package model

import "github.com/Bhinneka/user-service/src/shared"

// Payload for log
type Payload struct {
	ID         string   `json:"id,omitempty"`
	Module     string   `json:"module"`
	Action     string   `json:"action"`
	Logs       []Log    `json:"logs"`
	Pack       string   `json:"pack"`
	Target     string   `json:"target"`
	CreatorID  string   `json:"creatorId"`
	CreatorIP  string   `json:"creatorIP"`
	EditorID   string   `json:"editorId"`
	ViewType   []string `json:"viewType"`
	ObjectType string   `json:"objectType"`
	User       Users    `json:"user"`
}

// Log basic
type Log struct {
	Field    string `json:"field"`
	OldValue string `json:"oldValue"`
	NewValue string `json:"newValue"`
}
type Users struct {
	Id       string `json:"id"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

// ServiceResult data structure
type ServiceResult struct {
	Result interface{}
	Error  error
	Meta   shared.Meta
}

// ResponseService data struct
type ResponseService struct {
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success,omitempty"`
	Meta    shared.Meta `json:"meta,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}
