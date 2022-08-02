package dwp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

const bitSize = 64

type Coordinate float64

func (c *Coordinate) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		b = b[1:]
		b = b[:len(b)-1]
	}

	float, err := strconv.ParseFloat(string(b), bitSize)
	if err != nil {
		return fmt.Errorf("coordinate.UnmarshalJSON: failed to unmarshal %s: %w", string(b), err)
	}

	*c = Coordinate(float)

	return nil
}

type Person struct {
	ID        int
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string
	IPAddress string `json:"ip_address"`
	Latitude  Coordinate
	Longitude Coordinate
}

type People []Person

// RetrievePeople returns all people from the '/users' endpoint. If any error is returned then People will be nil.
func (c client) RetrievePeople(ctx context.Context) (People, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/users", nil)
	if err != nil {
		return nil, fmt.Errorf("RetrievePeople: failed creating http request: %w", err)
	}

	people := People{}

	err = c.makeRequest(request, &people)
	if err != nil {
		return nil, fmt.Errorf("RetrievePeople: failed executing http request: %w", err)
	}

	return people, nil
}

// RetrievePeopleByCity returns all people from the '/city/{city}/users' endpoint. If any error is returned then People
// will be nil.
func (c client) RetrievePeopleByCity(ctx context.Context, city string) (People, error) {
	path := fmt.Sprintf("/city/%s/users", city)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("RetrievePeopleByCity: failed creating http request: %w", err)
	}

	people := People{}

	err = c.makeRequest(request, &people)
	if err != nil {
		return nil, fmt.Errorf("RetrievePeopleByCity: failed executing http request: %w", err)
	}

	return people, nil
}
