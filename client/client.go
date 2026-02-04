package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/accounts"
	"github.com/leidruid/botx-go/errors"
)

const defaultTimeout = 15 * time.Second

// Client is a low-level BotX API client.
type Client struct {
	HTTP     *http.Client
	Accounts *accounts.Storage
}

func NewClient(accounts *accounts.Storage) *Client {
	return &Client{
		HTTP:     &http.Client{Timeout: defaultTimeout},
		Accounts: accounts,
	}
}

func (c *Client) buildURL(botID uuid.UUID, path string) (string, error) {
	base, err := c.Accounts.CTSURL(botID)
	if err != nil {
		return "", err
	}
	return joinURL(base, path), nil
}

func joinURL(base, path string) string {
	return fmt.Sprintf("%s/%s", trimSlash(base), trimSlash(path))
}

func trimSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	for len(s) > 0 && s[0] == '/' {
		s = s[1:]
	}
	return s
}

func (c *Client) doJSON(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return resp, &errors.APIError{StatusCode: resp.StatusCode, Message: string(body)}
	}

	if v == nil {
		return resp, nil
	}

	if err := json.Unmarshal(body, v); err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *Client) doAuthorizedJSON(ctx context.Context, botID uuid.UUID, req *http.Request, v any) (*http.Response, error) {
	token, ok := c.Accounts.Token(botID)
	if !ok {
		var err error
		token, err = c.GetToken(ctx, botID)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.doJSON(ctx, req, v)
	if err == nil {
		return resp, nil
	}

	apiErr, ok := err.(*errors.APIError)
	if !ok || apiErr.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	// retry once with fresh token
	fresh, err := c.GetToken(ctx, botID)
	if err != nil {
		return resp, err
	}
	req.Header.Set("Authorization", "Bearer "+fresh)
	return c.doJSON(ctx, req, v)
}

// GetToken requests token for bot.
func (c *Client) GetToken(ctx context.Context, botID uuid.UUID) (string, error) {
	signature, err := c.Accounts.BuildSignature(botID)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("/api/v2/botx/bots/%s/token", botID.String())
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	q.Set("signature", signature)
	urlStr = urlStr + "?" + q.Encode()

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return "", err
	}

	var resp struct {
		Status string `json:"status"`
		Result string `json:"result"`
	}

	_, err = c.doJSON(ctx, req, &resp)
	if err != nil {
		return "", err
	}
	if resp.Status != "ok" {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	c.Accounts.SetToken(botID, resp.Result)
	return resp.Result, nil
}

// BuildJSONRequest helps create a JSON request.
func BuildJSONRequest(method, urlStr string, payload any) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}
