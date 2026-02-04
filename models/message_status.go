package models

import (
	"time"

	"github.com/google/uuid"
)

// MessageStatus describes delivery/read status for message.
type MessageStatus struct {
	GroupChatID uuid.UUID
	SentTo      []uuid.UUID
	ReadBy      map[uuid.UUID]time.Time
	ReceivedBy  map[uuid.UUID]time.Time
}
