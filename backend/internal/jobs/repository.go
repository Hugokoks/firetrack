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

func (r *Repository) Update(job *Job) (*Job, error) {
	job.UpdatedAt = time.Now()

	query := `
		UPDATE jobs
		SET
			title = $1,
			customer_name = $2,
			address = $3,
			city = $4,
			zip = $5,
			country = $6,
			latitude = $7,
			longitude = $8,
			scheduled_start = $9,
			scheduled_end = $10,
			completed_at = $11,
			status = $12,
			priority = $13,
			assigned_user_id = $14,
			description = $15,
			updated_at = $16
		WHERE id::text = $17
		RETURNING
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
	`

	var updatedJob Job

	err := r.db.QueryRow(
		query,
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
		job.CompletedAt,
		job.Status,
		job.Priority,
		job.AssignedUserID,
		job.Description,
		job.UpdatedAt,
		job.ID,
	).Scan(
		&updatedJob.ID,
		&updatedJob.JobNumber,
		&updatedJob.Title,
		&updatedJob.CustomerName,
		&updatedJob.Address,
		&updatedJob.City,
		&updatedJob.Zip,
		&updatedJob.Country,
		&updatedJob.Latitude,
		&updatedJob.Longitude,
		&updatedJob.ScheduledStart,
		&updatedJob.ScheduledEnd,
		&updatedJob.CompletedAt,
		&updatedJob.Status,
		&updatedJob.Priority,
		&updatedJob.AssignedUserID,
		&updatedJob.CreatedBy,
		&updatedJob.Description,
		&updatedJob.GoogleEventID,
		&updatedJob.CreatedAt,
		&updatedJob.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &updatedJob, nil
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM job_files WHERE id = $1`, id)
	return err
}