package model

// Health data structure
type Health struct {
	State        int      `json:"state"`
	Dependencies []string `json:"Dependencies"`
	ErrorCount   int
}
