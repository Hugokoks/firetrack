package files

import "time"

type File struct {
	ID         string
	JobID      string
	UploadedBy string
	FileName   string
	StoredName string
	FilePath   string
	MimeType   string
	FileSize   int64
	CreatedAt  time.Time
}
