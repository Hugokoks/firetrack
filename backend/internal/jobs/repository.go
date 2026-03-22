package jobs

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(input CreateJobInput, createdBy string) (*Job, error) {
	now := time.Now()

	job := &Job{
		ID:             uuid.NewString(),
		Title:          input.Title,
		CustomerName:   input.CustomerName,
		Address:        input.Address,
		City:           input.City,
		Zip:            input.Zip,
		Country:        input.Country,
		Latitude:       input.Latitude,
		Longitude:      input.Longitude,
		ScheduledStart: input.ScheduledStart,
		ScheduledEnd:   input.ScheduledEnd,
		Status:         "planned",
		Priority:       input.Priority,
		AssignedUserID: input.AssignedUserID,
		CreatedBy:      createdBy,
		Description:    input.Description,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	query := `
		INSERT INTO jobs (
			id,
			title,
			customer_name,
			address,
			city,
			zip,
			country,
			latitude,
			longitude,
			scheduled_start,
			scheduled_end,
			status,
			priority,
			assigned_user_id,
			created_by,
			description,
			created_at,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16, $17, $18
		)
	`

	_, err := r.db.Exec(
		query,
		job.ID,
		job.Title,
		job.CustomerName,
		job.Address,
		job.City,
		job.Zip,
		job.Country,
		job.Latitude,
		job.Longitude,
		job.ScheduledStart,
		job.ScheduledEnd,
		job.Status,
		job.Priority,
		job.AssignedUserID,
		job.CreatedBy,
		job.Description,
		job.CreatedAt,
		job.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (r *Repository) GetAll() ([]Job, error) {
	query := `
		SELECT
			id::text,
			job_number,
			title,
			customer_name,
			address,
			city,
			zip,
			country,
			latitude,
			longitude,
			scheduled_start,
			scheduled_end,
			completed_at,
			status,
			priority,
			assigned_user_id::text,
			created_by::text,
			description,
			google_event_id,
			created_at,
			updated_at
		FROM jobs
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job

		err := rows.Scan(
			&job.ID,
			&job.JobNumber,
			&job.Title,
			&job.CustomerName,
			&job.Address,
			&job.City,
			&job.Zip,
			&job.Country,
			&job.Latitude,
			&job.Longitude,
			&job.ScheduledStart,
			&job.ScheduledEnd,
			&job.CompletedAt,
			&job.Status,
			&job.Priority,
			&job.AssignedUserID,
			&job.CreatedBy,
			&job.Description,
			&job.GoogleEventID,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *Repository) GetByID(id string) (*Job, error) {
	query := `
		SELECT
			id::text,
			job_number,
			title,
			customer_name,
			address,
			city,
			zip,
			country,
			latitude,
			longitude,
			scheduled_start,
			scheduled_end,
			completed_at,
			status,
			priority,
			assigned_user_id::text,
			created_by::text,
			description,
			google_event_id,
			created_at,
			updated_at
		FROM jobs
		WHERE id::text = $1
		LIMIT 1
	`

	var job Job

	err := r.db.QueryRow(query, id).Scan(
		&job.ID,
		&job.JobNumber,
		&job.Title,
		&job.CustomerName,
		&job.Address,
		&job.City,
		&job.Zip,
		&job.Country,
		&job.Latitude,
		&job.Longitude,
		&job.ScheduledStart,
		&job.ScheduledEnd,
		&job.CompletedAt,
		&job.Status,
		&job.Priority,
		&job.AssignedUserID,
		&job.CreatedBy,
		&job.Description,
		&job.GoogleEventID,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &job, nil
}
