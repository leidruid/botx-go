package models

import "github.com/google/uuid"

// Sticker represents sticker data.
type Sticker struct {
	ID     uuid.UUID
	Emoji  string
	Link   string
	PackID uuid.UUID
}

// StickerPack represents full sticker pack.
type StickerPack struct {
	ID       uuid.UUID
	Name     string
	IsPublic bool
	Stickers []Sticker
}

// StickerPackFromList represents sticker pack in list.
type StickerPackFromList struct {
	ID            uuid.UUID
	Name          string
	IsPublic      bool
	StickersCount int
	StickerIDs    []uuid.UUID
}

// StickerPackPage represents paginated list.
type StickerPackPage struct {
	StickerPacks []StickerPackFromList
	After        string
}
