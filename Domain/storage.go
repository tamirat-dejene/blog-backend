package domain

import (
	"context"
)

type StorageService interface {
	UploadFile(ctx context.Context, fileName string, fileData []byte) (string, error)
}
