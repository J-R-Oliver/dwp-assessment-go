package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/J-R-Oliver/dwp-assessment-go/internal/configuration"
	"github.com/J-R-Oliver/dwp-assessment-go/internal/handler"
	"github.com/J-R-Oliver/dwp-assessment-go/internal/middleware"
	"github.com/J-R-Oliver/dwp-assessment-go/internal/people"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/dwp"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
	"github.com/umahmood/haversine"
)

func main() {
	configurationPath, ok := os.LookupEnv("CONFIGURATION_PATH")
	if !ok {
		configurationPath = "./configuration.yaml"
	}

	c, err := configuration.LoadConfiguration(configurationPath)
	if err != nil {
		log.Fatal(err)
	}

	l := logging.New(c.LoggingLevel)

	cities := convertCities(c)

	client := dwp.NewClient(c.PeopleConfiguration.BaseURL, http.Client{})

	s := people.Service{
		DwpClient: client,
		Cities:    cities,
		Logger:    l,
	}

	h := handler.Handlers{
		Service:         s,
		DefaultDistance: c.PeopleConfiguration.Distance,
		Cities:          c.Cities,
		Logger:          l,
	}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/api/people", h.GetPeople)
	serveMux.HandleFunc("/api/people/", h.GetPeopleByCity("/api/people/"))
	serveMux.HandleFunc("/health", h.Health)
	serveMux.HandleFunc("/", h.NotFound)

	middlewareChain := middleware.PanicHandler(serveMux, h.InternalServerError)
	middlewareChain = middleware.LogRequestHandler(middlewareChain, l)

	srv := &http.Server{
		Addr:              ":" + c.Port,
		Handler:           middlewareChain,
		ReadHeaderTimeout: time.Minute,
	}

	l.Info("Starting server on :" + c.Port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func convertCities(c configuration.Configuration) map[string]haversine.Coord {
	cities := make(map[string]haversine.Coord)
	bitSize := 64

	for cityName, coordinates := range c.Cities {
		lat, err := strconv.ParseFloat(coordinates.Latitude, bitSize)
		if err != nil {
			log.Fatalf("fatal error: unable to parse %s latitude %s to float64: %s", cityName, coordinates.Latitude, err)
		}

		lon, err := strconv.ParseFloat(coordinates.Longitude, bitSize)
		if err != nil {
			log.Fatalf("fatal error: unable to parse %s longtitude %s to float64: %s", cityName, coordinates.Latitude, err)
		}

		cities[cityName] = haversine.Coord{
			Lat: lat,
			Lon: lon,
		}
	}

	return cities
}
