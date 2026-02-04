package bot

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/accounts"
	"github.com/leidruid/botx-go/auth"
	"github.com/leidruid/botx-go/callbacks"
	"github.com/leidruid/botx-go/client"
	"github.com/leidruid/botx-go/models"
)

const defaultCallbackTimeout = 30 * time.Second

// Bot is the main entrypoint for BotX runtime and API calls.
type Bot struct {
	Client                 *client.Client
	Accounts               *accounts.Storage
	Collector              *HandlerCollector
	Callbacks              *callbacks.Manager
	DefaultCallbackTimeout time.Duration
}

func NewBot(botAccounts []models.BotAccountWithSecret, collectors ...*HandlerCollector) *Bot {
	storage := accounts.NewStorage(botAccounts)
	c := client.NewClient(storage)
	mainCollector := NewHandlerCollector()
	mainCollector.Include(collectors...)
	return &Bot{
		Client:                 c,
		Accounts:               storage,
		Collector:              mainCollector,
		Callbacks:              callbacks.NewManager(callbacks.NewMemoryRepo()),
		DefaultCallbackTimeout: defaultCallbackTimeout,
	}
}

// VerifyRequest verifies incoming BotX request using Authorization header.
func (b *Bot) VerifyRequest(botID uuid.UUID, authHeader string, trustedIssuers map[string]struct{}) error {
	acc, err := b.Accounts.GetAccount(botID)
	if err != nil {
		return err
	}
	host, err := acc.Host()
	if err != nil {
		return err
	}
	return auth.VerifyRequest(authHeader, acc.SecretKey, host, botID.String(), trustedIssuers)
}

// HandleIncomingMessage routes message to command/default handler.
func (b *Bot) HandleIncomingMessage(ctx context.Context, msg models.IncomingMessage) error {
	if _, err := b.Accounts.GetAccount(msg.Bot.ID); err != nil {
		return err
	}
	return b.Collector.HandleIncomingMessage(ctx, msg, b)
}

// HandleSystemEvent routes system event handler.
func (b *Bot) HandleSystemEvent(ctx context.Context, ev models.SystemEvent) error {
	if _, err := b.Accounts.GetAccount(ev.Bot.ID); err != nil {
		return err
	}
	return b.Collector.HandleSystemEvent(ctx, ev, b)
}

// HandleSyncSmartApp routes sync smartapp event handler.
func (b *Bot) HandleSyncSmartApp(ctx context.Context, ev models.SmartAppEvent) (models.SyncSmartAppResponse, error) {
	if _, err := b.Accounts.GetAccount(ev.Bot.ID); err != nil {
		return models.SyncSmartAppResponse{}, err
	}
	return b.Collector.HandleSyncSmartApp(ctx, ev, b)
}

// Status returns bot menu for recipient.
func (b *Bot) Status(ctx context.Context, recipient models.StatusRecipient) (models.BotMenu, error) {
	if _, err := b.Accounts.GetAccount(recipient.BotID); err != nil {
		return nil, err
	}
	return b.Collector.BotMenu(ctx, recipient, b)
}

// SetCallback stores async method callback.
func (b *Bot) SetCallback(callback models.MethodCallback) error {
	return b.Callbacks.Set(callback)
}

// WaitCallback waits for callback with timeout.
func (b *Bot) WaitCallback(syncID uuid.UUID, timeout time.Duration) (models.MethodCallback, error) {
	if timeout == 0 {
		timeout = b.DefaultCallbackTimeout
	}
	return b.Callbacks.Wait(syncID, timeout)
}

// EnsureNonNilCollector ensures bot has a collector.
func (b *Bot) EnsureNonNilCollector() error {
	if b.Collector == nil {
		return errors.New("handler collector is nil")
	}
	return nil
}
