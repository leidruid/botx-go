package bot

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/leidruid/botx-go/models"
)

var (
	ErrInvalidCommandName    = errors.New("command must start with '/' and must not include spaces")
	ErrCommandExists         = errors.New("command handler already registered")
	ErrNoSyncSmartAppHandler = errors.New("sync smartapp handler not found")
)

var commandNameRe = regexp.MustCompile(`^/[^\s/]+$`)

type VisibleFunc func(ctx context.Context, recipient models.StatusRecipient, b *Bot) (bool, error)

type commandHandler struct {
	handler     MessageHandler
	description string
	visible     bool
	visibleFunc VisibleFunc
	middlewares []Middleware
}

type HandlerCollector struct {
	commands       map[string]commandHandler
	defaultHandler *commandHandler
	systemHandlers map[string]SystemEventHandler
	syncSmartApp   SyncSmartAppHandler
	middlewares    []Middleware
}

func NewHandlerCollector(middlewares ...Middleware) *HandlerCollector {
	return &HandlerCollector{
		commands:       make(map[string]commandHandler),
		systemHandlers: make(map[string]SystemEventHandler),
		middlewares:    middlewares,
	}
}

// Include merges handlers from other collectors.
func (c *HandlerCollector) Include(others ...*HandlerCollector) {
	for _, o := range others {
		for k, v := range o.commands {
			c.commands[k] = v
		}
		if o.defaultHandler != nil {
			c.defaultHandler = o.defaultHandler
		}
		for k, v := range o.systemHandlers {
			c.systemHandlers[k] = v
		}
		if o.syncSmartApp != nil {
			c.syncSmartApp = o.syncSmartApp
		}
	}
}

// OnCommand registers a command handler.
func (c *HandlerCollector) OnCommand(name string, handler MessageHandler, opts ...CommandOption) error {
	if !commandNameRe.MatchString(name) {
		return ErrInvalidCommandName
	}
	if _, ok := c.commands[name]; ok {
		return ErrCommandExists
	}

	cfg := commandOptions{visible: true}
	for _, opt := range opts {
		opt(&cfg)
	}

	c.commands[name] = commandHandler{
		handler:     handler,
		description: cfg.description,
		visible:     cfg.visible,
		visibleFunc: cfg.visibleFunc,
		middlewares: append([]Middleware{}, cfg.middlewares...),
	}
	return nil
}

// Default sets default handler.
func (c *HandlerCollector) Default(handler MessageHandler, middlewares ...Middleware) {
	c.defaultHandler = &commandHandler{
		handler:     handler,
		visible:     false,
		middlewares: middlewares,
	}
}

// OnSystemEvent registers system event handler by type.
func (c *HandlerCollector) OnSystemEvent(eventType string, handler SystemEventHandler) {
	c.systemHandlers[eventType] = handler
}

// OnSyncSmartApp registers sync smartapp handler.
func (c *HandlerCollector) OnSyncSmartApp(handler SyncSmartAppHandler) {
	c.syncSmartApp = handler
}

func (c *HandlerCollector) HandleIncomingMessage(ctx context.Context, msg models.IncomingMessage, b *Bot) error {
	cmd := firstToken(msg.Body)
	h, ok := c.commands[cmd]
	if !ok {
		if c.defaultHandler == nil {
			return nil
		}
		h = *c.defaultHandler
	}
	return applyMiddlewares(ctx, msg, b, h.handler, append(c.middlewares, h.middlewares...)...)
}

func (c *HandlerCollector) HandleSystemEvent(ctx context.Context, ev models.SystemEvent, b *Bot) error {
	h, ok := c.systemHandlers[ev.Type]
	if !ok {
		return nil
	}
	return h(ctx, ev, b)
}

func (c *HandlerCollector) HandleSyncSmartApp(ctx context.Context, ev models.SmartAppEvent, b *Bot) (models.SyncSmartAppResponse, error) {
	if c.syncSmartApp == nil {
		return models.SyncSmartAppResponse{}, ErrNoSyncSmartAppHandler
	}
	return c.syncSmartApp(ctx, ev, b)
}

func (c *HandlerCollector) BotMenu(ctx context.Context, recipient models.StatusRecipient, b *Bot) (models.BotMenu, error) {
	menu := models.BotMenu{}
	for name, h := range c.commands {
		if !h.visible {
			continue
		}
		if h.visibleFunc != nil {
			ok, err := h.visibleFunc(ctx, recipient, b)
			if err != nil || !ok {
				continue
			}
		}
		if h.description == "" {
			continue
		}
		menu[name] = h.description
	}
	return menu, nil
}

func applyMiddlewares(ctx context.Context, msg models.IncomingMessage, b *Bot, handler MessageHandler, mws ...Middleware) error {
	wrapped := handler
	for i := len(mws) - 1; i >= 0; i-- {
		mw := mws[i]
		next := wrapped
		wrapped = func(ctx context.Context, msg models.IncomingMessage, b *Bot) error {
			return mw(ctx, msg, b, next)
		}
	}
	return wrapped(ctx, msg, b)
}

func firstToken(s string) string {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

type commandOptions struct {
	description string
	visible     bool
	visibleFunc VisibleFunc
	middlewares []Middleware
}

type CommandOption func(*commandOptions)

func WithDescription(desc string) CommandOption {
	return func(o *commandOptions) { o.description = desc }
}

func WithVisible(visible bool) CommandOption {
	return func(o *commandOptions) { o.visible = visible }
}

func WithVisibleFunc(fn VisibleFunc) CommandOption {
	return func(o *commandOptions) { o.visibleFunc = fn }
}

func WithMiddlewares(mws ...Middleware) CommandOption {
	return func(o *commandOptions) { o.middlewares = append(o.middlewares, mws...) }
}
