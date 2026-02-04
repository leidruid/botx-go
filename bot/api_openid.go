package bot

import (
	"context"

	"github.com/google/uuid"
)

// RefreshAccessToken refreshes OpenID access token.
func (b *Bot) RefreshAccessToken(ctx context.Context, botID uuid.UUID, userHUID uuid.UUID, ref *uuid.UUID) (bool, error) {
	return b.Client.RefreshAccessToken(ctx, botID, userHUID, ref)
}
