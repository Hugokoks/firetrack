package jobs

import (
	"net/http"

	"firetrack/internal/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateJobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	rawUser, exists := c.Get(auth.ContextUserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user, ok := rawUser.(*auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid auth context",
		})
		return
	}

	job, err := h.service.Create(input, user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create job",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "job created successfully",
		"job":     job,
	})
}

func (h *Handler) GetAll(c *gin.Context) {
	jobs, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch jobs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing job id",
		})
		return
	}

	job, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch job",
		})
		return
	}

	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "job not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job": job,
	})
}
