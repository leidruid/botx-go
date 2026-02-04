package models

import "github.com/google/uuid"

// BotContext identifies bot.
type BotContext struct {
	ID uuid.UUID `json:"bot_id"`
}

// ChatContext identifies chat.
type ChatContext struct {
	ID uuid.UUID `json:"chat_id"`
}

// UserContext identifies user.
type UserContext struct {
	HUID uuid.UUID `json:"huid"`
}
