package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// MethodCallback represents async method callback payload.
type MethodCallback struct {
	SyncID    uuid.UUID       `json:"sync_id"`
	Status    string          `json:"status"`
	Reason    string          `json:"reason,omitempty"`
	Errors    []string        `json:"errors,omitempty"`
	Result    json.RawMessage `json:"result,omitempty"`
	ErrorData json.RawMessage `json:"error_data,omitempty"`
}
