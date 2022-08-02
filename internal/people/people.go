package people

import (
	"context"
	"fmt"

	"github.com/J-R-Oliver/dwp-assessment-go/pkg/dwp"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
	"github.com/umahmood/haversine"
	"golang.org/x/sync/errgroup"
)

type peopleClient interface {
	RetrievePeople(ctx context.Context) (dwp.People, error)
	RetrievePeopleByCity(ctx context.Context, city string) (dwp.People, error)
}

type Service struct {
	DwpClient peopleClient
	Cities    map[string]haversine.Coord
	Logger    logging.Logger
}

func (s Service) RetrievePeople(ctx context.Context) (dwp.People, error) {
	s.Logger.Info("Attempting to retrieve all people")

	people, err := s.DwpClient.RetrievePeople(ctx)
	if err != nil {
		return nil, err
	}

	s.Logger.Info("All people retrieved successfully")

	return people, nil
}

func (s Service) RetrievePeopleByCity(ctx context.Context, city string, distance int) (dwp.People, error) {
	cityCoordinates, ok := s.Cities[city]
	if !ok {
		return nil, fmt.Errorf("%s's coordinates have not been configured", city)
	}

	c := make(chan dwp.People, 2) //nolint:gomnd

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		s.Logger.Info("Attempting to retrieve all people")

		people, err := s.DwpClient.RetrievePeople(ctx)
		if err != nil {
			return err
		}

		s.Logger.Info("All people retrieved successfully")
		c <- filterPeople(people, distance, cityCoordinates)

		return nil
	})

	eg.Go(func() error {
		s.Logger.Info("Attempting to retrieve people by city")
		cityPeople, err := s.DwpClient.RetrievePeopleByCity(ctx, city)

		s.Logger.Info("People by city retrieved successfully")
		c <- cityPeople

		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	close(c)

	var allPeople dwp.People
	for people := range c {
		allPeople = append(allPeople, people...)
	}

	return allPeople, nil
}

func filterPeople(people dwp.People, distance int, cityCoordinates haversine.Coord) dwp.People {
	var filteredPeople dwp.People

	for _, person := range people {
		personCoordinates := haversine.Coord{
			Lat: float64(person.Latitude),
			Lon: float64(person.Longitude),
		}

		miles, _ := haversine.Distance(cityCoordinates, personCoordinates)

		if miles <= float64(distance) {
			filteredPeople = append(filteredPeople, person)
		}
	}

	return filteredPeople
}
