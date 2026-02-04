package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

type createStickerPackResponse struct {
	Status string `json:"status"`
	Result struct {
		ID     uuid.UUID `json:"id"`
		Name   string    `json:"name"`
		Public bool      `json:"public"`
	} `json:"result"`
}

type stickerPackResult struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Public   bool      `json:"public"`
	Stickers []struct {
		ID    uuid.UUID `json:"id"`
		Emoji string    `json:"emoji"`
		Link  string    `json:"link"`
	} `json:"stickers"`
}

type getStickerPackResponse struct {
	Status string            `json:"status"`
	Result stickerPackResult `json:"result"`
}

type getStickerPacksResponse struct {
	Status string `json:"status"`
	Result struct {
		Packs []struct {
			ID            uuid.UUID   `json:"id"`
			Name          string      `json:"name"`
			Public        bool        `json:"public"`
			StickersCount int         `json:"stickers_count"`
			StickersOrder []uuid.UUID `json:"stickers_order"`
		} `json:"packs"`
		Pagination struct {
			After string `json:"after"`
		} `json:"pagination"`
	} `json:"result"`
}

type getStickerResponse struct {
	Status string `json:"status"`
	Result struct {
		ID    uuid.UUID `json:"id"`
		Emoji string    `json:"emoji"`
		Link  string    `json:"link"`
	} `json:"result"`
}

// CreateStickerPack creates sticker pack.
func (c *Client) CreateStickerPack(ctx context.Context, botID uuid.UUID, name string, userHUID optional.Optional[uuid.UUID]) (models.StickerPack, error) {
	path := "/api/v3/botx/stickers/packs"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.StickerPack{}, err
	}
	payload := map[string]any{"name": name}
	if userHUID.Set {
		payload["user_huid"] = userHUID.Value
	}

	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return models.StickerPack{}, err
	}
	var resp createStickerPackResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.StickerPack{}, err
	}
	if resp.Status != "ok" {
		return models.StickerPack{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return models.StickerPack{ID: resp.Result.ID, Name: resp.Result.Name, IsPublic: resp.Result.Public, Stickers: []models.Sticker{}}, nil
}

// EditStickerPack edits sticker pack.
func (c *Client) EditStickerPack(ctx context.Context, botID uuid.UUID, packID uuid.UUID, name string, preview uuid.UUID, stickersOrder []uuid.UUID) (models.StickerPack, error) {
	path := fmt.Sprintf("/api/v3/botx/stickers/packs/%s", packID.String())
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.StickerPack{}, err
	}
	payload := map[string]any{"name": name, "preview": preview}
	if stickersOrder != nil {
		payload["stickers_order"] = stickersOrder
	}
	req, err := BuildJSONRequest(http.MethodPut, urlStr, payload)
	if err != nil {
		return models.StickerPack{}, err
	}
	var resp getStickerPackResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.StickerPack{}, err
	}
	if resp.Status != "ok" {
		return models.StickerPack{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return mapStickerPack(resp.Result), nil
}

// DeleteStickerPack deletes sticker pack.
func (c *Client) DeleteStickerPack(ctx context.Context, botID uuid.UUID, packID uuid.UUID) error {
	path := fmt.Sprintf("/api/v3/botx/stickers/packs/%s", packID.String())
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, urlStr, nil)
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

// GetStickerPack returns sticker pack by id.
func (c *Client) GetStickerPack(ctx context.Context, botID uuid.UUID, packID uuid.UUID) (models.StickerPack, error) {
	path := fmt.Sprintf("/api/v3/botx/stickers/packs/%s", packID.String())
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.StickerPack{}, err
	}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return models.StickerPack{}, err
	}
	var resp getStickerPackResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.StickerPack{}, err
	}
	if resp.Status != "ok" {
		return models.StickerPack{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return mapStickerPack(resp.Result), nil
}

// GetStickerPacks returns sticker packs list.
func (c *Client) GetStickerPacks(ctx context.Context, botID uuid.UUID, userHUID uuid.UUID, limit int, after optional.Optional[string]) (models.StickerPackPage, error) {
	path := "/api/v3/botx/stickers/packs"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.StickerPackPage{}, err
	}
	params := url.Values{"user_huid": []string{userHUID.String()}, "limit": []string{fmt.Sprintf("%d", limit)}}
	if after.Set {
		params.Set("after", after.Value)
	}
	urlStr = urlStr + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return models.StickerPackPage{}, err
	}
	var resp getStickerPacksResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.StickerPackPage{}, err
	}
	if resp.Status != "ok" {
		return models.StickerPackPage{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	packs := make([]models.StickerPackFromList, 0, len(resp.Result.Packs))
	for _, p := range resp.Result.Packs {
		packs = append(packs, models.StickerPackFromList{ID: p.ID, Name: p.Name, IsPublic: p.Public, StickersCount: p.StickersCount, StickerIDs: p.StickersOrder})
	}
	return models.StickerPackPage{StickerPacks: packs, After: resp.Result.Pagination.After}, nil
}

// GetSticker returns sticker by pack and sticker id.
func (c *Client) GetSticker(ctx context.Context, botID uuid.UUID, packID uuid.UUID, stickerID uuid.UUID) (models.Sticker, error) {
	path := fmt.Sprintf("/api/v3/botx/stickers/packs/%s/stickers/%s", packID.String(), stickerID.String())
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.Sticker{}, err
	}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return models.Sticker{}, err
	}
	var resp getStickerResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.Sticker{}, err
	}
	if resp.Status != "ok" {
		return models.Sticker{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return models.Sticker{ID: resp.Result.ID, Emoji: resp.Result.Emoji, Link: resp.Result.Link, PackID: packID}, nil
}

// AddSticker adds sticker to pack.
func (c *Client) AddSticker(ctx context.Context, botID uuid.UUID, packID uuid.UUID, emoji string, png []byte) (models.Sticker, error) {
	path := fmt.Sprintf("/api/v3/botx/stickers/packs/%s/stickers", packID.String())
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.Sticker{}, err
	}

	imageData := models.EncodeRFC2397(png, "image/png")
	payload := map[string]any{"emoji": emoji, "image": imageData}
	bodyBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(bodyBytes))
	if err != nil {
		return models.Sticker{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	var resp getStickerResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.Sticker{}, err
	}
	if resp.Status != "ok" {
		return models.Sticker{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return models.Sticker{ID: resp.Result.ID, Emoji: resp.Result.Emoji, Link: resp.Result.Link, PackID: packID}, nil
}

// DeleteSticker deletes sticker from pack.
func (c *Client) DeleteSticker(ctx context.Context, botID uuid.UUID, packID uuid.UUID, stickerID uuid.UUID) error {
	path := fmt.Sprintf("/api/v3/botx/stickers/packs/%s/stickers/%s", packID.String(), stickerID.String())
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, urlStr, nil)
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

func mapStickerPack(r stickerPackResult) models.StickerPack {
	stickers := make([]models.Sticker, 0, len(r.Stickers))
	for _, s := range r.Stickers {
		stickers = append(stickers, models.Sticker{ID: s.ID, Emoji: s.Emoji, Link: s.Link, PackID: r.ID})
	}
	return models.StickerPack{ID: r.ID, Name: r.Name, IsPublic: r.Public, Stickers: stickers}
}
