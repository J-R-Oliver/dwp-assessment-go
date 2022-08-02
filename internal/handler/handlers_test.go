package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/J-R-Oliver/dwp-assessment-go/internal/configuration"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/dwp"
	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
)

const london = "London"

var mockRetrievePeople func() (dwp.People, error)
var mockRetrievePeopleByCity func(city string, distance int) (dwp.People, error)

type mockService struct{}

func (m mockService) RetrievePeople(ctx context.Context) (dwp.People, error) {
	return mockRetrievePeople()
}

func (m mockService) RetrievePeopleByCity(ctx context.Context, city string, distance int) (dwp.People, error) {
	return mockRetrievePeopleByCity(city, distance)
}

func TestHandlers_GetPeople(t *testing.T) {
	t.Run("Given a valid request when there are no errors then people are returned", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/people", nil)

		mockRetrievePeople = func() (dwp.People, error) {
			p := dwp.People{
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
			return p, nil
		}

		h := Handlers{
			Service:         mockService{},
			DefaultDistance: 0,
			Cities:          nil,
			Logger:          nil,
		}
		h.GetPeople(w, r)

		resp := w.Result()

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GetPeople() = %v, want %v", resp.StatusCode, http.StatusOK)
		}

		if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
			t.Errorf("GetPeople() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
		}

		b, _ := io.ReadAll(resp.Body)

		f, err := os.ReadFile("testdata/people.txt")
		if err != nil {
			t.Errorf("GetPeople() error reading testdata/people.txt = %v", err)
			return
		}

		if !bytes.Equal(b, f) {
			t.Errorf("GetPeople() = %v, want %v", b, f)
		}
	})

	t.Run("Given a request with HTTP method post then method not allowed response", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/people", nil)

		h := Handlers{
			Service:         mockService{},
			DefaultDistance: 0,
			Cities:          nil,
			Logger:          nil,
		}
		h.GetPeople(w, r)

		resp := w.Result()

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("GetPeople() = %v, want %v", resp.StatusCode, http.StatusMethodNotAllowed)
		}

		if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
			t.Errorf("GetPeople() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
		}

		b, _ := io.ReadAll(resp.Body)
		body := string(b)

		expectedBody := `"status":405,"message":"Method Not Allowed","path":"/api/people"`

		if !strings.Contains(body, expectedBody) {
			t.Errorf("GetPeople() = %v, want %v", body, expectedBody)
		}
	})

	t.Run("Given a valid request when there is RetrievePeople error then internal server error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/people", nil)

		mockRetrievePeople = func() (dwp.People, error) {
			return nil, errors.New("test error")
		}

		h := Handlers{
			Service:         mockService{},
			DefaultDistance: 0,
			Cities:          nil,
			Logger:          logging.New(logging.Info),
		}
		h.GetPeople(w, r)

		resp := w.Result()

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("GetPeople() = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
		}

		if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
			t.Errorf("GetPeople() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
		}

		b, _ := io.ReadAll(resp.Body)
		body := string(b)

		expectedBody := `"status":500,"message":"Internal Server Error","path":"/api/people"`

		if !strings.Contains(body, expectedBody) {
			t.Errorf("GetPeople() = %v, want %v", body, expectedBody)
		}
	})
}

func TestHandlers_GetPeopleByCity(t *testing.T) {
	t.Run("Given a valid request when there are no errors then people are returned", func(t *testing.T) { //nolint:dupl
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/people/london", nil)

		mockRetrievePeopleByCity = func(city string, distance int) (dwp.People, error) {
			if city != london {
				t.Errorf("GetPeopleByCity() = %v, London", city)
			}

			if distance != 50 {
				t.Errorf("GetPeopleByCity() = %v, 50", distance)
			}

			p := dwp.People{
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
			return p, nil
		}

		h := Handlers{
			Service:         mockService{},
			DefaultDistance: 50,
			Cities:          map[string]configuration.City{london: {}},
			Logger:          nil,
		}
		h.GetPeopleByCity("/api/people/")(w, r)

		resp := w.Result()

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.StatusCode, http.StatusOK)
		}

		if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
		}

		b, _ := io.ReadAll(resp.Body)

		f, err := os.ReadFile("testdata/people.txt")
		if err != nil {
			t.Errorf("GetPeopleByCity() error reading testdata/people.txt = %v", err)
			return
		}

		if !bytes.Equal(b, f) {
			t.Errorf("GetPeopleByCity() = %v, want %v", b, f)
		}
	})

	t.Run("Given a valid request with distance query when there are no errors then people are returned", func(t *testing.T) { //nolint:dupl
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/people/london?distance=5", nil)

		mockRetrievePeopleByCity = func(city string, distance int) (dwp.People, error) {
			if city != london {
				t.Errorf("GetPeopleByCity() = %v, London", city)
			}

			if distance != 5 {
				t.Errorf("GetPeopleByCity() = %v, 5", distance)
			}

			p := dwp.People{
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
			return p, nil
		}

		h := Handlers{
			Service:         mockService{},
			DefaultDistance: 0,
			Cities:          map[string]configuration.City{london: {}},
			Logger:          nil,
		}
		h.GetPeopleByCity("/api/people/")(w, r)

		resp := w.Result()

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.StatusCode, http.StatusOK)
		}

		if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
		}

		b, _ := io.ReadAll(resp.Body)

		f, err := os.ReadFile("testdata/people.txt")
		if err != nil {
			t.Errorf("GetPeopleByCity() error reading testdata/people.txt = %v", err)
			return
		}

		if !bytes.Equal(b, f) {
			t.Errorf("GetPeopleByCity() = %v, want %v", b, f)
		}
	})

	t.Run("Given a request with HTTP method post then method not allowed response", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/people/london", nil)

		h := Handlers{
			Service:         mockService{},
			DefaultDistance: 0,
			Cities:          nil,
			Logger:          nil,
		}
		h.GetPeopleByCity("/api/people/")(w, r)

		resp := w.Result()

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.StatusCode, http.StatusMethodNotAllowed)
		}

		if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
		}

		b, _ := io.ReadAll(resp.Body)
		body := string(b)

		expectedBody := `"status":405,"message":"Method Not Allowed","path":"/api/people/london"`

		if !strings.Contains(body, expectedBody) {
			t.Errorf("GetPeopleByCity() = %v, want %v", body, expectedBody)
		}
	})

	t.Run("Given a request with an invalid distance query then bad request response", func(t *testing.T) { //nolint:dupl
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/people/london?distance=not-an-int", nil)

		h := Handlers{
			Service:         mockService{},
			DefaultDistance: 0,
			Cities:          nil,
			Logger:          logging.New(logging.Info),
		}
		h.GetPeopleByCity("/api/people/")(w, r)

		resp := w.Result()

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
		}

		if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
		}

		b, _ := io.ReadAll(resp.Body)
		body := string(b)

		expectedBody := `"status":400,"message":"Invalid distance query - not-an-int is not an integer","path":"/api/people/london"`

		if !strings.Contains(body, expectedBody) {
			t.Errorf("GetPeopleByCity() = %v, want %v", body, expectedBody)
		}
	})

	t.Run("Given a request with an unknown city then not found response", func(t *testing.T) { //nolint:dupl
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/people/timbuctoo", nil)

		h := Handlers{
			Service:         mockService{},
			DefaultDistance: 0,
			Cities:          nil,
			Logger:          logging.New(logging.Info),
		}
		h.GetPeopleByCity("/api/people/")(w, r)

		resp := w.Result()

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.StatusCode, http.StatusNotFound)
		}

		if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
		}

		b, _ := io.ReadAll(resp.Body)
		body := string(b)

		expectedBody := `"status":404,"message":"City Not Found","path":"/api/people/timbuctoo"`

		if !strings.Contains(body, expectedBody) {
			t.Errorf("GetPeopleByCity() = %v, want %v", body, expectedBody)
		}
	})

	t.Run("Given a valid request when there is RetrievePeopleByCity error then internal server error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/people/london", nil)

		mockRetrievePeopleByCity = func(city string, distance int) (dwp.People, error) {
			if city != london {
				t.Errorf("GetPeopleByCity() = %v, London", city)
			}

			if distance != 50 {
				t.Errorf("GetPeopleByCity() = %v, 50", distance)
			}

			return nil, errors.New("test error")
		}

		h := Handlers{
			Service:         mockService{},
			DefaultDistance: 50,
			Cities:          map[string]configuration.City{london: {}},
			Logger:          logging.New(logging.Info),
		}
		h.GetPeopleByCity("/api/people/")(w, r)

		resp := w.Result()

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
		}

		if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
			t.Errorf("GetPeopleByCity() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
		}

		b, _ := io.ReadAll(resp.Body)
		body := string(b)

		expectedBody := `"status":500,"message":"Internal Server Error","path":"/api/people/london"`

		if !strings.Contains(body, expectedBody) {
			t.Errorf("GetPeopleByCity() = %v, want %v", body, expectedBody)
		}
	})
}

func TestHandlers_Health(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	h := Handlers{
		Service:         nil,
		DefaultDistance: 0,
		Cities:          nil,
		Logger:          nil,
	}

	h.Health(w, r)

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Health() = %v, want %v", w.Code, http.StatusNoContent)
	}
}
