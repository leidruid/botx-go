package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/botx"
	"github.com/leidruid/botx-go/client"
)

func main() {
	botID := uuid.MustParse("123e4567-e89b-12d3-a456-426655440000")
	chatID := uuid.MustParse("123e4567-e89b-12d3-a456-426655440001")

	bot := botx.NewBot([]botx.BotAccountWithSecret{
		{ID: botID, CTSURL: "https://cts.example.com", SecretKey: "secret"},
	})

	ctx := context.Background()

	// Optional: explicitly request token (authorization)
	_, err := bot.Client.GetToken(ctx, botID)
	if err != nil {
		log.Fatalf("get token failed: %v", err)
	}

	// Send message
	_, err = bot.SendMessage(ctx, botID, chatID, "Hello from botx-go", client.SendMessageOptions{}, false, 0)
	if err != nil {
		log.Fatalf("send message failed: %v", err)
	}
}
