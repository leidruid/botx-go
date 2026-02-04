package models

import (
	"encoding/base64"
	"errors"
	"strings"
)

// FileAttachment is a decoded incoming file attachment.
type FileAttachment struct {
	Type     AttachmentType
	Filename string
	Size     int
	Content  []byte
	Duration int
}

// Location attachment.
type Location struct {
	Name      string
	Address   string
	Latitude  string
	Longitude string
}

// Contact attachment.
type Contact struct {
	Name string
}

// Link attachment.
type Link struct {
	URL     string
	Title   string
	Preview string
	Text    string
}

// StickerAttachment represents sticker attached to message.
type StickerAttachment struct {
	ID       string
	ImageURL string
	PackID   string
	Emoji    string
}

// OutgoingAttachment represents a file to send.
type OutgoingAttachment struct {
	Filename string
	Content  []byte
}

// APINotificationAttachment represents BotX file attachment payload.
type APINotificationAttachment struct {
	FileName string `json:"file_name"`
	Data     string `json:"data"`
}

// AttachmentFromAPI decodes BotX attachment (rfc2397 data URL) into FileAttachment.
func AttachmentFromAPI(typ AttachmentType, fileName string, data string, duration int) (FileAttachment, error) {
	content, err := DecodeRFC2397(data)
	if err != nil {
		return FileAttachment{}, err
	}
	return FileAttachment{
		Type:     typ,
		Filename: fileName,
		Size:     len(content),
		Content:  content,
		Duration: duration,
	}, nil
}

// AttachmentToAPI encodes a file attachment for BotX request.
func AttachmentToAPI(att OutgoingAttachment) (APINotificationAttachment, error) {
	if att.Filename == "" {
		return APINotificationAttachment{}, errors.New("filename is required")
	}
	mimetype := MimeTypeByFilename(att.Filename)
	encoded := EncodeRFC2397(att.Content, mimetype)
	return APINotificationAttachment{FileName: att.Filename, Data: encoded}, nil
}

// EncodeRFC2397 builds data URL for bytes.
func EncodeRFC2397(content []byte, mimetype string) string {
	b64 := base64.StdEncoding.EncodeToString(content)
	return "data:" + mimetype + ";base64," + b64
}

// DecodeRFC2397 decodes data URL into bytes.
func DecodeRFC2397(encoded string) ([]byte, error) {
	if encoded == "" {
		return nil, nil
	}
	parts := strings.SplitN(encoded, ",", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid data url")
	}
	return base64.StdEncoding.DecodeString(parts[1])
}

// MimeTypeByFilename returns mimetype by extension.
func MimeTypeByFilename(filename string) string {
	ext := strings.ToLower(strings.TrimPrefix(strings.ToLower(extFromFilename(filename)), "."))
	if ext == "" {
		return defaultMimeType
	}
	if mt, ok := extToMime[ext]; ok {
		return mt
	}
	return defaultMimeType
}

func extFromFilename(filename string) string {
	i := strings.LastIndex(filename, ".")
	if i < 0 || i == len(filename)-1 {
		return ""
	}
	return filename[i+1:]
}

const defaultMimeType = "application/octet-stream"

var extToMime = map[string]string{
	"png":  "image/png",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"gif":  "image/gif",
	"webp": "image/webp",
	"mp4":  "video/mp4",
	"mp3":  "audio/mpeg",
	"wav":  "audio/wav",
	"pdf":  "application/pdf",
	"txt":  "text/plain",
	"json": "application/json",
}
