package httputil

import (
	"net/http"

	"firetrack/internal/auth"

	"github.com/gin-gonic/gin"
)

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": message,
	})
}

func GetCurrentUser(c *gin.Context) (*auth.User, bool) {
	rawUser, exists := c.Get(auth.ContextUserKey)
	if !exists {
		RespondError(c, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}

	user, ok := rawUser.(*auth.User)
	if !ok {
		RespondError(c, http.StatusInternalServerError, "invalid auth context")
		return nil, false
	}

	return user, true
}

func RequireParam(c *gin.Context, key, errMsg string) (string, bool) {
	value := c.Param(key)
	if value == "" {
		RespondError(c, http.StatusBadRequest, errMsg)
		return "", false
	}
	return value, true
}
