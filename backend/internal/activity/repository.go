package activity

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(payload Payload) error {
	metaJSON, err := json.Marshal(payload.Meta)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO job_activity (
			id,
			job_id,
			user_id,
			action_type,
			action_label,
			meta,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(
		query,
		uuid.NewString(),
		payload.JobID,
		payload.UserID,
		payload.ActionType,
		payload.ActionLabel,
		metaJSON,
		time.Now(),
	)

	return err
}
