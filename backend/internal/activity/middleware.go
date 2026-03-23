package activity

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Set(c *gin.Context, payload Payload) {
	c.Set(ContextActivityKey, payload)
}

func Middleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() >= http.StatusBadRequest {
			return
		}

		raw, exists := c.Get(ContextActivityKey)
		if !exists {
			return
		}

		payload, ok := raw.(Payload)
		if !ok {
			return
		}

		_ = service.Log(payload)
	}
}
