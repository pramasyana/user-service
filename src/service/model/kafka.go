package model

type QueuePayload struct {
	GeneralPayload
	Auth string `json:"token,omitempty"`
}

type Messages struct {
	Key     string
	Content []byte
}

type GeneralPayload struct {
	EventType string      `json:"eventType"`
	Payload   interface{} `json:"payload"`
}
