package main

import (
	"log"

	"firetrack/internal/config"
	"firetrack/internal/db"
	"firetrack/internal/server"
)

func main() {
	// load config
	cfg := config.Load()
	// connect DB
	database, err := db.NewPostgres(cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	defer database.Close()

	app := server.NewApp(database, cfg)
	r := server.NewRouter(app)

	log.Println("🚀 Server running on", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
