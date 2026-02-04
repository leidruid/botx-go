package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// RefreshAccessToken refreshes OpenID access token.
func (c *Client) RefreshAccessToken(ctx context.Context, botID uuid.UUID, userHUID uuid.UUID, ref *uuid.UUID) (bool, error) {
	path := "/api/v3/botx/openid/refresh_access_token"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return false, err
	}
	payload := map[string]any{"user_huid": userHUID}
	if ref != nil {
		payload["ref"] = *ref
	}
	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return false, err
	}
	var resp struct {
		Status string `json:"status"`
		Result bool   `json:"result"`
	}
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return false, err
	}
	if resp.Status != "ok" {
		return false, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return resp.Result, nil
}
