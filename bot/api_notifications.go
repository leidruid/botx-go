package bot

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/client"
)

// SendMessage sends a message and optionally waits for callback.
func (b *Bot) SendMessage(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, body string, opts client.SendMessageOptions, waitCallback bool, callbackTimeout time.Duration) (uuid.UUID, error) {
	syncID, err := b.Client.SendMessage(ctx, botID, chatID, body, opts)
	if err != nil {
		return uuid.Nil, err
	}
	if !waitCallback {
		return syncID, nil
	}
	_, err = b.WaitCallback(syncID, callbackTimeout)
	if err != nil {
		return syncID, err
	}
	return syncID, nil
}
