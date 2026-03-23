package notes

import (
	"net/http"

	"firetrack/internal/activity"
	"firetrack/internal/httputil"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	jobID, ok := httputil.RequireParam(c, "id", "missing job id")
	if !ok {
		return
	}

	var input NoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httputil.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, ok := httputil.GetCurrentUser(c)
	if !ok {
		return
	}

	note, err := h.service.Create(jobID, user.ID, input.Content)
	if err != nil {
		if err == ErrEmptyContent {
			httputil.RespondError(c, http.StatusBadRequest, "content cannot be empty")
			return
		}

		httputil.RespondError(c, http.StatusInternalServerError, "failed to create note")
		return
	}

	activity.Set(c, activity.Payload{
		JobID:       jobID,
		UserID:      user.ID,
		ActionType:  "note_created",
		ActionLabel: "Note created",
		Meta: gin.H{
			"note_id": note.ID,
		},
	})

	c.JSON(http.StatusCreated, gin.H{
		"note": note,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	note, ok := getContextNote(c)
	if !ok {
		return
	}

	user, ok := httputil.GetCurrentUser(c)
	if !ok {
		return
	}

	if err := h.service.Delete(note.ID); err != nil {
		httputil.RespondError(c, http.StatusInternalServerError, "failed to delete note")
		return
	}

	activity.Set(c, activity.Payload{
		JobID:       note.JobID,
		UserID:      user.ID,
		ActionType:  "note_deleted",
		ActionLabel: "Note deleted",
		Meta: gin.H{
			"note_id": note.ID,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "note deleted successfully",
	})
}

func (h *Handler) GetByJobID(c *gin.Context) {
	jobID, ok := httputil.RequireParam(c, "id", "missing job id")
	if !ok {
		return
	}

	notes, err := h.service.GetByJobID(jobID)
	if err != nil {
		httputil.RespondError(c, http.StatusInternalServerError, "failed to fetch notes")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

func (h *Handler) Update(c *gin.Context) {
	note, ok := getContextNote(c)
	if !ok {
		return
	}

	user, ok := httputil.GetCurrentUser(c)
	if !ok {
		return
	}

	var input NoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httputil.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedNote, err := h.service.Update(note.ID, input.Content)
	if err != nil {
		if err == ErrEmptyContent {
			httputil.RespondError(c, http.StatusBadRequest, "content cannot be empty")
			return
		}

		if err == ErrNoteNotFound {
			httputil.RespondError(c, http.StatusNotFound, "note not found")
			return
		}

		httputil.RespondError(c, http.StatusInternalServerError, "failed to update note")
		return
	}

	activity.Set(c, activity.Payload{
		JobID:       updatedNote.JobID,
		UserID:      user.ID,
		ActionType:  "note_updated",
		ActionLabel: "Note updated",
		Meta: gin.H{
			"note_id": updatedNote.ID,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"note": updatedNote,
	})
}
