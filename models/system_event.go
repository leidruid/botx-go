package models

import "github.com/google/uuid"

// SystemEvent represents a system event command.
type SystemEvent struct {
	Type   string
	Bot    BotContext
	Chat   ChatContext
	User   UserContext
	SyncID uuid.UUID
	Data   map[string]any
	Raw    map[string]any
}
