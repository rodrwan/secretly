package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/rodrwan/secretly/internal/database"
	"github.com/rodrwan/secretly/internal/web/templates"
)

//go:embed static
var staticFiles embed.FS

// Handler maneja las rutas web
type Handler struct {
	queries *database.Queries
}

// NewHandler crea una nueva instancia del manejador web
func NewHandler(queries *database.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

// RegisterRoutes registra las rutas web
func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	// Servir archivos estáticos
	staticFs, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
		return
	}

	fs := http.FileServer(http.FS(staticFs))
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
