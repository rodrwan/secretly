package web

import (
	"net/http"

	"github.com/rodrwan/secretly/internal/web/templates"
)

// Handler maneja las rutas web
type Handler struct {
	envManager interface {
		Load() (map[string]string, error)
		Save(map[string]string) error
	}
}

// NewHandler crea una nueva instancia del manejador web
func NewHandler(envManager interface {
	Load() (map[string]string, error)
	Save(map[string]string) error
}) *Handler {
	return &Handler{
		envManager: envManager,
	}
}

// RegisterRoutes registra las rutas web
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("internal/web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Rutas de la aplicación
	mux.HandleFunc("/", h.handleIndex)
}

// handleIndex maneja la ruta principal
func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	component := templates.Index()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
