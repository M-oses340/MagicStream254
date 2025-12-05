package models

import "time"

// APIResponse defines a standard JSON response
type APIResponse struct {
	Status    string      `json:"status"`            // "success" or "error"
	Error     bool        `json:"error"`             // true for errors, false for success
	Message   string      `json:"message"`           // Human-readable message
	Content   interface{} `json:"content,omitempty"` // Payload (UserContent, etc.)
	Timestamp time.Time   `json:"timestamp"`         // Response timestamp
}
