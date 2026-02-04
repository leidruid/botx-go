package main

import (
	"context"
	"log"
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

	log.Fatal(http.ListenAndServe(":8080", mux))
}
