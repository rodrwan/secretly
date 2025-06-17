package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/rodrwan/secretly/cmd/server/handlers"
	"github.com/rodrwan/secretly/internal/config"
	"github.com/rodrwan/secretly/internal/database"
	"github.com/rodrwan/secretly/internal/env"
	"github.com/rodrwan/secretly/internal/web"

	_ "modernc.org/sqlite"
)

func main() {
	// Load configuration
	cfg := config.New()
	envManager := env.NewManager(cfg.EnvPath)

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	queries := database.New(db)

	// Run migrations
	goose.SetBaseFS(database.Migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}
	goose.SetLogger(log.New(os.Stdout, "goose: ", log.LstdFlags))

	// Run migrations
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}

	// Server configuration
	router := http.NewServeMux()

	// Configure web handler
	webHandler := web.NewHandler(envManager, queries)
	webHandler.RegisterRoutes(router)
	handlers.RegisterRoutes(router, queries)

	// Start server
	log.Printf("Server started on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
