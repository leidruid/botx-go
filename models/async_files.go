package models

import "github.com/google/uuid"

// AsyncFile represents uploaded file metadata for later download.
type AsyncFile struct {
	Type     AttachmentType
	FileID   uuid.UUID
	FileURL  string
	FileName string
	FileSize int
	FileHash string
	MimeType string
	Duration int
}
