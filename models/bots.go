package models

import "github.com/google/uuid"

// BotsListItem represents bot entry in catalog.
type BotsListItem struct {
	ID          uuid.UUID
	Name        string
	Description string
	Avatar      string
	Enabled     bool
}
