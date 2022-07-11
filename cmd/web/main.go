package main

import (
	"log"
	"net/http"

	"github.com/J-R-Oliver/dwp-assessment-go/internal/people"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/dwp"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
)

func main() {
	c := loadConfiguration()

	l := logging.New(logging.Info) // ToDo - this should be configured by configuration

	httpClient := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	client := dwp.NewClient(c.PeopleConfiguration.BaseURL, httpClient)
	s := people.Service{DwpClient: client}
	h := handlers{service: s, defaultDistance: c.PeopleConfiguration.Distance}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/api/people", h.getPeople)
	serveMux.HandleFunc("/api/people/", h.getPeopleByCity("/api/people/"))

	handler := panicHandler(serveMux)
	logRequestHandler(handler)

	srv := &http.Server{
		Addr: ":" + c.Port,
		//ErrorLog: log.Error, // ToDo - custom error logging?
		Handler: serveMux,
	}

	l.Info("Starting server on :" + c.Port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
