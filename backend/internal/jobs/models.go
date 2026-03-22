package jobs

import "time"

type Job struct {
	ID             string
	JobNumber      *string
	Title          string
	CustomerName   *string
	Address        string
	City           *string
	Zip            *string
	Country        string
	Latitude       *float64
	Longitude      *float64
	ScheduledStart time.Time
	ScheduledEnd   *time.Time
	CompletedAt    *time.Time
	Status         string
	Priority       string
	AssignedUserID *string
	CreatedBy      string
	Description    *string
	GoogleEventID  *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
