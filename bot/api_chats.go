package bot

import (
	"context"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

// CreateChat creates a new chat.
func (b *Bot) CreateChat(ctx context.Context, botID uuid.UUID, name string, chatType models.ChatType, members []uuid.UUID, sharedHistory optional.Optional[bool], description string) (uuid.UUID, error) {
	return b.Client.CreateChat(ctx, botID, name, chatType, members, sharedHistory, description)
}

// ChatInfo returns chat details.
func (b *Bot) ChatInfo(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) (models.ChatInfo, error) {
	return b.Client.ChatInfo(ctx, botID, chatID)
}

// ListChats returns list of chats.
func (b *Bot) ListChats(ctx context.Context, botID uuid.UUID) ([]models.ChatListItem, error) {
	return b.Client.ListChats(ctx, botID)
}

// AddUser adds users to chat.
func (b *Bot) AddUser(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, users []uuid.UUID) error {
	return b.Client.AddUser(ctx, botID, chatID, users)
}

// RemoveUser removes users from chat.
func (b *Bot) RemoveUser(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, users []uuid.UUID) error {
	return b.Client.RemoveUser(ctx, botID, chatID, users)
}

// AddAdmin adds admins to chat.
func (b *Bot) AddAdmin(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, users []uuid.UUID) error {
	return b.Client.AddAdmin(ctx, botID, chatID, users)
}

// CreateThread creates thread for message.
func (b *Bot) CreateThread(ctx context.Context, botID uuid.UUID, syncID uuid.UUID) (uuid.UUID, error) {
	return b.Client.CreateThread(ctx, botID, syncID)
}

// SetStealth enables stealth mode.
func (b *Bot) SetStealth(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, disableWeb optional.Optional[bool], burnIn optional.Optional[int], expireIn optional.Optional[int]) error {
	return b.Client.SetStealth(ctx, botID, chatID, disableWeb, burnIn, expireIn)
}

// DisableStealth disables stealth mode.
func (b *Bot) DisableStealth(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) error {
	return b.Client.DisableStealth(ctx, botID, chatID)
}

// PinMessage pins a message in chat.
func (b *Bot) PinMessage(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, syncID uuid.UUID) error {
	return b.Client.PinMessage(ctx, botID, chatID, syncID)
}

// UnpinMessage unpins a message.
func (b *Bot) UnpinMessage(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) error {
	return b.Client.UnpinMessage(ctx, botID, chatID)
}
