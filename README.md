# botx-go

Go SDK for building chat-bots and SmartApps for eXpress BotX.

## Quickstart

Create a bot with a command handler and HTTP endpoints.

```go
package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/botx"
)

func main() {
	collector := botx.NewHandlerCollector()
	_ = collector.OnCommand("/echo", func(ctx context.Context, msg botx.IncomingMessage, b *botx.Bot) error {
		_, _ = b.SendMessage(ctx, msg.Bot.ID, msg.Chat.ID, msg.Body, botx.SendMessageOptions{}, false, 0)
		return nil
	}, botx.WithDescription("Echo back text"))

	bot := botx.NewBot([]botx.BotAccountWithSecret{
		{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426655440000"), CTSURL: "https://cts.example.com", SecretKey: "secret"},
	}, collector)

	mux := http.NewServeMux()
	mux.HandleFunc("/command", bot.CommandHandler(true, nil))
	mux.HandleFunc("/smartapps/request", bot.SyncSmartAppHandler(true, nil))
	mux.HandleFunc("/status", bot.StatusHandler(true, nil))
	mux.HandleFunc("/notification/callback", bot.CallbackHandler(false, nil))

	_ = http.ListenAndServe(":8080", mux)
}
```

## Examples

- `examples/internal/echobot/main.go`
- `examples/internal/smartapp_sync/main.go`
- `examples/internal/callback_listener/main.go`
- `examples/simple_auth_send/main.go`

## Notes

- Optional fields use `optional.Optional[T]`.
- Callback-based methods can be called with `waitCallback=true` in `bot` methods.

## API coverage

| Domain | SDK Method | Path |
| --- | --- | --- |
| Bots | `Bot.BotsList` | `/api/v1/botx/bots/catalog` |
| Notifications | `Bot.SendMessage` | `/api/v4/botx/notifications/direct` |
| Events | `Bot.ReplyEvent` | `/api/v3/botx/events/reply_event` |
| Events | `Bot.EditEvent` | `/api/v3/botx/events/edit_event` |
| Events | `Bot.DeleteEvent` | `/api/v3/botx/events/delete_event` |
| Events | `Bot.Typing` | `/api/v3/botx/events/typing` |
| Events | `Bot.StopTyping` | `/api/v3/botx/events/stop_typing` |
| Events | `Bot.MessageStatus` | `/api/v3/botx/events/{sync_id}/status` |
| Chats | `Bot.CreateChat` | `/api/v3/botx/chats/create` |
| Chats | `Bot.ChatInfo` | `/api/v3/botx/chats/info` |
| Chats | `Bot.ListChats` | `/api/v3/botx/chats/list` |
| Chats | `Bot.AddUser` | `/api/v3/botx/chats/add_user` |
| Chats | `Bot.RemoveUser` | `/api/v3/botx/chats/remove_user` |
| Chats | `Bot.AddAdmin` | `/api/v3/botx/chats/add_admin` |
| Chats | `Bot.CreateThread` | `/api/v3/botx/chats/create_thread` |
| Chats | `Bot.SetStealth` | `/api/v3/botx/chats/stealth_set` |
| Chats | `Bot.DisableStealth` | `/api/v3/botx/chats/stealth_disable` |
| Chats | `Bot.PinMessage` | `/api/v3/botx/chats/pin_message` |
| Chats | `Bot.UnpinMessage` | `/api/v3/botx/chats/unpin_message` |
| Users | `Bot.SearchUserByHUID` | `/api/v3/botx/users/by_huid` |
| Users | `Bot.SearchUserByEmail` | `/api/v3/botx/users/by_email` |
| Users | `Bot.SearchUserByEmails` | `/api/v3/botx/users/by_email` (POST) |
| Users | `Bot.SearchUserByLogin` | `/api/v3/botx/users/by_login` |
| Users | `Bot.SearchUserByOtherID` | `/api/v3/botx/users/by_other_id` |
| Users | `Bot.UpdateUserProfile` | `/api/v3/botx/users/update_profile` |
| Users | `Bot.UsersAsCSV` | `/api/v3/botx/users/users_as_csv` |
| Files | `Bot.UploadFile` | `/api/v3/botx/files/upload` |
| Files | `Bot.DownloadFile` | `/api/v3/botx/files/download` |
| SmartApps | `Bot.SmartAppsList` | `/api/v3/botx/smartapps/list` |
| SmartApps | `Bot.SendSmartAppEvent` | `/api/v3/botx/smartapps/event` |
| SmartApps | `Bot.SendSmartAppNotification` | `/api/v3/botx/smartapps/notification` |
| SmartApps | `Bot.SendSmartAppCustomNotification` | `/api/v4/botx/smartapps/notification` |
| SmartApps | `Bot.SendSmartAppUnreadCounter` | `/api/v4/botx/smartapps/unread_counter` |
| SmartApps | `Bot.SendSmartAppManifest` | `/api/v1/botx/smartapps/manifest` |
| SmartApps | `Bot.SmartAppUploadFile` | `/api/v3/botx/smartapps/upload_file` |
| Stickers | `Bot.CreateStickerPack` | `/api/v3/botx/stickers/packs` |
| Stickers | `Bot.EditStickerPack` | `/api/v3/botx/stickers/packs/{id}` |
| Stickers | `Bot.DeleteStickerPack` | `/api/v3/botx/stickers/packs/{id}` |
| Stickers | `Bot.GetStickerPack` | `/api/v3/botx/stickers/packs/{id}` |
| Stickers | `Bot.GetStickerPacks` | `/api/v3/botx/stickers/packs` |
| Stickers | `Bot.GetSticker` | `/api/v3/botx/stickers/packs/{id}/stickers/{id}` |
| Stickers | `Bot.AddSticker` | `/api/v3/botx/stickers/packs/{id}/stickers` |
| Stickers | `Bot.DeleteSticker` | `/api/v3/botx/stickers/packs/{id}/stickers/{id}` |
| OpenID | `Bot.RefreshAccessToken` | `/api/v3/botx/openid/refresh_access_token` |
| Metrics | `Bot.CollectBotFunction` | `/api/v3/botx/metrics/bot_function` |
