package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rodrwan/secretly/internal/database"
	"go.uber.org/zap"
)

type handlerFunc func(db database.Querier, w http.ResponseWriter, r *http.Request) (Response, error)

type Handler struct {
	db database.Querier
}

func NewHandler(db database.Querier) *Handler {
	return &Handler{db: db}
}

func (eh *Handler) Call(handler handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := handler(eh.db, w, r)
		if err != nil {
			Error(w, r, resp.Code, resp.Message, err)
			return
		}

		Success(w, r, resp.Code, resp.Message, resp.Data)
	}
}

type Response struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
}

func Success(w http.ResponseWriter, r *http.Request, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, r *http.Request, code int, message string, err error) {
	zap.L().Error("Error",
		zap.String("path", r.URL.Path),
		zap.Int("code", code),
		zap.String("message", message),
		zap.Error(err),
	)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Code:  code,
		Error: message,
	})
}
