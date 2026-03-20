package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const ContextUserKey = "auth_user"

func AuthMiddleware(repo *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(SessionCookieName)
		if err != nil || sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		user, err := repo.GetUserBySessionID(sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to validate session",
			})
			c.Abort()
			return
		}

		if user == nil || !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		c.Set(ContextUserKey, user)
		c.Next()
	}
}