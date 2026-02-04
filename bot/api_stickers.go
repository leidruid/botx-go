package bot

import (
	"context"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

// CreateStickerPack creates sticker pack.
func (b *Bot) CreateStickerPack(ctx context.Context, botID uuid.UUID, name string, userHUID optional.Optional[uuid.UUID]) (models.StickerPack, error) {
	return b.Client.CreateStickerPack(ctx, botID, name, userHUID)
}

// EditStickerPack edits sticker pack.
func (b *Bot) EditStickerPack(ctx context.Context, botID uuid.UUID, packID uuid.UUID, name string, preview uuid.UUID, stickersOrder []uuid.UUID) (models.StickerPack, error) {
	return b.Client.EditStickerPack(ctx, botID, packID, name, preview, stickersOrder)
}

// DeleteStickerPack deletes sticker pack.
func (b *Bot) DeleteStickerPack(ctx context.Context, botID uuid.UUID, packID uuid.UUID) error {
	return b.Client.DeleteStickerPack(ctx, botID, packID)
}

// GetStickerPack returns sticker pack by id.
func (b *Bot) GetStickerPack(ctx context.Context, botID uuid.UUID, packID uuid.UUID) (models.StickerPack, error) {
	return b.Client.GetStickerPack(ctx, botID, packID)
}

// GetStickerPacks returns sticker packs list.
func (b *Bot) GetStickerPacks(ctx context.Context, botID uuid.UUID, userHUID uuid.UUID, limit int, after optional.Optional[string]) (models.StickerPackPage, error) {
	return b.Client.GetStickerPacks(ctx, botID, userHUID, limit, after)
}

// GetSticker returns sticker by pack and sticker id.
func (b *Bot) GetSticker(ctx context.Context, botID uuid.UUID, packID uuid.UUID, stickerID uuid.UUID) (models.Sticker, error) {
	return b.Client.GetSticker(ctx, botID, packID, stickerID)
}

// AddSticker adds sticker to pack.
func (b *Bot) AddSticker(ctx context.Context, botID uuid.UUID, packID uuid.UUID, emoji string, png []byte) (models.Sticker, error) {
	return b.Client.AddSticker(ctx, botID, packID, emoji, png)
}

// DeleteSticker deletes sticker from pack.
func (b *Bot) DeleteSticker(ctx context.Context, botID uuid.UUID, packID uuid.UUID, stickerID uuid.UUID) error {
	return b.Client.DeleteSticker(ctx, botID, packID, stickerID)
}
