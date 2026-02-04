package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

type notificationMessageOpts struct {
	SilentResponse    optional.Optional[bool] `json:"silent_response,omitempty"`
	ButtonsAutoAdjust optional.Optional[bool] `json:"buttons_auto_adjust,omitempty"`
}

type notificationNestedOpts struct {
	Send     optional.Optional[bool] `json:"send,omitempty"`
	ForceDND optional.Optional[bool] `json:"force_dnd,omitempty"`
}

type notificationOpts struct {
	StealthMode      optional.Optional[bool]                   `json:"stealth_mode,omitempty"`
	NotificationOpts optional.Optional[notificationNestedOpts] `json:"notification_opts,omitempty"`
}

type directNotificationRequest struct {
	GroupChatID  uuid.UUID `json:"group_chat_id"`
	Notification struct {
		Status   string                                     `json:"status"`
		Body     string                                     `json:"body"`
		Metadata optional.Optional[map[string]any]          `json:"metadata,omitempty"`
		Opts     optional.Optional[notificationMessageOpts] `json:"opts,omitempty"`
		Bubble   optional.Optional[models.APIMarkup]        `json:"bubble,omitempty"`
		Keyboard optional.Optional[models.APIMarkup]        `json:"keyboard,omitempty"`
		Mentions optional.Optional[[]models.APIMention]     `json:"mentions,omitempty"`
	} `json:"notification"`
	File       optional.Optional[models.APINotificationAttachment] `json:"file,omitempty"`
	Recipients optional.Optional[[]uuid.UUID]                      `json:"recipients,omitempty"`
	Opts       optional.Optional[notificationOpts]                 `json:"opts,omitempty"`
}

type directNotificationResponse struct {
	Status string `json:"status"`
	Result struct {
		SyncID uuid.UUID `json:"sync_id"`
	} `json:"result"`
}

// SendMessageOptions configures message sending.
type SendMessageOptions struct {
	Metadata         optional.Optional[map[string]any]
	Bubble           optional.Optional[*models.Markup]
	Keyboard         optional.Optional[*models.Markup]
	Attachment       optional.Optional[models.OutgoingAttachment]
	Recipients       optional.Optional[[]uuid.UUID]
	SilentResponse   optional.Optional[bool]
	MarkupAutoAdjust optional.Optional[bool]
	StealthMode      optional.Optional[bool]
	SendPush         optional.Optional[bool]
	IgnoreMute       optional.Optional[bool]
}

// SendMessage sends a message to a chat and returns sync_id.
func (c *Client) SendMessage(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, body string, opts SendMessageOptions) (uuid.UUID, error) {
	path := "/api/v4/botx/notifications/direct"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return uuid.Nil, err
	}

	cleanBody, mentions := models.FindAndReplaceEmbedMentions(body)

	var reqPayload directNotificationRequest
	reqPayload.GroupChatID = chatID
	reqPayload.Notification.Status = "ok"
	reqPayload.Notification.Body = cleanBody
	reqPayload.Notification.Metadata = opts.Metadata

	if opts.Bubble.Set && opts.Bubble.Value != nil {
		reqPayload.Notification.Bubble = optional.Some(opts.Bubble.Value.ToAPI())
	}
	if opts.Keyboard.Set && opts.Keyboard.Value != nil {
		reqPayload.Notification.Keyboard = optional.Some(opts.Keyboard.Value.ToAPI())
	}
	if len(mentions) > 0 {
		reqPayload.Notification.Mentions = optional.Some(mentions)
	}

	msgOpts := notificationMessageOpts{
		SilentResponse:    opts.SilentResponse,
		ButtonsAutoAdjust: opts.MarkupAutoAdjust,
	}
	if msgOpts.SilentResponse.Set || msgOpts.ButtonsAutoAdjust.Set {
		reqPayload.Notification.Opts = optional.Some(msgOpts)
	}

	if opts.Attachment.Set {
		apiAtt, err := models.AttachmentToAPI(opts.Attachment.Value)
		if err != nil {
			return uuid.Nil, err
		}
		reqPayload.File = optional.Some(apiAtt)
	}

	if opts.Recipients.Set {
		reqPayload.Recipients = opts.Recipients
	}

	if opts.StealthMode.Set || opts.SendPush.Set || opts.IgnoreMute.Set {
		nested := notificationNestedOpts{Send: opts.SendPush, ForceDND: opts.IgnoreMute}
		reqPayload.Opts = optional.Some(notificationOpts{
			StealthMode:      opts.StealthMode,
			NotificationOpts: optional.Some(nested),
		})
	}

	req, err := BuildJSONRequest(http.MethodPost, urlStr, reqPayload)
	if err != nil {
		return uuid.Nil, err
	}
	var resp directNotificationResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return uuid.Nil, err
	}
	if resp.Status != "ok" {
		return uuid.Nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return resp.Result.SyncID, nil
}
