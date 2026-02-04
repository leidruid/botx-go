package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

type uploadFileResponse struct {
	Status string `json:"status"`
	Result struct {
		Type     string    `json:"type"`
		File     string    `json:"file"`
		FileMime string    `json:"file_mime_type"`
		FileID   uuid.UUID `json:"file_id"`
		FileName string    `json:"file_name"`
		FileSize int       `json:"file_size"`
		FileHash string    `json:"file_hash"`
		Duration int       `json:"duration"`
	} `json:"result"`
}

// UploadFile uploads file to chat and returns async file metadata.
func (c *Client) UploadFile(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, filename string, r io.Reader, duration optional.Optional[int], caption optional.Optional[string]) (models.AsyncFile, error) {
	path := "/api/v3/botx/files/upload"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return models.AsyncFile{}, err
	}

	meta := map[string]any{}
	if duration.Set {
		meta["duration"] = duration.Value
	}
	if caption.Set {
		meta["caption"] = caption.Value
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return models.AsyncFile{}, err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("group_chat_id", chatID.String())
	_ = writer.WriteField("meta", string(metaJSON))

	filePart, err := writer.CreateFormFile("content", filename)
	if err != nil {
		return models.AsyncFile{}, err
	}
	if _, err := io.Copy(filePart, r); err != nil {
		return models.AsyncFile{}, err
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, urlStr, &body)
	if err != nil {
		return models.AsyncFile{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var resp uploadFileResponse
	_, err = c.doAuthorizedJSON(ctx, botID, req, &resp)
	if err != nil {
		return models.AsyncFile{}, err
	}
	if resp.Status != "ok" {
		return models.AsyncFile{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return models.AsyncFile{
		Type:     models.AttachmentType(resp.Result.Type),
		FileID:   resp.Result.FileID,
		FileURL:  resp.Result.File,
		FileName: resp.Result.FileName,
		FileSize: resp.Result.FileSize,
		FileHash: resp.Result.FileHash,
		MimeType: resp.Result.FileMime,
		Duration: resp.Result.Duration,
	}, nil
}

// DownloadFile downloads a file into writer.
func (c *Client) DownloadFile(ctx context.Context, botID uuid.UUID, chatID uuid.UUID, fileID uuid.UUID, w io.Writer) error {
	path := "/api/v3/botx/files/download"
	urlStr, err := c.buildURL(botID, path)
	if err != nil {
		return err
	}
	params := url.Values{
		"group_chat_id": []string{chatID.String()},
		"file_id":       []string{fileID.String()},
		"is_preview":    []string{"false"},
	}
	urlStr = urlStr + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	resp, err := c.doAuthorizedJSON(ctx, botID, req, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}
