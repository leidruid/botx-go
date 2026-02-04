package bot

import (
	"testing"

	"github.com/google/uuid"
)

func TestParseIncomingMessage(t *testing.T) {
	botID := uuid.MustParse("123e4567-e89b-12d3-a456-426655440000")
	chatID := uuid.MustParse("123e4567-e89b-12d3-a456-426655440001")
	userHUID := uuid.MustParse("123e4567-e89b-12d3-a456-426655440002")

	raw := map[string]any{
		"bot_id":  botID.String(),
		"sync_id": uuid.MustParse("123e4567-e89b-12d3-a456-426655440003").String(),
		"command": map[string]any{
			"body":     "/echo hi",
			"data":     map[string]any{},
			"metadata": map[string]any{},
		},
		"from": map[string]any{
			"group_chat_id": chatID.String(),
			"user_huid":     userHUID.String(),
		},
	}

	msg, err := ParseIncomingMessage(raw)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if msg.Bot.ID != botID || msg.Chat.ID != chatID || msg.Sender.HUID != userHUID {
		t.Fatalf("unexpected ids")
	}
	if msg.Body != "/echo hi" {
		t.Fatalf("unexpected body")
	}
}
