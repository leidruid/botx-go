package bot

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

// UploadFile uploads a file and returns async file metadata.
func (b *Bot) UploadFile(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, filename string, r io.Reader, duration optional.Optional[int], caption optional.Optional[string]) (models.AsyncFile, error) {
	return b.Client.UploadFile(ctx, botID, chatID, filename, r, duration, caption)
}

// DownloadFile downloads a file into writer.
func (b *Bot) DownloadFile(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, fileID uuid.UUID, w io.Writer) error {
	return b.Client.DownloadFile(ctx, botID, chatID, fileID, w)
}
