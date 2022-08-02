package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/J-R-Oliver/dwp-assessment-go/internal/configuration"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/dwp"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
)

const ContentTypeApplicationJSON = "application/json"

type service interface {
	RetrievePeople(ctx context.Context) (dwp.People, error)
	RetrievePeopleByCity(ctx context.Context, city string, distance int) (dwp.People, error)
}

type Handlers struct {
	Service         service
	DefaultDistance int
	Cities          map[string]configuration.City
	Logger          logging.Logger
}

func (h Handlers) GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ContentTypeApplicationJSON)

	if r.Method != http.MethodGet {
		h.methodNotAllow(w, r)
		return
	}

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	people, err := h.Service.RetrievePeople(ctx)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	err = json.NewEncoder(w).Encode(people)
	if err != nil {
		h.InternalServerError(w, r, err)
	}
}

func (h Handlers) GetPeopleByCity(pathPrefix string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			h.methodNotAllow(w, r)
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
			distance = h.DefaultDistance
		}

		if err != nil {
			h.Logger.Info(fmt.Sprintf("bad query: %s distance", distanceQuery))
			h.badRequest(w, r, fmt.Sprintf("Invalid distance query - %s is not an integer", distanceQuery))

			return
		}

		_, ok := h.Cities[path]
		if !ok {
			h.Logger.Info(fmt.Sprintf("city not found - %s", path))
			h.notFound(w, r, "City Not Found")

			return
		}

		people, err := h.Service.RetrievePeopleByCity(ctx, path, distance)
		if err != nil {
			h.InternalServerError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", ContentTypeApplicationJSON)

		err = json.NewEncoder(w).Encode(people)
		if err != nil {
			h.InternalServerError(w, r, err)
		}
	}
}

func (h Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
