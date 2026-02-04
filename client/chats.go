package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

type createChatRequest struct {
	Name          string                  `json:"name"`
	Description   string                  `json:"description,omitempty"`
	ChatType      models.ChatType         `json:"chat_type"`
	Members       []uuid.UUID             `json:"members"`
	SharedHistory optional.Optional[bool] `json:"shared_history,omitempty"`
}

type createChatResponse struct {
	Status string `json:"status"`
	Result struct {
		ChatID uuid.UUID `json:"chat_id"`
	} `json:"result"`
}

type chatInfoResponse struct {
	Status string `json:"status"`
	Result struct {
		ChatType    models.ChatType `json:"chat_type"`
		Creator     *uuid.UUID      `json:"creator"`
		Description string          `json:"description"`
		GroupChatID uuid.UUID       `json:"group_chat_id"`
		InsertedAt  time.Time       `json:"inserted_at"`
		Members     []struct {
			Admin    bool      `json:"admin"`
			UserHUID uuid.UUID `json:"user_huid"`
			UserKind string    `json:"user_kind"`
		} `json:"members"`
		Name          string `json:"name"`
		SharedHistory bool   `json:"shared_history"`
	} `json:"result"`
}

type listChatsResponse struct {
	Status string `json:"status"`
	Result []struct {
		GroupChatID   uuid.UUID       `json:"group_chat_id"`
		ChatType      models.ChatType `json:"chat_type"`
		Name          string          `json:"name"`
		Description   string          `json:"description"`
		Members       []uuid.UUID     `json:"members"`
		InsertedAt    time.Time       `json:"inserted_at"`
		UpdatedAt     time.Time       `json:"updated_at"`
		SharedHistory bool            `json:"shared_history"`
	} `json:"result"`
}

// CreateChat creates a new chat.
func (c *Client) CreateChat(ctx context.Context, botID uuid.UUID, name string, chatType models.ChatType, members []uuid.UUID, sharedHistory optional.Optional[bool], description string) (uuid.UUID, error) {
	path := "/api/v3/botx/chats/create"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return uuid.Nil, err
	}
	payload := createChatRequest{
		Name:          name,
		ChatType:      chatType,
		Members:       members,
		SharedHistory: sharedHistory,
		Description:   description,
	}
	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return uuid.Nil, err
	}
	var resp createChatResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return uuid.Nil, err
	}
	if resp.Status != "ok" {
		return uuid.Nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return resp.Result.ChatID, nil
}

// ChatInfo returns chat details.
func (c *Client) ChatInfo(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) (models.ChatInfo, error) {
	path := "/api/v3/botx/chats/info"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.ChatInfo{}, err
	}
	q := fmt.Sprintf("?group_chat_id=%s", chatID.String())
	req, err := http.NewRequest(http.MethodGet, urlStr+q, nil)
	if err != nil {
		return models.ChatInfo{}, err
	}
	var resp chatInfoResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.ChatInfo{}, err
	}
	if resp.Status != "ok" {
		return models.ChatInfo{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	members := make([]models.ChatInfoMember, 0, len(resp.Result.Members))
	for _, m := range resp.Result.Members {
		members = append(members, models.ChatInfoMember{HUID: m.UserHUID, IsAdmin: m.Admin, Kind: m.UserKind})
	}
	return models.ChatInfo{
		ChatID:        resp.Result.GroupChatID,
		ChatType:      resp.Result.ChatType,
		Name:          resp.Result.Name,
		Description:   resp.Result.Description,
		CreatorID:     resp.Result.Creator,
		CreatedAt:     resp.Result.InsertedAt,
		Members:       members,
		SharedHistory: resp.Result.SharedHistory,
	}, nil
}

// ListChats returns list of chats.
func (c *Client) ListChats(ctx context.Context, botID uuid.UUID) ([]models.ChatListItem, error) {
	path := "/api/v3/botx/chats/list"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	var resp listChatsResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != "ok" {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	items := make([]models.ChatListItem, 0, len(resp.Result))
	for _, citem := range resp.Result {
		items = append(items, models.ChatListItem{
			ChatID:        citem.GroupChatID,
			ChatType:      citem.ChatType,
			Name:          citem.Name,
			Description:   citem.Description,
			Members:       citem.Members,
			CreatedAt:     citem.InsertedAt,
			UpdatedAt:     citem.UpdatedAt,
			SharedHistory: citem.SharedHistory,
		})
	}
	return items, nil
}

// AddUser adds users to chat.
func (c *Client) AddUser(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, users []uuid.UUID) error {
	path := "/api/v3/botx/chats/add_user"
	return c.simpleChatUsersRequest(ctx, botID, path, chatID, users)
}

// RemoveUser removes users from chat.
func (c *Client) RemoveUser(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, users []uuid.UUID) error {
	path := "/api/v3/botx/chats/remove_user"
	return c.simpleChatUsersRequest(ctx, botID, path, chatID, users)
}

// AddAdmin adds admins to chat.
func (c *Client) AddAdmin(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, users []uuid.UUID) error {
	path := "/api/v3/botx/chats/add_admin"
	return c.simpleChatUsersRequest(ctx, botID, path, chatID, users)
}

func (c *Client) simpleChatUsersRequest(ctx context.Context, botID uuid.UUID, path string, chatID uuid.UUID, users []uuid.UUID) error {
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{"group_chat_id": chatID, "user_huids": users}
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

// CreateThread creates a thread from message sync_id.
func (c *Client) CreateThread(ctx context.Context, botID uuid.UUID, syncID uuid.UUID) (uuid.UUID, error) {
	path := "/api/v3/botx/chats/create_thread"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return uuid.Nil, err
	}
	payload := map[string]any{"sync_id": syncID}
	req, err := BuildJSONRequest(http.MethodPost, urlStr, payload)
	if err != nil {
		return uuid.Nil, err
	}
	var resp struct {
		Status string `json:"status"`
		Result struct {
			ThreadID uuid.UUID `json:"thread_id"`
		} `json:"result"`
	}
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return uuid.Nil, err
	}
	if resp.Status != "ok" {
		return uuid.Nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return resp.Result.ThreadID, nil
}

// SetStealth enables stealth mode.
func (c *Client) SetStealth(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, disableWeb optional.Optional[bool], burnIn optional.Optional[int], expireIn optional.Optional[int]) error {
	path := "/api/v3/botx/chats/stealth_set"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"group_chat_id": chatID,
	}
	if disableWeb.Set {
		payload["disable_web"] = disableWeb.Value
	}
	if burnIn.Set {
		payload["burn_in"] = burnIn.Value
	}
	if expireIn.Set {
		payload["expire_in"] = expireIn.Value
	}
	return c.simpleStatusRequest(ctx, botID, urlStr, payload)
}

// DisableStealth disables stealth mode.
func (c *Client) DisableStealth(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) error {
	path := "/api/v3/botx/chats/stealth_disable"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{"group_chat_id": chatID}
	return c.simpleStatusRequest(ctx, botID, urlStr, payload)
}

// PinMessage pins a message in chat.
func (c *Client) PinMessage(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, syncID uuid.UUID) error {
	path := "/api/v3/botx/chats/pin_message"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{"chat_id": chatID, "sync_id": syncID}
	return c.simpleStatusRequest(ctx, botID, urlStr, payload)
}

// UnpinMessage unpins a message.
func (c *Client) UnpinMessage(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) error {
	path := "/api/v3/botx/chats/unpin_message"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{"chat_id": chatID}
	return c.simpleStatusRequest(ctx, botID, urlStr, payload)
}

func (c *Client) simpleStatusRequest(ctx context.Context, botID uuid.UUID, urlStr string, payload map[string]any) error {
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
