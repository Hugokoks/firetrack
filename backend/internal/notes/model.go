package notes

import "time"

type Note struct {
	ID        string
	JobID     string
	AuthorID  string
	Content   string
	CreatedAt time.Time
}
