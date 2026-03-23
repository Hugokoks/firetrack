package files

import (
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

var allowedMimeTypes = map[string]string{
	"application/pdf": ".pdf",
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
}

func detectAllowedMime(file multipartFile) (mimeType string, ext string, err error) {
	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		return "", "", err
	}

	mimeType = mtype.String()

	allowedExt, ok := allowedMimeTypes[mimeType]
	if !ok {
		return "", "", ErrFileTypeNotAllowed
	}

	return mimeType, allowedExt, nil
}

func sanitizeFileName(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, "\x00", "")
	return name
}
