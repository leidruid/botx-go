package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/botx"
)

func main() {
	bot := botx.NewBot([]botx.BotAccountWithSecret{
		{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426655440000"), CTSURL: "https://cts.example.com", SecretKey: "secret"},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/notification/callback", bot.CallbackHandler(false, nil))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
