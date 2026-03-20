package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const SessionCookieName = "session_id"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	result, err := h.service.Login(input.Email, input.Password)
	fmt.Println(err)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})
			return
		case ErrInactiveUser:
			c.JSON(http.StatusForbidden, gin.H{
				"error": "user account is inactive",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to login",
			})
			return
		}
	}

	secureCookie := false // na Railway/HTTPS pak dej true

	c.SetCookie(
		SessionCookieName,
		result.Session.ID,
		7*24*60*60, // 7 dní
		"/",
		"",
		secureCookie,
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": gin.H{
			"id":    result.User.ID,
			"name":  result.User.Name,
			"email": result.User.Email,
			"role":  result.User.Role,
		},
	})
}

func (h *Handler) Me(c *gin.Context) {
	rawUser, exists := c.Get(ContextUserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user, ok := rawUser.(*User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid auth context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}