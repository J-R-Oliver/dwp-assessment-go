package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type errorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
}

func (h Handlers) NotFound(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info(fmt.Sprintf("route not found - %s", r.URL.Path))
	h.notFound(w, r, "Route Not Found")
}

func (h Handlers) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Error(err)

	h.errorHandler(w, r, http.StatusInternalServerError, "Internal Server Error")
}

func (h Handlers) badRequest(w http.ResponseWriter, r *http.Request, message string) {
	h.errorHandler(w, r, http.StatusBadRequest, message)
}

func (h Handlers) notFound(w http.ResponseWriter, r *http.Request, message string) {
	h.errorHandler(w, r, http.StatusNotFound, message)
}

func (h Handlers) methodNotAllow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", http.MethodGet)
	h.errorHandler(w, r, http.StatusMethodNotAllowed, "Method Not Allowed")
}

func (h Handlers) errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", ContentTypeApplicationJSON)

	w.WriteHeader(status)

	response := errorResponse{
		Timestamp: time.Now(),
		Status:    status,
		Message:   message,
		Path:      r.URL.Path,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
