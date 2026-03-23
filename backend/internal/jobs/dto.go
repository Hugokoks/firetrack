package jobs

import "time"

type CreateJobInput struct {
	Title          string     `json:"title" binding:"required"`
	CustomerName   *string    `json:"customer_name"`
	Address        string     `json:"address" binding:"required"`
	City           *string    `json:"city"`
	Zip            *string    `json:"zip"`
	Country        string     `json:"country"`
	Latitude       *float64   `json:"latitude"`
	Longitude      *float64   `json:"longitude"`
	ScheduledStart time.Time  `json:"scheduled_start" binding:"required"`
	ScheduledEnd   *time.Time `json:"scheduled_end"`
	AssignedUserID *string    `json:"assigned_user_id"`
	Description    *string    `json:"description"`
	Priority       string     `json:"priority"`
}

type UpdateJobInput struct {
	Title          *string    `json:"title"`
	CustomerName   *string    `json:"customer_name"`
	Address        *string    `json:"address"`
	City           *string    `json:"city"`
	Zip            *string    `json:"zip"`
	Country        *string    `json:"country"`
	Latitude       *float64   `json:"latitude"`
	Longitude      *float64   `json:"longitude"`
	ScheduledStart *time.Time `json:"scheduled_start"`
	ScheduledEnd   *time.Time `json:"scheduled_end"`
	CompletedAt    *time.Time `json:"completed_at"`
	Status         *string    `json:"status"`
	Priority       *string    `json:"priority"`
	AssignedUserID *string    `json:"assigned_user_id"`
	Description    *string    `json:"description"`
}
