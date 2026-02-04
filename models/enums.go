package models

// AttachmentType represents attachment kind.
type AttachmentType string

const (
	AttachmentImage    AttachmentType = "image"
	AttachmentVideo    AttachmentType = "video"
	AttachmentDocument AttachmentType = "document"
	AttachmentVoice    AttachmentType = "voice"
	AttachmentLocation AttachmentType = "location"
	AttachmentContact  AttachmentType = "contact"
	AttachmentLink     AttachmentType = "link"
	AttachmentSticker  AttachmentType = "sticker"
)

// MentionType represents mention kind.
type MentionType string

const (
	MentionUser    MentionType = "user"
	MentionContact MentionType = "contact"
	MentionChat    MentionType = "chat"
	MentionChannel MentionType = "channel"
	MentionAll     MentionType = "all"
)
