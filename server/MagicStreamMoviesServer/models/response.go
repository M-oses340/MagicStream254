package models

import "time"

// APIResponse defines the standard JSON response format for all endpoints
type APIResponse struct {
	Status    string      `json:"status"`           // "success" or "error"
	Message   string      `json:"message"`          // Human-readable message
	Data      interface{} `json:"data,omitempty"`   // Payload for success
	Errors    interface{} `json:"errors,omitempty"` // Details for validation or server errors
	Timestamp time.Time   `json:"timestamp"`        // Response timestamp
}
