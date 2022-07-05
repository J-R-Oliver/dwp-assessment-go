package main

import (
	"net/http"

	"github.com/J-R-Oliver/dwp-assessment-go/dwp"
	"github.com/J-R-Oliver/dwp-assessment-go/internal/people"
	"github.com/J-R-Oliver/dwp-assessment-go/logging"
)

func main() {
	c := loadConfiguration()

	log := logging.New()

	client := dwp.NewClient(c.PeopleConfiguration.BaseURL)
	s := people.Service{DwpClient: client}
	h := handlers{service: s, defaultDistance: c.PeopleConfiguration.Distance}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/api/people", h.getPeople)
	serveMux.HandleFunc("/api/people/", h.getPeopleByCity("/api/people/"))

	handler := panicHandler(serveMux)
	logRequestHandler(handler)

	log.Info("Starting server on :8080")

	err := http.ListenAndServe(":8080", serveMux)
	log.Fatal(err)
}
