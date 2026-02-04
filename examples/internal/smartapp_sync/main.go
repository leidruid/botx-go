package main

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/botx"
	"github.com/leidruid/botx-go/models"
)

func main() {
	collector := botx.NewHandlerCollector()
	collector.OnSyncSmartApp(func(ctx context.Context, ev botx.SmartAppEvent, b *botx.Bot) (models.SyncSmartAppResponse, error) {
		return models.SyncSmartAppResponse{
			Status: "ok",
			Result: map[string]any{"echo": ev.Data},
		}, nil
	})

	bot := botx.NewBot([]botx.BotAccountWithSecret{
		{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426655440000"), CTSURL: "https://cts.example.com", SecretKey: "secret"},
	}, collector)

	mux := http.NewServeMux()
	mux.HandleFunc("/smartapps/request", bot.SyncSmartAppHandler(true, nil))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
