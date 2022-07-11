package people

import (
	"context"

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
	logger    logging.Logger
}

func (s Service) RetrievePeople(ctx context.Context) (dwp.People, error) {
	people, err := s.DwpClient.RetrievePeople(ctx)
	if err != nil {
		return nil, err
	}

	return people, nil
}

var london = haversine.Coord{
	Lat: 51.514248, //nolint:gomnd // ToDo - fix this error
	Lon: -0.093145,
}

func (s Service) RetrievePeopleByCity(ctx context.Context, city string, distance int) (dwp.People, error) {
	c := make(chan dwp.People, 2)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		people, err := s.DwpClient.RetrievePeople(ctx)
		if err != nil {
			return err
		}

		c <- filterPeople(people, distance)

		return nil
	})

	eg.Go(func() error {
		cityPeople, err := s.DwpClient.RetrievePeopleByCity(ctx, city)

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

func filterPeople(people dwp.People, distance int) dwp.People {
	var filteredPeople dwp.People

	for _, person := range people {
		coord := haversine.Coord{
			Lat: float64(person.Latitude),
			Lon: float64(person.Longitude),
		}

		miles, _ := haversine.Distance(london, coord)

		if miles <= float64(distance) {
			filteredPeople = append(filteredPeople, person)
		}
	}

	return filteredPeople
}
