package botx

import (
	"github.com/leidruid/botx-go/bot"
	"github.com/leidruid/botx-go/client"
	"github.com/leidruid/botx-go/models"
)

// Re-export common types for convenience.
type Bot = bot.Bot

type HandlerCollector = bot.HandlerCollector

type BotAccountWithSecret = models.BotAccountWithSecret

type IncomingMessage = models.IncomingMessage

type SystemEvent = models.SystemEvent

type SmartAppEvent = models.SmartAppEvent

type StatusRecipient = models.StatusRecipient

type Markup = models.Markup
type Button = models.Button
type ButtonTextAlign = models.ButtonTextAlign

type SendMessageOptions = client.SendMessageOptions

type ChatType = models.ChatType
type ChatListItem = models.ChatListItem
type ChatInfo = models.ChatInfo

type UserFromSearch = models.UserFromSearch
type UserFromCSV = models.UserFromCSV

type AsyncFile = models.AsyncFile
type SmartApp = models.SmartApp
type SmartAppManifest = models.SmartAppManifest
type SmartAppManifestWebParams = models.SmartAppManifestWebParams
type SmartAppManifestMobileParams = models.SmartAppManifestMobileParams
type SmartAppManifestUnreadCounterParams = models.SmartAppManifestUnreadCounterParams

type Sticker = models.Sticker
type StickerPack = models.StickerPack
type StickerPackFromList = models.StickerPackFromList
type StickerPackPage = models.StickerPackPage

var NewBot = bot.NewBot
var NewHandlerCollector = bot.NewHandlerCollector
var WithDescription = bot.WithDescription
var WithVisible = bot.WithVisible
var WithVisibleFunc = bot.WithVisibleFunc
var WithMiddlewares = bot.WithMiddlewares
