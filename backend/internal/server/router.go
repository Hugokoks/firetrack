package server

import (
	"database/sql"
	"firetrack/internal/auth"
	"firetrack/internal/jobs"
	"firetrack/internal/notes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	r := gin.New()

	// middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	jobsRepo := jobs.NewRepository(db)
	jobsService := jobs.NewService(jobsRepo)
	jobsHandler := jobs.NewHandler(jobsService)

	notesRepo := notes.NewRepositry(db)
	notesService := notes.NewService(notesRepo)
	notesHandler := notes.NewHandler(notesService)
	// API group
	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", auth.AuthMiddleware(authRepo), authHandler.Me)
			authGroup.POST("/logout", auth.AuthMiddleware(authRepo), authHandler.Logout)
		}
		jobsGroup := api.Group("/jobs")
		{
			jobsGroup.POST("", auth.AuthMiddleware(authRepo), jobsHandler.Create)
			jobsGroup.GET("", auth.AuthMiddleware(authRepo), jobsHandler.GetAll)
			jobsGroup.GET("/:id", auth.AuthMiddleware(authRepo), jobsHandler.GetByID)

			jobsGroup.POST("/:id/notes", auth.AuthMiddleware(authRepo), notesHandler.Create)
			jobsGroup.GET("/:id/notes", auth.AuthMiddleware(authRepo), notesHandler.GetByJobID)
			jobsGroup.DELETE("/:id/notes/:noteId",
				auth.AuthMiddleware(authRepo),
				notes.RequireNoteOwnerOrAdmin(notesRepo), notesHandler.Delete)

		}
	}

	return r
}
