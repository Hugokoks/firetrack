package notes

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepositry(db *sql.DB) *Repository {

	return &Repository{db: db}

}

func (r *Repository) Create(jobID, authorID, content string) (*Note, error) {

	note := &Note{

		ID:        uuid.NewString(),
		JobID:     jobID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO job_notes (id,job_id,author_id,content,created_at)
		VALUES($1,$2,$3,$4,$5)
	`
	_, err := r.db.Exec(query, note.ID, note.JobID, note.AuthorID, note.Content, note.CreatedAt)

	if err != nil {
		return nil, err
	}

	return note, nil

}

func (r *Repository) GetByID(noteID string) (*Note, error) {
	query := `
		SELECT
			id::text,
			job_id::text,
			author_id::text,
			content,
			created_at
		FROM job_notes
		WHERE id::text = $1
		LIMIT 1
	`

	var note Note

	err := r.db.QueryRow(query, noteID).Scan(
		&note.ID,
		&note.JobID,
		&note.AuthorID,
		&note.Content,
		&note.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &note, nil
}

func (r *Repository) DeleteByID(noteID string) error {
	query := `
		DELETE FROM job_notes
		WHERE id::text = $1
	`

	_, err := r.db.Exec(query, noteID)
	return err
}

func (r *Repository) GetByJobID(jobID string) ([]Note, error) {
	query := `
		SELECT
			id::text,
			job_id::text,
			author_id::text,
			content,
			created_at
		FROM job_notes
		WHERE job_id::text = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var note Note

		err := rows.Scan(
			&note.ID,
			&note.JobID,
			&note.AuthorID,
			&note.Content,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
