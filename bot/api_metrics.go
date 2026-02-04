package bot

import (
	"context"

	"github.com/google/uuid"
)

// CollectBotFunction reports bot function usage.
func (b *Bot) CollectBotFunction(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, huids []uuid.UUID, botFunction string) error {
	return b.Client.CollectBotFunction(ctx, botID, chatID, huids, botFunction)
}
