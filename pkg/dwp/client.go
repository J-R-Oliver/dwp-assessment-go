package dwp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client interface {
	RetrievePeople(ctx context.Context) (People, error)
	RetrievePeopleByCity(ctx context.Context, city string) (People, error)
}

type client struct {
	baseURL    string
	httpClient http.Client
}

func NewClient(baseURL string, httpClient http.Client) Client {
	return &client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c client) makeRequest(r *http.Request, v interface{}) error {
	r.Header.Set("Accept-Encoding", "application/json")

	response, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d and body %s", response.StatusCode, response.Body)
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
