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
	"github.com/rodrwan/secretly/internal/web"

	_ "modernc.org/sqlite"
)

func main() {
	// Load configuration
	cfg := config.New()

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
	webHandler := web.NewHandler(queries)
	webHandler.RegisterRoutes(router)
	handlers.RegisterRoutes(router, queries)

	// Wrap the router with middleware
	handler := panicMiddleware(router)

	// Start server
	log.Printf("Server started on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}

// Create middleware wrapper to handle panics
func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
