package files

import (
	"database/sql"
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
