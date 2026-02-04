package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

const smartappAPIVersion = 1

type smartappsListResponse struct {
	Status string `json:"status"`
	Result struct {
		PhonebookVersion int `json:"phonebook_version"`
		SmartApps        []struct {
			AppID         string    `json:"app_id"`
			Enabled       bool      `json:"enabled"`
			ID            uuid.UUID `json:"id"`
			Name          string    `json:"name"`
			Avatar        string    `json:"avatar"`
			AvatarPreview string    `json:"avatar_preview"`
		} `json:"smartapps"`
	} `json:"result"`
}

// SmartAppsList returns list of smartapps.
func (c *Client) SmartAppsList(ctx context.Context, botID uuid.UUID, version optional.Optional[int]) ([]models.SmartApp, int, error) {
	path := "/api/v3/botx/smartapps/list"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return nil, 0, err
	}
	if version.Set {
		urlStr = urlStr + "?" + url.Values{"version": []string{fmt.Sprintf("%d", version.Value)}}.Encode()
	}
	payloadReq, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, 0, err
	}
	var resp smartappsListResponse
	_, err = c.doAuthorizedJSON(ctx, botID, payloadReq, &resp)
	if err != nil {
		return nil, 0, err
	}
	if resp.Status != "ok" {
		return nil, 0, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	apps := make([]models.SmartApp, 0, len(resp.Result.SmartApps))
	for _, a := range resp.Result.SmartApps {
		apps = append(apps, models.SmartApp{AppID: a.AppID, Enabled: a.Enabled, ID: a.ID, Name: a.Name, Avatar: a.Avatar, AvatarPreview: a.AvatarPreview})
	}
	return apps, resp.Result.PhonebookVersion, nil
}

// SendSmartAppEvent sends event to smartapp.
func (c *Client) SendSmartAppEvent(ctx context.Context, botID uuid.UUID, ref optional.Optional[uuid.UUID], smartappID uuid.UUID, chatID uuid.UUID, data map[string]any, opts optional.Optional[map[string]any], files optional.Optional[[]models.AsyncFile], encrypted bool) error {
	path := "/api/v3/botx/smartapps/event"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"smartapp_id":          smartappID,
		"group_chat_id":        chatID,
		"data":                 data,
		"smartapp_api_version": smartappAPIVersion,
		"encrypted":            encrypted,
	}
	if ref.Set {
		payload["ref"] = ref.Value
	}
	if opts.Set {
		payload["opts"] = opts.Value
	} else {
		payload["opts"] = map[string]any{}
	}
	if files.Set {
		asyncFiles := make([]map[string]any, 0, len(files.Value))
		for _, f := range files.Value {
			asyncFiles = append(asyncFiles, map[string]any{
				"type":           string(f.Type),
				"file":           f.FileURL,
				"file_mime_type": f.MimeType,
				"file_id":        f.FileID,
				"file_name":      f.FileName,
				"file_size":      f.FileSize,
				"file_hash":      f.FileHash,
				"duration":       f.Duration,
			})
		}
		payload["async_files"] = asyncFiles
	}

	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return err
	}
	var resp struct {
		Status string `json:"status"`
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

// SendSmartAppNotification sends smartapp notification.
func (c *Client) SendSmartAppNotification(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, counter int, body optional.Optional[string], opts optional.Optional[map[string]any], meta optional.Optional[map[string]any]) error {
	path := "/api/v3/botx/smartapps/notification"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"group_chat_id":        chatID,
		"smartapp_counter":     counter,
		"smartapp_api_version": smartappAPIVersion,
	}
	if body.Set {
		payload["body"] = body.Value
	}
	if opts.Set {
		payload["opts"] = opts.Value
	} else {
		payload["opts"] = map[string]any{}
	}
	if meta.Set {
		payload["meta"] = meta.Value
	}
	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return err
	}
	var resp struct {
		Status string `json:"status"`
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

// SendSmartAppCustomNotification sends custom notification and returns sync_id.
func (c *Client) SendSmartAppCustomNotification(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, title string, body string, meta optional.Optional[map[string]any]) (uuid.UUID, error) {
	path := "/api/v4/botx/smartapps/notification"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return uuid.Nil, err
	}
	payload := map[string]any{
		"group_chat_id": chatID,
		"payload":       map[string]any{"title": title, "body": body},
	}
	if meta.Set {
		payload["meta"] = meta.Value
	}
	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return uuid.Nil, err
	}
	var resp struct {
		Status string `json:"status"`
		Result struct {
			SyncID uuid.UUID `json:"sync_id"`
		} `json:"result"`
	}
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return uuid.Nil, err
	}
	if resp.Status != "ok" {
		return uuid.Nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return resp.Result.SyncID, nil
}

// SendSmartAppUnreadCounter sets unread counter and returns sync_id.
func (c *Client) SendSmartAppUnreadCounter(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, counter int) (uuid.UUID, error) {
	path := "/api/v4/botx/smartapps/unread_counter"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return uuid.Nil, err
	}
	payload := map[string]any{"group_chat_id": chatID, "counter": counter}
	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return uuid.Nil, err
	}
	var resp struct {
		Status string `json:"status"`
		Result struct {
			SyncID uuid.UUID `json:"sync_id"`
		} `json:"result"`
	}
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return uuid.Nil, err
	}
	if resp.Status != "ok" {
		return uuid.Nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return resp.Result.SyncID, nil
}

// SendSmartAppManifest sends manifest payload and returns current manifest.
func (c *Client) SendSmartAppManifest(ctx context.Context, botID uuid.UUID, manifest models.SmartAppManifest) (models.SmartAppManifest, error) {
	path := "/api/v1/botx/smartapps/manifest"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.SmartAppManifest{}, err
	}
	payload := map[string]any{"manifest": manifest}
	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return models.SmartAppManifest{}, err
	}
	var resp struct {
		Status string                  `json:"status"`
		Result models.SmartAppManifest `json:"result"`
	}
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.SmartAppManifest{}, err
	}
	if resp.Status != "ok" {
		return models.SmartAppManifest{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return resp.Result, nil
}

// SmartAppUploadFile uploads file for smartapp and returns link.
func (c *Client) SmartAppUploadFile(ctx context.Context, botID uuid.UUID, filename string, r io.Reader) (string, error) {
	path := "/api/v3/botx/smartapps/upload_file"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	filePart, err := writer.CreateFormFile("content", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(filePart, r); err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, urlStr, &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var resp struct {
		Status string `json:"status"`
		Result struct {
			Link string `json:"link"`
		} `json:"result"`
	}
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return "", err
	}
	if resp.Status != "ok" {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return resp.Result.Link, nil
}
