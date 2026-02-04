package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

type botsListResponse struct {
	Status string `json:"status"`
	Result struct {
		GeneratedAt time.Time `json:"generated_at"`
		Bots        []struct {
			UserHUID    uuid.UUID `json:"user_huid"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			Avatar      string    `json:"avatar"`
			Enabled     bool      `json:"enabled"`
		} `json:"bots"`
	} `json:"result"`
}

// BotsList returns bots catalog.
func (c *Client) BotsList(ctx context.Context, botID uuid.UUID, since optional.Optional[time.Time]) ([]models.BotsListItem, time.Time, error) {
	path := "/api/v1/botx/bots/catalog"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return nil, time.Time{}, err
	}
	q := url.Values{}
	if since.Set {
		q.Set("since", since.Value.Format(time.RFC3339))
	}
	if enc := q.Encode(); enc != "" {
		urlStr = urlStr + "?" + enc
	}

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, time.Time{}, err
	}
	var resp botsListResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return nil, time.Time{}, err
	}
	if resp.Status != "ok" {
		return nil, time.Time{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	items := make([]models.BotsListItem, 0, len(resp.Result.Bots))
	for _, b := range resp.Result.Bots {
		items = append(items, models.BotsListItem{
			ID:          b.UserHUID,
			Name:        b.Name,
			Description: b.Description,
			Avatar:      b.Avatar,
			Enabled:     b.Enabled,
		})
	}
	return items, resp.Result.GeneratedAt, nil
}
