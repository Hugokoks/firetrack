package files

import (
	"fmt"
	"io"
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

	user, ok := httputil.GetCurrentUser(c)
	if !ok {
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

func (h *Handler) GetByJobID(c *gin.Context) {

	jobID, ok := httputil.RequireParam(c, "id", "missing job id")
	if !ok {

		return
	}

	files, err := h.service.GetByJobID(jobID)

	if err != nil {
		switch err {
		case ErrJobNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "job not found",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch files",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{

		"files": files,
	})

}

func (h *Handler) View(c *gin.Context) {
	
	jobID, ok := httputil.RequireParam(c, "id", "missing job id")
	if !ok {

		return
	}

	fileID, ok := httputil.RequireParam(c, "fileId", "missing job id")
	if !ok {

		return
	}

	file, reader, err := h.service.View(jobID, fileID)
	
	if err != nil {
		switch err {
		case ErrJobNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "job not found",
			})
			return
		case ErrFileNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "file not found",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to open file",
			})
			return
		}
	}
	defer reader.Close()

	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", `inline; filename="`+file.FileName+`"`)
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Length", fmt.Sprintf("%d", file.FileSize))

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}


func (h *Handler) Download(c *gin.Context) {

	jobID, ok := httputil.RequireParam(c, "id", "missing job id")
	if !ok {
		return
	}
	fileID, ok := httputil.RequireParam(c, "fileId", "missing job id")
	if !ok {
		return
	}

	file, reader, err := h.service.Download(jobID, fileID)
	if err != nil {
		switch err {
		case ErrJobNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "job not found",
			})
			return
		case ErrFileNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "file not found",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to download file",
			})
			return
		}
	}
	defer reader.Close()

	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", `attachment; filename="`+file.FileName+`"`)
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Length", fmt.Sprintf("%d", file.FileSize))

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Delete(c *gin.Context) {
	jobID, ok := httputil.RequireParam(c, "id", "missing job id")
	if !ok {
		return
	}
	fileID, ok := httputil.RequireParam(c, "fileId", "missing file id")
	if !ok {
		return
	}
	user, ok := httputil.GetCurrentUser(c)
	if !ok {
		return
	}
	file, err := h.service.Delete(jobID, fileID)
	if err != nil {
		switch err {
		case ErrJobNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "job not found",
			})
			return
		case ErrFileNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "file not found",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete file",
			})
			return
		}
	}
	activity.Set(c, activity.Payload{
		JobID:       jobID,
		UserID:      user.ID,
		ActionType:  "file_deleted",
		ActionLabel: "File deleted",
		Meta: gin.H{
			"file_id":   file.ID,
			"file_name": file.FileName,
		},
	})
	
	c.JSON(http.StatusOK, gin.H{
		"message": "file deleted successfully",
	})
}