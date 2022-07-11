package main

import (
	"encoding/json"
	"errors"
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

func badRequest(w http.ResponseWriter, r *http.Request, err error) {
	errorHandler(w, r, err, http.StatusBadRequest)
}

// func _notFound(w http.ResponseWriter, r *http.Request, err error) { //nolint:deadcode
//	errorHandler(w, r, err, http.StatusNotFound)
// }

func methodNotAllow(w http.ResponseWriter, r *http.Request) {
	// ToDo - sort out the error message here
	w.Header().Set("Allow", http.MethodGet)
	errorHandler(w, r, errors.New("method not allowed"), http.StatusMethodNotAllowed)
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	errorHandler(w, r, err, http.StatusInternalServerError)
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error, status int) {
	// ToDo - log the error!
	fmt.Println(err)

	w.Header().Set("Content-Type", "application/json")

	// ToDo - Should probably store the StatusText on an Error key and a response error on message
	response := errorResponse{
		Timestamp: time.Now(),
		Status:    status,
		Message:   http.StatusText(status),
		Path:      r.URL.Path,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
