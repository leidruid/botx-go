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

type replyEventRequest struct {
	SourceSyncID uuid.UUID `json:"source_sync_id"`
	Reply        struct {
		Status   string                                     `json:"status"`
		Body     string                                     `json:"body"`
		Metadata optional.Optional[map[string]any]          `json:"metadata,omitempty"`
		Opts     optional.Optional[notificationMessageOpts] `json:"opts,omitempty"`
		Bubble   optional.Optional[models.APIMarkup]        `json:"bubble,omitempty"`
		Keyboard optional.Optional[models.APIMarkup]        `json:"keyboard,omitempty"`
		Mentions optional.Optional[[]models.APIMention]     `json:"mentions,omitempty"`
	} `json:"reply"`
	File optional.Optional[models.APINotificationAttachment] `json:"file,omitempty"`
	Opts struct {
		RawMentions      bool                                      `json:"raw_mentions"`
		StealthMode      optional.Optional[bool]                   `json:"stealth_mode,omitempty"`
		NotificationOpts optional.Optional[notificationNestedOpts] `json:"notification_opts,omitempty"`
	} `json:"opts"`
}

type editEventRequest struct {
	SyncID  uuid.UUID `json:"sync_id"`
	Payload struct {
		Body     optional.Optional[string]         `json:"body,omitempty"`
		Metadata optional.Optional[map[string]any] `json:"metadata,omitempty"`
		Opts     optional.Optional[struct {
			ButtonsAutoAdjust optional.Optional[bool] `json:"buttons_auto_adjust,omitempty"`
		}] `json:"opts,omitempty"`
		Bubble   optional.Optional[models.APIMarkup]    `json:"bubble,omitempty"`
		Keyboard optional.Optional[models.APIMarkup]    `json:"keyboard,omitempty"`
		Mentions optional.Optional[[]models.APIMention] `json:"mentions,omitempty"`
	} `json:"payload"`
	File optional.Optional[models.APINotificationAttachment] `json:"file,omitempty"`
	Opts optional.Optional[struct {
		RawMentions optional.Optional[bool] `json:"raw_mentions,omitempty"`
	}] `json:"opts,omitempty"`
}

type statusResponse struct {
	Status string `json:"status"`
	Result struct {
		GroupChatID uuid.UUID   `json:"group_chat_id"`
		SentTo      []uuid.UUID `json:"sent_to"`
		ReadBy      []struct {
			UserHUID uuid.UUID `json:"user_huid"`
			ReadAt   time.Time `json:"read_at"`
		} `json:"read_by"`
		ReceivedBy []struct {
			UserHUID   uuid.UUID `json:"user_huid"`
			ReceivedAt time.Time `json:"received_at"`
		} `json:"received_by"`
	} `json:"result"`
}

// ReplyEvent replies to an existing message by sync_id.
func (c *Client) ReplyEvent(ctx context.Context, botID uuid.UUID, sourceSyncID uuid.UUID, body string, opts SendMessageOptions) error {
	path := "/api/v3/botx/events/reply_event"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}

	cleanBody, mentions := models.FindAndReplaceEmbedMentions(body)

	var reqPayload replyEventRequest
	reqPayload.SourceSyncID = sourceSyncID
	reqPayload.Reply.Status = "ok"
	reqPayload.Reply.Body = cleanBody
	reqPayload.Reply.Metadata = opts.Metadata

	if opts.Bubble.Set && opts.Bubble.Value != nil {
		reqPayload.Reply.Bubble = optional.Some(opts.Bubble.Value.ToAPI())
	}
	if opts.Keyboard.Set && opts.Keyboard.Value != nil {
		reqPayload.Reply.Keyboard = optional.Some(opts.Keyboard.Value.ToAPI())
	}
	if len(mentions) > 0 {
		reqPayload.Reply.Mentions = optional.Some(mentions)
	}

	msgOpts := notificationMessageOpts{
		SilentResponse:    opts.SilentResponse,
		ButtonsAutoAdjust: opts.MarkupAutoAdjust,
	}
	if msgOpts.SilentResponse.Set || msgOpts.ButtonsAutoAdjust.Set {
		reqPayload.Reply.Opts = optional.Some(msgOpts)
	}

	if opts.Attachment.Set {
		apiAtt, err := models.AttachmentToAPI(opts.Attachment.Value)
		if err != nil {
			return err
		}
		reqPayload.File = optional.Some(apiAtt)
	}

	reqPayload.Opts.RawMentions = true
	if opts.StealthMode.Set || opts.SendPush.Set || opts.IgnoreMute.Set {
		nested := notificationNestedOpts{Send: opts.SendPush, ForceDND: opts.IgnoreMute}
		reqPayload.Opts.NotificationOpts = optional.Some(nested)
		reqPayload.Opts.StealthMode = opts.StealthMode
	}

	req, err := BuildJSONRequest(http.MethodPost, urlStr, reqPayload)
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

// EditEvent edits an existing message by sync_id.
func (c *Client) EditEvent(ctx context.Context, botID uuid.UUID, syncID uuid.UUID, body optional.Optional[string], metadata optional.Optional[map[string]any], bubble optional.Optional[*models.Markup], keyboard optional.Optional[*models.Markup], attachment optional.Optional[*models.OutgoingAttachment], markupAutoAdjust optional.Optional[bool]) error {
	path := "/api/v3/botx/events/edit_event"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}

	mentions := []models.APIMention{}
	if body.Set {
		cleanBody, m := models.FindAndReplaceEmbedMentions(body.Value)
		body = optional.Some(cleanBody)
		mentions = m
	}

	var reqPayload editEventRequest
	reqPayload.SyncID = syncID
	reqPayload.Payload.Body = body
	reqPayload.Payload.Metadata = metadata
	if bubble.Set && bubble.Value != nil {
		reqPayload.Payload.Bubble = optional.Some(bubble.Value.ToAPI())
	}
	if keyboard.Set && keyboard.Value != nil {
		reqPayload.Payload.Keyboard = optional.Some(keyboard.Value.ToAPI())
	}
	if len(mentions) > 0 {
		reqPayload.Payload.Mentions = optional.Some(mentions)
		reqPayload.Opts = optional.Some(struct {
			RawMentions optional.Optional[bool] `json:"raw_mentions,omitempty"`
		}{RawMentions: optional.Some(true)})
	}

	if markupAutoAdjust.Set {
		reqPayload.Payload.Opts = optional.Some(struct {
			ButtonsAutoAdjust optional.Optional[bool] `json:"buttons_auto_adjust,omitempty"`
		}{ButtonsAutoAdjust: markupAutoAdjust})
	}

	if attachment.Set && attachment.Value != nil {
		apiAtt, err := models.AttachmentToAPI(*attachment.Value)
		if err != nil {
			return err
		}
		reqPayload.File = optional.Some(apiAtt)
	} else if attachment.Set && attachment.Value == nil {
		// Explicitly remove file by sending null is not supported in this simplified client.
	}

	req, err := BuildJSONRequest(http.MethodPost, urlStr, reqPayload)
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

// DeleteEvent deletes a message by sync_id.
func (c *Client) DeleteEvent(ctx context.Context, botID uuid.UUID, syncID uuid.UUID) error {
	path := "/api/v3/botx/events/delete_event"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{"sync_id": syncID}
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

// Typing sends typing event.
func (c *Client) Typing(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) error {
	path := "/api/v3/botx/events/typing"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{"group_chat_id": chatID}
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

// StopTyping stops typing event.
func (c *Client) StopTyping(ctx context.Context, botID uuid.UUID, chatID uuid.UUID) error {
	path := "/api/v3/botx/events/stop_typing"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	payload := map[string]any{"group_chat_id": chatID}
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

// MessageStatus returns message delivery/read status.
func (c *Client) MessageStatus(ctx context.Context, botID uuid.UUID, syncID uuid.UUID) (models.MessageStatus, error) {
	path := fmt.Sprintf("/api/v3/botx/events/%s/status", syncID.String())
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.MessageStatus{}, err
	}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return models.MessageStatus{}, err
	}
	var resp statusResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.MessageStatus{}, err
	}
	if resp.Status != "ok" {
		return models.MessageStatus{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	readBy := make(map[uuid.UUID]time.Time, len(resp.Result.ReadBy))
	for _, r := range resp.Result.ReadBy {
		readBy[r.UserHUID] = r.ReadAt
	}
	receivedBy := make(map[uuid.UUID]time.Time, len(resp.Result.ReceivedBy))
	for _, r := range resp.Result.ReceivedBy {
		receivedBy[r.UserHUID] = r.ReceivedAt
	}

	return models.MessageStatus{
		GroupChatID: resp.Result.GroupChatID,
		SentTo:      resp.Result.SentTo,
		ReadBy:      readBy,
		ReceivedBy:  receivedBy,
	}, nil
}
