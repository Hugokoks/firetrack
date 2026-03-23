package notes

import (
	"net/http"

	"firetrack/internal/httputil"

	"github.com/gin-gonic/gin"
)

func getContextNote(c *gin.Context) (*Note, bool) {
	rawNote, exists := c.Get(ContextNoteKey)
	if !exists {
		httputil.RespondError(c, http.StatusInternalServerError, "note not found in context")
		return nil, false
	}

	note, ok := rawNote.(*Note)
	if !ok {
		httputil.RespondError(c, http.StatusInternalServerError, "invalid note context")
		return nil, false
	}

	return note, true
}
