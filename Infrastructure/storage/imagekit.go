package storage

import (
	"bytes"
	"context"

	// "embed"
	"fmt"
	domain "g6/blog-api/Domain"
	"image"
	_ "image/gif"  // Register GIF decoder (optional, for future)
	_ "image/jpeg" // Register JPEG decoder (optional, for future)
	_ "image/png"  // Register PNG decoder

	"path/filepath"
	"strings"

	_ "embed"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

type ImageKitStorage struct {
	client *imagekit.ImageKit
}

func NewImageKitStorage(privateKey, publicKey, urlEndpoint string) domain.StorageService {

	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		UrlEndpoint: urlEndpoint,
	})

	return &ImageKitStorage{client: ik}
}

func (s *ImageKitStorage) UploadFile(ctx context.Context, fileName string, fileData []byte) (string, error) {
	// Upload the image to ImageKit
	if fileName == "" || len(fileData) == 0 {
		return "", domain.ErrInvalidFile
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", domain.ErrInvalidFile
	}

	// Validate image data
	_, format, err := image.DecodeConfig(bytes.NewReader(fileData))
	if err != nil {
		return "", domain.ErrInvalidFile
	}
	if format != "jpeg" && format != "png" && format != "gif" {
		return "", fmt.Errorf("unsupported image format: %s", format)
	}
	// Use the fileData passed to the function
	file := bytes.NewReader(fileData)

	isPrivateFile := false
	useUniqueFileName := true
	resp, err := s.client.Uploader.Upload(ctx, file, uploader.UploadParam{
		FileName:          fileName,
		IsPrivateFile:     &isPrivateFile,
		UseUniqueFileName: &useUniqueFileName,
		Tags:              "image",
		Folder:            "/blog-profile/",
	})
	if err != nil {
		return "", fmt.Errorf("imagekit upload failed: %w", err)
	}

	return resp.Data.Url, nil
}
