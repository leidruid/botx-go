package models

import (
	"time"

	"github.com/google/uuid"
)

// ChatType represents chat type.
type ChatType string

const (
	ChatPersonal ChatType = "chat"
	ChatGroup    ChatType = "group_chat"
	ChatChannel  ChatType = "channel"
	ChatThread   ChatType = "thread"
)

// ChatListItem represents chat summary.
type ChatListItem struct {
	ChatID        uuid.UUID
	ChatType      ChatType
	Name          string
	Description   string
	Members       []uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	SharedHistory bool
}

// ChatInfoMember represents chat member.
type ChatInfoMember struct {
	HUID    uuid.UUID
	IsAdmin bool
	Kind    string
}

// ChatInfo represents detailed chat info.
type ChatInfo struct {
	ChatID        uuid.UUID
	ChatType      ChatType
	Name          string
	Description   string
	CreatorID     *uuid.UUID
	CreatedAt     time.Time
	Members       []ChatInfoMember
	SharedHistory bool
}
