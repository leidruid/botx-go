package bot

import (
	"context"

	"github.com/leidruid/botx-go/models"
)

type MessageHandler func(ctx context.Context, msg models.IncomingMessage, b *Bot) error

type SystemEventHandler func(ctx context.Context, ev models.SystemEvent, b *Bot) error

type SyncSmartAppHandler func(ctx context.Context, ev models.SmartAppEvent, b *Bot) (models.SyncSmartAppResponse, error)

type Middleware func(ctx context.Context, msg models.IncomingMessage, b *Bot, next MessageHandler) error

type ExceptionHandler func(ctx context.Context, msg models.IncomingMessage, b *Bot, err error) error
