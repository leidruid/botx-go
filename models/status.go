package models

import "github.com/google/uuid"

// StatusRecipient describes user requesting status.
type StatusRecipient struct {
	BotID uuid.UUID
	HUID  uuid.UUID
}

// BotMenu is a map of command -> description.
type BotMenu map[string]string
