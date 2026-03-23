package files

import "errors"

var (
	ErrMissingFile        = errors.New("missing file")
	ErrFileTooLarge       = errors.New("file too large")
	ErrFileTypeNotAllowed = errors.New("file type not allowed")
	ErrJobNotFound        = errors.New("job not found")
)
