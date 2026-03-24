package files

import (
	"io"

	"firetrack/internal/jobs"

	"github.com/google/uuid"
)

type multipartFile interface {
	io.Reader
	io.Seeker
	io.ReaderAt
}

type Service struct {
	repo     *Repository
	storage  *Storage
	jobsRepo *jobs.Repository
	maxSize  int64
}

func NewService(repo *Repository, storage *Storage, jobsRepo *jobs.Repository, maxSize int64) *Service {
	return &Service{
		repo:     repo,
		storage:  storage,
		jobsRepo: jobsRepo,
		maxSize:  maxSize,
	}
}

func (s *Service) Create(input CreateFileInput) (*File, error) {
	if input.FileHeader == nil {
		return nil, ErrMissingFile
	}

	job, err := s.jobsRepo.GetByID(input.JobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrJobNotFound
	}

	if input.FileHeader.Size > s.maxSize {
		return nil, ErrFileTooLarge
	}

	src, err := input.FileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	file, ok := src.(multipartFile)
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}

	mimeType, ext, err := detectAllowedMime(file)
	if err != nil {
		return nil, err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	storedName := uuid.NewString() + ext
	relPath := s.storage.BuildRelativePath(input.JobID, storedName)

	if err := s.storage.EnsureDirFor(relPath); err != nil {
		return nil, err
	}

	written, err := s.storage.Save(relPath, file)
	if err != nil {
		return nil, err
	}

	record := &File{
		JobID:      input.JobID,
		UploadedBy: input.UploadedBy,
		FileName:   sanitizeFileName(input.FileHeader.Filename),
		StoredName: storedName,
		FilePath:   relPath,
		MimeType:   mimeType,
		FileSize:   written,
	}

	created, err := s.repo.Create(record)
	if err != nil {
		_ = s.storage.Remove(relPath)
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByJobID(jobID string) ([]File, error) {
	job, err := s.jobsRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrJobNotFound
	}

	return s.repo.GetByJobID(jobID)
}
