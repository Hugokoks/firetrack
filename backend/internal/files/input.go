package files

import "mime/multipart"

type CreateFileInput struct {
	JobID      string
	UploadedBy string
	FileHeader *multipart.FileHeader
}
