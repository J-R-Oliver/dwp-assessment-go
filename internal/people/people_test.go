package people

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/J-R-Oliver/dwp-assessment-go/pkg/dwp"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
	"github.com/umahmood/haversine"
)

var mockRetrievePeople func() (dwp.People, error)
var mockRetrievePeopleByCity func() (dwp.People, error)

type MockDwpClient struct{}

func (m MockDwpClient) RetrievePeople(ctx context.Context) (dwp.People, error) {
	return mockRetrievePeople()
}

func (m MockDwpClient) RetrievePeopleByCity(ctx context.Context, city string) (dwp.People, error) {
	return mockRetrievePeopleByCity()
}

func TestService_RetrievePeople(t *testing.T) {
	t.Run("Given RetrievePeople is invoked when DwpClient returns people then people are returned", func(t *testing.T) {
		expectedPeople := dwp.People{
			{
				ID:        1,
				FirstName: "Maurise",
				LastName:  "Shieldon",
				Email:     "mshieldon0@squidoo.com",
				IPAddress: "192.57.232.111",
				Latitude:  dwp.Coordinate(34.003135),
				Longitude: dwp.Coordinate(-117.7228641),
			},
			{
				ID:        2,
				FirstName: "Bendix",
				LastName:  "Halgarth",
				Email:     "bhalgarth1@timesonline.co.uk",
				IPAddress: "4.185.73.82",
				Latitude:  dwp.Coordinate(-2.9623869),
				Longitude: dwp.Coordinate(104.7399789),
			},
		}

		mockRetrievePeople = func() (dwp.People, error) {
			return expectedPeople, nil
		}

		m := MockDwpClient{}

		s := Service{
			DwpClient: m,
			Cities:    nil,
			Logger:    logging.New(logging.Info),
		}

		actualPeople, err := s.RetrievePeople(context.Background())

		if err != nil {
			t.Errorf("RetrievePeople() error = %v", err)
		}

		if !reflect.DeepEqual(expectedPeople, actualPeople) {
			t.Errorf("RetrievePeople() = %v, want %v", actualPeople, expectedPeople)
		}
	})

	t.Run("Given RetrievePeople is invoked when DwpClient returns err then err is returned", func(t *testing.T) {
		expectedErr := errors.New("test error")

		mockRetrievePeople = func() (dwp.People, error) {
			return nil, expectedErr
		}

		m := MockDwpClient{}

		s := Service{
			DwpClient: m,
			Cities:    nil,
			Logger:    logging.New(logging.Info),
		}

		people, err := s.RetrievePeople(context.Background())

		if !errors.Is(err, expectedErr) {
			t.Errorf("RetrievePeople() error = %v, want %v", err, expectedErr)
		}

		if people != nil {
			t.Errorf("RetrievePeople() people = %v, want nil", people)
		}
	})
}

func TestService_RetrievePeopleByCity(t *testing.T) {
	t.Run("Given city has been configured when RetrievePeople and RetrievePeopleByCity are successful then returns filtered people", func(t *testing.T) {
		p := dwp.People{
			{
				ID:        2,
				FirstName: "Bendix",
				LastName:  "Halgarth",
				Email:     "bhalgarth1@timesonline.co.uk",
				IPAddress: "4.185.73.82",
				Latitude:  dwp.Coordinate(-2.9623869),
				Longitude: dwp.Coordinate(104.7399789),
			},
			{
				ID:        1,
				FirstName: "Maurise",
				LastName:  "Shieldon",
				Email:     "mshieldon0@squidoo.com",
				IPAddress: "192.57.232.111",
				Latitude:  dwp.Coordinate(51.6553959),
				Longitude: dwp.Coordinate(0.0572553),
			},
			{
				ID:        135,
				FirstName: "Mechelle",
				LastName:  "Boam",
				Email:     "mboam3q@thetimes.co.uk",
				IPAddress: "113.71.242.187",
				Latitude:  dwp.Coordinate(-6.5115909),
				Longitude: dwp.Coordinate(105.652983),
			},
			{
				ID:        396,
				FirstName: "Terry",
				LastName:  "Stowgill",
				Email:     "tstowgillaz@webeden.co.uk",
				IPAddress: "143.190.50.240",
				Latitude:  dwp.Coordinate(-6.7098551),
				Longitude: dwp.Coordinate(111.3479498),
			},
		}

		mockRetrievePeople = func() (dwp.People, error) {
			return p[:2], nil
		}

		mockRetrievePeopleByCity = func() (dwp.People, error) {
			return p[2:], nil
		}

		m := MockDwpClient{}

		s := Service{
			DwpClient: m,
			Cities:    map[string]haversine.Coord{"london": {Lat: 51.514248, Lon: -0.093145}},
			Logger:    logging.New(logging.Info),
		}

		actualPeople, err := s.RetrievePeopleByCity(context.Background(), "london", 50)

		if err != nil {
			t.Errorf("RetrievePeople() error = %v", err)
		}

		var hasID1, hasID135, hasID396 bool

		for _, person := range p {
			if person.ID == 1 {
				hasID1 = true
			}

			if person.ID == 135 {
				hasID135 = true
			}

			if person.ID == 396 {
				hasID396 = true
			}
		}

		if !hasID1 || !hasID135 || !hasID396 || len(actualPeople) != 3 {
			t.Errorf("RetrievePeople() = %v, want %v", actualPeople, p[1:])
		}
	})

	t.Run("Given city has not been configured then returns error", func(t *testing.T) {
		s := Service{
			DwpClient: nil,
			Cities:    map[string]haversine.Coord{},
			Logger:    logging.New(logging.Info),
		}

		p, err := s.RetrievePeopleByCity(context.Background(), "timbuctoo", 50)

		if err.Error() != "timbuctoo's coordinates have not been configured" {
			t.Errorf("RetrievePeople() error = %v, want = timbuctoo's coordinates have not been configured", err)
		}

		if p != nil {
			t.Errorf("RetrievePeople() = %v, want nil", p)
		}
	})

	t.Run("Given city has been configured when RetrievePeople is unsuccessful then returns error", func(t *testing.T) {
		expectedError := errors.New("test error")

		mockRetrievePeople = func() (dwp.People, error) {
			return nil, expectedError
		}

		mockRetrievePeopleByCity = func() (dwp.People, error) {
			p := dwp.People{
				{
					ID:        2,
					FirstName: "Bendix",
					LastName:  "Halgarth",
					Email:     "bhalgarth1@timesonline.co.uk",
					IPAddress: "4.185.73.82",
					Latitude:  dwp.Coordinate(-2.9623869),
					Longitude: dwp.Coordinate(104.7399789),
				},
			}
			return p, nil
		}

		m := MockDwpClient{}

		s := Service{
			DwpClient: m,
			Cities:    map[string]haversine.Coord{"london": {Lat: 51.514248, Lon: -0.093145}},
			Logger:    logging.New(logging.Info),
		}

		p, err := s.RetrievePeopleByCity(context.Background(), "london", 50)

		if !errors.Is(err, expectedError) {
			t.Errorf("RetrievePeople() error = %v, want = %v", err, expectedError)
		}

		if p != nil {
			t.Errorf("RetrievePeople() = %v, want nil", p)
		}
	})

	t.Run("Given city has been configured when RetrievePeopleByCity is unsuccessful then returns error", func(t *testing.T) {
		expectedError := errors.New("test error")

		mockRetrievePeople = func() (dwp.People, error) {
			p := dwp.People{
				{
					ID:        2,
					FirstName: "Bendix",
					LastName:  "Halgarth",
					Email:     "bhalgarth1@timesonline.co.uk",
					IPAddress: "4.185.73.82",
					Latitude:  dwp.Coordinate(-2.9623869),
					Longitude: dwp.Coordinate(104.7399789),
				},
			}
			return p, nil
		}

		mockRetrievePeopleByCity = func() (dwp.People, error) {
			return nil, expectedError
		}

		m := MockDwpClient{}

		s := Service{
			DwpClient: m,
			Cities:    map[string]haversine.Coord{"london": {Lat: 51.514248, Lon: -0.093145}},
			Logger:    logging.New(logging.Info),
		}

		p, err := s.RetrievePeopleByCity(context.Background(), "london", 50)

		if !errors.Is(err, expectedError) {
			t.Errorf("RetrievePeople() error = %v, want = %v", err, expectedError)
		}

		if p != nil {
			t.Errorf("RetrievePeople() = %v, want nil", p)
		}
	})
}

func Test_filterPeople(t *testing.T) {
	p := dwp.People{
		{
			ID:        1,
			FirstName: "Maurise",
			LastName:  "Shieldon",
			Email:     "mshieldon0@squidoo.com",
			IPAddress: "192.57.232.111",
			Latitude:  dwp.Coordinate(51.6553959),
			Longitude: dwp.Coordinate(0.0572553),
		},
		{
			ID:        2,
			FirstName: "Bendix",
			LastName:  "Halgarth",
			Email:     "bhalgarth1@timesonline.co.uk",
			IPAddress: "4.185.73.82",
			Latitude:  dwp.Coordinate(-2.9623869),
			Longitude: dwp.Coordinate(104.7399789),
		},
	}

	expectedPeople := p[:1]

	if actualPeople := filterPeople(p, 50, haversine.Coord{Lat: 51.514248, Lon: -0.093145}); !reflect.DeepEqual(actualPeople, expectedPeople) {
		t.Errorf("filterPeople() = %v, want %v", actualPeople, expectedPeople)
	}
}
