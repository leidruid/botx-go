package models

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

// Mention describes a mention in message.
type Mention struct {
	Type     MentionType
	EntityID uuid.UUID
	Name     string
}

// MentionList stores mentions.
type MentionList []Mention

// BuildEmbedMention creates <embed_mention> markup for message body.
func BuildEmbedMention(typ MentionType, entityID uuid.UUID, name string) string {
	return fmt.Sprintf("<embed_mention>%s:%s:%s</embed_mention>", typ, entityID.String(), name)
}

// APIMention represents BotX mention payload.
type APIMention struct {
	MentionType string         `json:"mention_type"`
	MentionID   uuid.UUID      `json:"mention_id"`
	MentionData map[string]any `json:"mention_data,omitempty"`
}

func buildAPIMention(t MentionType, entityID string, name string) APIMention {
	mentionID := uuid.New()
	switch t {
	case MentionUser, MentionContact:
		data := map[string]any{"user_huid": entityID}
		if name != "" {
			data["name"] = name
		}
		return APIMention{MentionType: string(t), MentionID: mentionID, MentionData: data}
	case MentionChat, MentionChannel:
		data := map[string]any{"group_chat_id": entityID}
		if name != "" {
			data["name"] = name
		}
		return APIMention{MentionType: string(t), MentionID: mentionID, MentionData: data}
	case MentionAll:
		return APIMention{MentionType: string(t), MentionID: mentionID}
	default:
		return APIMention{MentionType: string(t), MentionID: mentionID}
	}
}

var embedMentionRe = regexp.MustCompile(`(<embed_mention>(?P<mention_type>.+?):(?P<mentioned_entity_id>[0-9a-f-]*?):(?P<mention_name>.*?)</embed_mention>)`)

// FindAndReplaceEmbedMentions replaces embed mentions with BotX inline format and builds payload mentions.
func FindAndReplaceEmbedMentions(body string) (string, []APIMention) {
	mentions := []APIMention{}
	matches := embedMentionRe.FindAllStringSubmatchIndex(body, -1)
	if len(matches) == 0 {
		return body, mentions
	}

	// Build result incrementally.
	out := ""
	last := 0
	for _, idx := range matches {
		fullStart, fullEnd := idx[0], idx[1]
		sub := body[fullStart:fullEnd]
		group := embedMentionRe.FindStringSubmatch(sub)
		if len(group) < 5 {
			continue
		}
		mentionType := MentionType(group[2])
		entityID := group[3]
		name := group[4]

		mention := buildAPIMention(mentionType, entityID, name)
		mentions = append(mentions, mention)

		out += body[last:fullStart]
		out += mentionToInlineFormat(mentionType, mention.MentionID)
		last = fullEnd
	}
	out += body[last:]

	return out, mentions
}

func mentionToInlineFormat(t MentionType, mentionID uuid.UUID) string {
	switch t {
	case MentionContact:
		return fmt.Sprintf("@@{mention:%s}", mentionID.String())
	case MentionChat, MentionChannel:
		return fmt.Sprintf("##{mention:%s}", mentionID.String())
	default:
		return fmt.Sprintf("@{mention:%s}", mentionID.String())
	}
}
