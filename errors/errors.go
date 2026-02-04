package errors

import (
	"fmt"
)

// APIError represents an error response from BotX API.
type APIError struct {
	StatusCode int
	Status     string
	Reason     string
	Message    string
	Data       any
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("botx api error: status=%s reason=%s message=%s", e.Status, e.Reason, e.Message)
	}
	return fmt.Sprintf("botx api error: status=%s reason=%s", e.Status, e.Reason)
}

// CallbackError represents a failed callback for async methods.
type CallbackError struct {
	SyncID string
	Reason string
	Data   any
}

func (e *CallbackError) Error() string {
	return fmt.Sprintf("botx callback error: sync_id=%s reason=%s", e.SyncID, e.Reason)
}
