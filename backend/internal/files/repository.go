package files

import (
	"database/sql"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(file *File) (*File, error) {
	file.ID = uuid.NewString()
	file.CreatedAt = time.Now()

	query := `
		INSERT INTO job_files (
			id,
			job_id,
			uploaded_by,
			file_name,
			stored_name,
			file_path,
			mime_type,
			file_size,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(
		query,
		file.ID,
		file.JobID,
		file.UploadedBy,
		file.FileName,
		file.StoredName,
		file.FilePath,
		file.MimeType,
		file.FileSize,
		file.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r *Repository) GetByJobID(jobID string) ([]File, error) {
	query := `
		SELECT
			id,
			job_id,
			uploaded_by,
			file_name,
			stored_name,
			file_path,
			mime_type,
			file_size,
			created_at
		FROM job_files
		WHERE job_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File

	for rows.Next() {
		var file File

		err := rows.Scan(
			&file.ID,
			&file.JobID,
			&file.UploadedBy,
			&file.FileName,
			&file.StoredName,
			&file.FilePath,
			&file.MimeType,
			&file.FileSize,
			&file.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *Repository) GetByID(id string) (*File, error) {
	query := `
		SELECT
			id,
			job_id,
			uploaded_by,
			file_name,
			stored_name,
			file_path,
			mime_type,
			file_size,
			created_at
		FROM job_files
		WHERE id = $1
		LIMIT 1
	`

	var file File

	err := r.db.QueryRow(query, id).Scan(
		&file.ID,
		&file.JobID,
		&file.UploadedBy,
		&file.FileName,
		&file.StoredName,
		&file.FilePath,
		&file.MimeType,
		&file.FileSize,
		&file.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &file, nil
}

func (s *Service) Download(jobID, fileID string) (*File, io.ReadCloser, error) {
	job, err := s.jobsRepo.GetByID(jobID)
	if err != nil {
		return nil, nil, err
	}
	if job == nil {
		return nil, nil, ErrJobNotFound
	}

	file, err := s.repo.GetByID(fileID)
	if err != nil {
		return nil, nil, err
	}
	if file == nil || file.JobID != jobID {
		return nil, nil, ErrFileNotFound
	}

	reader, err := s.storage.Open(file.FilePath)
	if err != nil {
		return nil, nil, err
	}

	return file, reader, nil
}

func (r * Repository) Delete (id string) error{

	_,err := r.db.Exec(`DELETE FROM job_files WHERE id = $1`,id)
	return  err
}