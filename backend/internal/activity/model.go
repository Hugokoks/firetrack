package activity

type Payload struct {
	JobID       string
	UserID      string
	ActionType  string
	ActionLabel string
	Meta        any
}
