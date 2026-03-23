package server

import (
	"database/sql"
	"firetrack/internal/activity"
	"firetrack/internal/auth"
	"firetrack/internal/config"
	"firetrack/internal/files"
	"firetrack/internal/jobs"
	"firetrack/internal/notes"

	"github.com/gin-gonic/gin"
)

type App struct {
	AuthHandler  *auth.Handler
	JobsHandler  *jobs.Handler
	NotesHandler *notes.Handler
	FilesHandler *files.Handler

	AuthMiddleware     gin.HandlerFunc
	NoteOwnerCheck     gin.HandlerFunc
	ActivityMiddleware gin.HandlerFunc
}

func NewApp(db *sql.DB, cfg *config.Config) *App {
	// AUTH
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// JOBS
	jobsRepo := jobs.NewRepository(db)
	jobsService := jobs.NewService(jobsRepo)
	jobsHandler := jobs.NewHandler(jobsService)

	// NOTES
	notesRepo := notes.NewRepository(db)
	notesService := notes.NewService(notesRepo)
	notesHandler := notes.NewHandler(notesService)

	// ACTIVITY
	activityRepo := activity.NewRepository(db)
	activityService := activity.NewService(activityRepo)
	activityMiddleware := activity.Middleware(activityService)

	// FILES
	filesRepo := files.NewRepository(db)
	filesStorage := files.NewStorage(cfg.FilesRoot)
	filesService := files.NewService(filesRepo, filesStorage, jobsRepo, 10<<20) // 10 MB
	filesHandler := files.NewHandler(filesService)

	return &App{
		AuthHandler:  authHandler,
		JobsHandler:  jobsHandler,
		NotesHandler: notesHandler,
		FilesHandler: filesHandler,

		AuthMiddleware:     auth.AuthMiddleware(authRepo),
		NoteOwnerCheck:     notes.RequireNoteOwnerOrAdmin(notesRepo),
		ActivityMiddleware: activityMiddleware,
	}
}
