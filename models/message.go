package models

import "github.com/google/uuid"

// IncomingMessage represents a parsed incoming bot command message.
type IncomingMessage struct {
	Bot      BotContext
	Chat     ChatContext
	Sender   UserContext
	SyncID   uuid.UUID
	Body     string
	Data     map[string]any
	Meta     map[string]any
	Mentions MentionList
	Forward  *Forward
	Reply    *Reply
	File     *FileAttachment
	Location *Location
	Contact  *Contact
	Link     *Link
	Sticker  *StickerAttachment
	Raw      map[string]any
}

// Argument returns message text without command prefix.
func (m IncomingMessage) Argument() string {
	if m.Body == "" {
		return ""
	}
	parts := splitWhitespace(m.Body)
	if len(parts) == 0 {
		return ""
	}
	cmdLen := len(parts[0])
	if len(m.Body) <= cmdLen {
		return ""
	}
	return trimWhitespace(m.Body[cmdLen:])
}

// Arguments returns arguments split by whitespace.
func (m IncomingMessage) Arguments() []string {
	arg := m.Argument()
	if arg == "" {
		return nil
	}
	return splitWhitespace(arg)
}

func splitWhitespace(s string) []string {
	out := []string{}
	start := -1
	for i, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if start >= 0 {
				out = append(out, s[start:i])
				start = -1
			}
			continue
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		out = append(out, s[start:])
	}
	return out
}

func trimWhitespace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
