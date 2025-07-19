package main

import (
	"log"

	"github.com/sinasadeghi83/aut-grader/internal/api/config"
	"github.com/sinasadeghi83/aut-grader/internal/api/platform/database"
	"github.com/sinasadeghi83/aut-grader/internal/api/server"
)

func main() {
	cfg := config.LoadConfig()

	//Connect to DB
	db, err := database.OpenDatabase(cfg.DbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	srv := server.NewServer(cfg.ServerPort, db, cfg)
	srv.SetupRoutes()

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
