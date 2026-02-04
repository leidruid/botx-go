package bot

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
)

var (
	ErrInvalidPayload = errors.New("invalid payload")
)

// ParseIncomingMessage parses raw BotX command payload into IncomingMessage.
func ParseIncomingMessage(raw map[string]any) (models.IncomingMessage, error) {
	syncID, _ := getUUID(raw, "sync_id")
	botID, ok := getUUID(raw, "bot_id")
	if !ok {
		return models.IncomingMessage{}, fmt.Errorf("%w: bot_id", ErrInvalidPayload)
	}

	command, _ := raw["command"].(map[string]any)
	body, _ := command["body"].(string)
	data, _ := command["data"].(map[string]any)
	metadata, _ := command["metadata"].(map[string]any)

	from, _ := raw["from"].(map[string]any)
	chatID, _ := getUUID(from, "group_chat_id")
	userHUID, _ := getUUID(from, "user_huid")

	var file *models.FileAttachment
	var location *models.Location
	var contact *models.Contact
	var link *models.Link
	var sticker *models.StickerAttachment

	if attachments, ok := raw["attachments"].([]any); ok && len(attachments) > 0 {
		if att, ok := attachments[0].(map[string]any); ok {
			parseIncomingAttachment(att, body, &file, &location, &contact, &link, &sticker)
		}
	}

	mentions := models.MentionList{}
	var forward *models.Forward
	var reply *models.Reply
	if entities, ok := raw["entities"].([]any); ok {
		for _, e := range entities {
			if ent, ok := e.(map[string]any); ok {
				parseEntity(ent, &mentions, &forward, &reply)
			}
		}
	}

	return models.IncomingMessage{
		Bot:      models.BotContext{ID: botID},
		Chat:     models.ChatContext{ID: chatID},
		Sender:   models.UserContext{HUID: userHUID},
		SyncID:   syncID,
		Body:     body,
		Data:     data,
		Meta:     metadata,
		Mentions: mentions,
		Forward:  forward,
		Reply:    reply,
		File:     file,
		Location: location,
		Contact:  contact,
		Link:     link,
		Sticker:  sticker,
		Raw:      raw,
	}, nil
}

// ParseSystemEvent parses raw BotX command payload into SystemEvent.
func ParseSystemEvent(raw map[string]any) (models.SystemEvent, error) {
	syncID, _ := getUUID(raw, "sync_id")
	botID, ok := getUUID(raw, "bot_id")
	if !ok {
		return models.SystemEvent{}, fmt.Errorf("%w: bot_id", ErrInvalidPayload)
	}

	command, _ := raw["command"].(map[string]any)
	body, _ := command["body"].(string)
	data, _ := command["data"].(map[string]any)

	from, _ := raw["from"].(map[string]any)
	chatID, _ := getUUID(from, "group_chat_id")
	userHUID, _ := getUUID(from, "user_huid")

	return models.SystemEvent{
		Type:   body,
		Bot:    models.BotContext{ID: botID},
		Chat:   models.ChatContext{ID: chatID},
		User:   models.UserContext{HUID: userHUID},
		SyncID: syncID,
		Data:   data,
		Raw:    raw,
	}, nil
}

// ParseSmartAppEvent parses sync smartapp event payload.
func ParseSmartAppEvent(raw map[string]any) (models.SmartAppEvent, error) {
	botID, ok := getUUID(raw, "bot_id")
	if !ok {
		return models.SmartAppEvent{}, fmt.Errorf("%w: bot_id", ErrInvalidPayload)
	}
	chatID, _ := getUUID(raw, "group_chat_id")

	senderInfo, _ := raw["sender_info"].(map[string]any)
	userHUID, _ := getUUID(senderInfo, "user_huid")

	data := map[string]any{}
	if m, ok := raw["method"].(string); ok {
		data["method"] = m
	}
	if p, ok := raw["payload"].(map[string]any); ok {
		data["payload"] = p
	}

	return models.SmartAppEvent{
		Bot:    models.BotContext{ID: botID},
		Chat:   models.ChatContext{ID: chatID},
		Sender: models.UserContext{HUID: userHUID},
		Data:   data,
		Raw:    raw,
	}, nil
}

// ParseStatusRecipient parses status query params.
func ParseStatusRecipient(q map[string]string) (models.StatusRecipient, error) {
	botID, err := uuid.Parse(q["bot_id"])
	if err != nil {
		return models.StatusRecipient{}, err
	}
	huid, err := uuid.Parse(q["user_huid"])
	if err != nil {
		return models.StatusRecipient{}, err
	}
	return models.StatusRecipient{BotID: botID, HUID: huid}, nil
}

// ParseCallback parses callback payload.
func ParseCallback(raw map[string]any) (models.MethodCallback, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return models.MethodCallback{}, err
	}
	var cb models.MethodCallback
	if err := json.Unmarshal(b, &cb); err != nil {
		return models.MethodCallback{}, err
	}
	if cb.SyncID == uuid.Nil {
		return models.MethodCallback{}, fmt.Errorf("%w: sync_id", ErrInvalidPayload)
	}
	return cb, nil
}

func getUUID(m map[string]any, key string) (uuid.UUID, bool) {
	v, ok := m[key]
	if !ok {
		return uuid.Nil, false
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func parseIncomingAttachment(att map[string]any, body string, file **models.FileAttachment, location **models.Location, contact **models.Contact, link **models.Link, sticker **models.StickerAttachment) {
	attType, _ := att["type"].(string)
	data, _ := att["data"].(map[string]any)

	switch attType {
	case string(models.AttachmentImage):
		content, _ := data["content"].(string)
		fileName, _ := data["file_name"].(string)
		fa, err := models.AttachmentFromAPI(models.AttachmentImage, fileName, content, 0)
		if err == nil {
			*file = &fa
		}
	case string(models.AttachmentVideo):
		content, _ := data["content"].(string)
		fileName, _ := data["file_name"].(string)
		duration, _ := toInt(data["duration"])
		fa, err := models.AttachmentFromAPI(models.AttachmentVideo, fileName, content, duration)
		if err == nil {
			*file = &fa
		}
	case string(models.AttachmentDocument):
		content, _ := data["content"].(string)
		fileName, _ := data["file_name"].(string)
		fa, err := models.AttachmentFromAPI(models.AttachmentDocument, fileName, content, 0)
		if err == nil {
			*file = &fa
		}
	case string(models.AttachmentVoice):
		content, _ := data["content"].(string)
		duration, _ := toInt(data["duration"])
		fa, err := models.AttachmentFromAPI(models.AttachmentVoice, "record", content, duration)
		if err == nil {
			*file = &fa
		}
	case string(models.AttachmentLocation):
		*location = &models.Location{
			Name:      toString(data["location_name"]),
			Address:   toString(data["location_address"]),
			Latitude:  toString(data["location_lat"]),
			Longitude: toString(data["location_lng"]),
		}
	case string(models.AttachmentContact):
		*contact = &models.Contact{Name: toString(data["contact_name"])}
	case string(models.AttachmentLink):
		*link = &models.Link{
			URL:     toString(data["url"]),
			Title:   toString(data["url_title"]),
			Preview: toString(data["url_preview"]),
			Text:    toString(data["url_text"]),
		}
	case string(models.AttachmentSticker):
		*sticker = &models.StickerAttachment{
			ID:       toString(data["id"]),
			ImageURL: toString(data["link"]),
			PackID:   toString(data["pack"]),
			Emoji:    body,
		}
	}
}

func parseEntity(ent map[string]any, mentions *models.MentionList, forward **models.Forward, reply **models.Reply) {
	entType, _ := ent["type"].(string)
	data, _ := ent["data"].(map[string]any)

	switch entType {
	case "mention":
		mentionType, _ := data["mention_type"].(string)
		mentionData, _ := data["mention_data"].(map[string]any)
		var idStr string
		var name string
		switch mentionType {
		case string(models.MentionUser), string(models.MentionContact):
			idStr = toString(mentionData["user_huid"])
			name = toString(mentionData["name"])
		case string(models.MentionChat), string(models.MentionChannel):
			idStr = toString(mentionData["group_chat_id"])
			name = toString(mentionData["name"])
		case string(models.MentionAll):
			idStr = ""
		}
		if mentionType == string(models.MentionAll) {
			*mentions = append(*mentions, models.Mention{Type: models.MentionType(mentionType), EntityID: uuid.Nil, Name: name})
			return
		}
		if id, err := uuid.Parse(idStr); err == nil {
			*mentions = append(*mentions, models.Mention{Type: models.MentionType(mentionType), EntityID: id, Name: name})
		}
	case "forward":
		if *forward != nil {
			return
		}
		*forward = &models.Forward{
			ChatID:     uuidOrNil(data["group_chat_id"]),
			AuthorID:   uuidOrNil(data["sender_huid"]),
			SyncID:     uuidOrNil(data["source_sync_id"]),
			ChatName:   toString(data["source_chat_name"]),
			Type:       toString(data["forward_type"]),
			InsertedAt: toString(data["source_inserted_at"]),
		}
	case "reply":
		if *reply != nil {
			return
		}
		*reply = &models.Reply{
			AuthorID: uuidOrNil(data["sender"]),
			SyncID:   uuidOrNil(data["source_sync_id"]),
			Body:     toString(data["body"]),
		}
	}
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toInt(v any) (int, bool) {
	switch t := v.(type) {
	case int:
		return t, true
	case int64:
		return int(t), true
	case float64:
		return int(t), true
	default:
		return 0, false
	}
}

func uuidOrNil(v any) uuid.UUID {
	s := toString(v)
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}
