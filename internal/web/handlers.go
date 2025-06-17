package web

import (
	"net/http"

	"github.com/rodrwan/secretly/internal/database"
	"github.com/rodrwan/secretly/internal/env"
	"github.com/rodrwan/secretly/internal/web/templates"
)

// Handler maneja las rutas web
type Handler struct {
	envManager *env.Manager
	queries    *database.Queries
}

// NewHandler crea una nueva instancia del manejador web
func NewHandler(envManager *env.Manager, queries *database.Queries) *Handler {
	return &Handler{
		envManager: envManager,
		queries:    queries,
	}
}

// RegisterRoutes registra las rutas web
func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("internal/web/static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Rutas de la aplicación
	router.HandleFunc("/", h.handleIndex)
}

// handleIndex maneja la ruta principal
func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	component := templates.Index()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
