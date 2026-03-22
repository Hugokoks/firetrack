package notes

import (
	"net/http"

	"firetrack/internal/auth"

	"github.com/gin-gonic/gin"
)

const ContextNoteKey = "note"

func RequireNoteOwnerOrAdmin(repo *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		noteID := c.Param("noteId")
		if noteID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing note id",
			})
			c.Abort()
			return
		}

		note, err := repo.GetByID(noteID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch note",
			})
			c.Abort()
			return
		}

		if note == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "note not found",
			})
			c.Abort()
			return
		}

		rawUser, exists := c.Get(auth.ContextUserKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		user, ok := rawUser.(*auth.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid auth context",
			})
			c.Abort()
			return
		}

		isOwner := note.AuthorID == user.ID
		isAdmin := user.Role == "admin"

		if !isOwner && !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			c.Abort()
			return
		}

		c.Set(ContextNoteKey, note)
		c.Next()
	}
}
