package bot

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

// SmartAppsList returns list of smartapps.
func (b *Bot) SmartAppsList(ctx context.Context, botID uuid.UUID, version optional.Optional[int]) ([]models.SmartApp, int, error) {
	return b.Client.SmartAppsList(ctx, botID, version)
}

// SendSmartAppEvent sends event to smartapp.
func (b *Bot) SendSmartAppEvent(ctx context.Context, botID uuid.UUID, ref optional.Optional[uuid.UUID], smartappID uuid.UUID, chatID uuid.UUID, data map[string]any, opts optional.Optional[map[string]any], files optional.Optional[[]models.AsyncFile], encrypted bool) error {
	return b.Client.SendSmartAppEvent(ctx, botID, ref, smartappID, chatID, data, opts, files, encrypted)
}

// SendSmartAppNotification sends smartapp notification.
func (b *Bot) SendSmartAppNotification(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, counter int, body optional.Optional[string], opts optional.Optional[map[string]any], meta optional.Optional[map[string]any]) error {
	return b.Client.SendSmartAppNotification(ctx, botID, chatID, counter, body, opts, meta)
}

// SendSmartAppCustomNotification sends custom notification and returns sync_id.
func (b *Bot) SendSmartAppCustomNotification(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, title string, body string, meta optional.Optional[map[string]any], waitCallback bool, timeout time.Duration) (uuid.UUID, error) {
	syncID, err := b.Client.SendSmartAppCustomNotification(ctx, botID, chatID, title, body, meta)
	if err != nil {
		return uuid.Nil, err
	}
	if !waitCallback {
		return syncID, nil
	}
	_, err = b.WaitCallback(syncID, timeout)
	return syncID, err
}

// SendSmartAppUnreadCounter sets unread counter and returns sync_id.
func (b *Bot) SendSmartAppUnreadCounter(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, counter int, waitCallback bool, timeout time.Duration) (uuid.UUID, error) {
	syncID, err := b.Client.SendSmartAppUnreadCounter(ctx, botID, chatID, counter)
	if err != nil {
		return uuid.Nil, err
	}
	if !waitCallback {
		return syncID, nil
	}
	_, err = b.WaitCallback(syncID, timeout)
	return syncID, err
}

// SendSmartAppManifest sends manifest payload and returns current manifest.
func (b *Bot) SendSmartAppManifest(ctx context.Context, botID uuid.UUID, manifest models.SmartAppManifest) (models.SmartAppManifest, error) {
	return b.Client.SendSmartAppManifest(ctx, botID, manifest)
}

// SmartAppUploadFile uploads file for smartapp and returns link.
func (b *Bot) SmartAppUploadFile(ctx context.Context, botID uuid.UUID, filename string, r io.Reader) (string, error) {
	return b.Client.SmartAppUploadFile(ctx, botID, filename, r)
}
