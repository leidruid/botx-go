package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// CollectBotFunction reports bot function usage.
func (c *Client) CollectBotFunction(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, huids []uuid.UUID, botFunction string) error {
	path := "/api/v3/botx/metrics/bot_function"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{"group_chat_id": chatID, "user_huids": huids, "bot_function": botFunction}
	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return err
	}
	var resp struct {
		Status string `json:"status"`
		Result bool   `json:"result"`
	}
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return err
	}
	if resp.Status != "ok" {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
