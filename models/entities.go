package models

import "github.com/google/uuid"

// Forward represents a forwarded message.
type Forward struct {
	ChatID     uuid.UUID
	AuthorID   uuid.UUID
	SyncID     uuid.UUID
	ChatName   string
	Type       string
	InsertedAt string
}

// Reply represents a reply entity.
type Reply struct {
	AuthorID uuid.UUID
	SyncID   uuid.UUID
	Body     string
}
