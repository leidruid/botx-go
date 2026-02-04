package models

import "github.com/google/uuid"

// SmartAppEvent represents smartapp event (sync/async).
type SmartAppEvent struct {
	Bot    BotContext
	Chat   ChatContext
	Sender UserContext
	SyncID uuid.UUID
	Data   map[string]any
	Raw    map[string]any
}

// SyncSmartAppResponse represents response for sync smartapp.
type SyncSmartAppResponse struct {
	Status string         `json:"status"`
	Result map[string]any `json:"result,omitempty"`
	Reason string         `json:"reason,omitempty"`
	Data   map[string]any `json:"data,omitempty"`
}
