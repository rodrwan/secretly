package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rodrwan/secretly/internal/config"
	"github.com/rodrwan/secretly/internal/env"
	"github.com/rodrwan/secretly/internal/web"
)

func main() {
	// Load configuration
	cfg := config.New()
	envManager := env.NewManager(cfg.EnvPath)

	// Server configuration
	mux := http.NewServeMux()

	// Configure web handler
	webHandler := web.NewHandler(envManager)
	webHandler.RegisterRoutes(mux)

	// API endpoints
	mux.HandleFunc(cfg.BasePath+"/env", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			vars, err := envManager.Load()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(vars)

		case http.MethodPost:
			var vars map[string]string
			if err := json.NewDecoder(r.Body).Decode(&vars); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := envManager.Save(vars); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start server
	log.Printf("Server started on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
