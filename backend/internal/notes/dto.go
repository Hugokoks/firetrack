package notes

type NoteInput struct {
	Content string `json:"content" binding:"required"`
}
