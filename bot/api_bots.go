package bot

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

// BotsList returns bots catalog.
func (b *Bot) BotsList(ctx context.Context, botID uuid.UUID, since optional.Optional[time.Time]) ([]models.BotsListItem, time.Time, error) {
	return b.Client.BotsList(ctx, botID, since)
}
