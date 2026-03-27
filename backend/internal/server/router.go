package server

import "github.com/gin-gonic/gin"

func NewRouter(app *App) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")

	// AUTH
	auth := api.Group("/auth")
	{
		auth.POST("/login", app.AuthHandler.Login)
		auth.GET("/me", app.AuthMiddleware, app.AuthHandler.Me)
		auth.POST("/logout", app.AuthMiddleware, app.AuthHandler.Logout)
	}

	// JOBS
	jobs := api.Group("/jobs")
	jobs.Use(app.AuthMiddleware, app.ActivityMiddleware)
	{	
		// JOBS
		jobs.POST("", app.JobsHandler.Create)
		jobs.GET("", app.JobsHandler.GetAll)
		jobs.GET("/:id", app.JobsHandler.GetByID)
		jobs.PATCH("/:id", app.JobsHandler.Update)

		// NOTES
		jobs.POST("/:id/notes", app.NotesHandler.Create)
		jobs.GET("/:id/notes", app.NotesHandler.GetByJobID)
		jobs.DELETE("/:id/notes/:noteId",
			app.NoteOwnerCheck,
			app.NotesHandler.Delete,
		)
		jobs.PATCH("/:id/notes/:noteId", app.NoteOwnerCheck, app.NotesHandler.Update)

		// FILES
		jobs.POST("/:id/files", app.FilesHandler.Create)
		jobs.GET("/:id/files", app.FilesHandler.GetByJobID)
		jobs.GET("/:id/files/:fileId/view",app.FilesHandler.View)
		jobs.GET("/:id/files/:fileId/download",app.FilesHandler.Download)
		jobs.DELETE("/:id/files/:fileId", app.FilesHandler.Delete)
		
	}

	return r
}
