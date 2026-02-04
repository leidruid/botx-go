package middleware

import (
	"context"
	"reflect"

	"github.com/leidruid/botx-go/bot"
	"github.com/leidruid/botx-go/models"
)

// ExceptionHandler handles errors from message handlers.
type ExceptionHandler func(ctx context.Context, msg models.IncomingMessage, b *bot.Bot, err error) error

// ExceptionHandlers maps error types to handlers.
type ExceptionHandlers map[reflect.Type]ExceptionHandler

// NewExceptionMiddleware creates a middleware that dispatches errors by type.
func NewExceptionMiddleware(handlers ExceptionHandlers) bot.Middleware {
	return func(ctx context.Context, msg models.IncomingMessage, b *bot.Bot, next bot.MessageHandler) error {
		err := next(ctx, msg, b)
		if err == nil {
			return nil
		}

		for t := reflect.TypeOf(err); t != nil; t = t.Elem() {
			if h, ok := handlers[t]; ok {
				return h(ctx, msg, b, err)
			}
			if t.Kind() != reflect.Ptr {
				break
			}
		}
		// Fallback: match by assignable types
		for et, h := range handlers {
			if reflect.TypeOf(err).AssignableTo(et) {
				return h(ctx, msg, b, err)
			}
		}
		return err
	}
}
