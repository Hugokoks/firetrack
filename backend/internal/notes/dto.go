package notes

type CreateNoteInput struct {
	Content string `json:"content" binding:"required"`
}
