package notes

import (
	"firetrack/internal/auth"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {

	return &Handler{service: service}

}

func (h *Handler) Create(c *gin.Context) {

	jobID := c.Param("id")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing job id",
		})
		return
	}

	var input CreateNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	rawUser, _ := c.Get(auth.ContextUserKey)
	user := rawUser.(*auth.User)

	note, err := h.service.Create(jobID, user.ID, input.Content)
	fmt.Println(err)
	if err != nil {
		if err == ErrEmptyContent {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "content cannot be empty",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create note",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"note": note,
	})

}

func (h *Handler) Delete(c *gin.Context) {
	rawNote, exists := c.Get(ContextNoteKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "note not found in context",
		})
		return
	}

	note, ok := rawNote.(*Note)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid note context",
		})
		return
	}

	if err := h.service.Delete(note.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete note",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "note deleted successfully",
	})
}

func (h *Handler) GetByJobID(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing job id",
		})
		return
	}

	notes, err := h.service.GetByJobID(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch notes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}
