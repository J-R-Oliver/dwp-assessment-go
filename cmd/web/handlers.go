package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/J-R-Oliver/dwp-assessment-go/dwp"
)

// func _createErrorResponse(status int, //nolint:deadcode
//	message string,
//	path string) string {
//	response := errorResponse{
//		time.Now(),
//		status,
//		message,
//		path,
//	}
//
//	marshal, err := json.Marshal(response)
//	if err != nil {
//		return ""
//	}
//
//	return string(marshal)
// }

type service interface {
	RetrievePeople(ctx context.Context) (dwp.People, error)
	RetrievePeopleByCity(ctx context.Context, city string, distance int) (dwp.People, error)
}

type handlers struct {
	service         service
	defaultDistance int
}

func (h handlers) getPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		methodNotAllow(w, r)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	people, err := h.service.RetrievePeople(ctx)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	err = json.NewEncoder(w).Encode(people)
	if err != nil {
		internalServerError(w, r, err)
	}
}

func (h handlers) getPeopleByCity(pathPrefix string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllow(w, r)
			return
		}

		ctx := r.Context()

		path := strings.TrimPrefix(r.URL.Path, pathPrefix)
		path = strings.ToUpper(path[:1]) + path[1:]

		var distance int

		var err error

		query := r.URL.Query()

		distanceQuery := query.Get("distance")
		if distanceQuery != "" {
			distance, err = strconv.Atoi(distanceQuery)
		} else {
			distance = h.defaultDistance
		}

		if err != nil {
			badRequest(w, r, errors.New("bad query: distance"))
			return
		}

		people, err := h.service.RetrievePeopleByCity(ctx, path, distance)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(people)
		if err != nil {
			internalServerError(w, r, err)
		}
	}
}
