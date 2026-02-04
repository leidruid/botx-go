package bot

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/client"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

// ReplyEvent replies to a message by sync_id.
func (b *Bot) ReplyEvent(ctx context.Context, botID uuid.UUID, sourceSyncID uuid.UUID, body string, opts client.SendMessageOptions) error {
	return b.Client.ReplyEvent(ctx, botID, sourceSyncID, body, opts)
}

// EditEvent edits a message by sync_id.
func (b *Bot) EditEvent(ctx context.Context, botID uuid.UUID, syncID uuid.UUID, body optional.Optional[string], metadata optional.Optional[map[string]any], bubble optional.Optional[*models.Markup], keyboard optional.Optional[*models.Markup], attachment optional.Optional[*models.OutgoingAttachment], markupAutoAdjust optional.Optional[bool]) error {
	return b.Client.EditEvent(ctx, botID, syncID, body, metadata, bubble, keyboard, attachment, markupAutoAdjust)
}

// DeleteEvent deletes a message by sync_id.
func (b *Bot) DeleteEvent(ctx context.Context, botID uuid.UUID, syncID uuid.UUID) error {
	return b.Client.DeleteEvent(ctx, botID, syncID)
}

// Typing sends typing event.
func (b *Bot) Typing(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) error {
	return b.Client.Typing(ctx, botID, chatID)
}

// StopTyping stops typing event.
func (b *Bot) StopTyping(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) error {
	return b.Client.StopTyping(ctx, botID, chatID)
}

// MessageStatus gets delivery/read status for message.
func (b *Bot) MessageStatus(ctx context.Context, botID uuid.UUID, syncID uuid.UUID) (models.MessageStatus, error) {
	return b.Client.MessageStatus(ctx, botID, syncID)
}

// WaitCallbackWithTimeout waits for callback with explicit timeout.
func (b *Bot) WaitCallbackWithTimeout(syncID uuid.UUID, timeout time.Duration) (models.MethodCallback, error) {
	return b.WaitCallback(syncID, timeout)
}
