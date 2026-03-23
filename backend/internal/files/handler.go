package files

import (
	"net/http"

	"firetrack/internal/activity"
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
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing job id",
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

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing file",
		})
		return
	}

	file, err := h.service.Create(CreateFileInput{
		JobID:      jobID,
		UploadedBy: user.ID,
		FileHeader: fileHeader,
	})
	if err != nil {
		switch err {
		case ErrMissingFile:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case ErrFileTooLarge:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case ErrFileTypeNotAllowed:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case ErrJobNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to upload file",
			})
			return
		}
	}

	activity.Set(c, activity.Payload{
		JobID:       file.JobID,
		UserID:      user.ID,
		ActionType:  "file_uploaded",
		ActionLabel: "File uploaded",
		Meta: gin.H{
			"file_id":   file.ID,
			"file_name": file.FileName,
			"mime_type": file.MimeType,
			"file_size": file.FileSize,
		},
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "file uploaded successfully",
		"file":    file,
	})
}
