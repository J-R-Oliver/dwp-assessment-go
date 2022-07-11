package dwp

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func Test_coordinate_UnmarshalJSON(t *testing.T) {
	var b bytes.Buffer

	binary.Write(&b, binary.BigEndian, 1.23) //nolint:errcheck

	type args struct {
		b []byte
	}

	tests := []struct {
		name    string
		c       coordinate
		args    args
		wantErr bool
		want    float64
	}{
		{
			"When passed float as JSON number as byte array then parses float",
			coordinate(0.0),
			args{[]byte("1.23")},
			false,
			1.23,
		},
		{
			"When passed float as JSON string as byte array then parses float",
			coordinate(0.0),
			args{[]byte("\"1.23\"")},
			false,
			1.23,
		},
		{
			"When passed string as JSON string as byte array then return error",
			coordinate(0.0),
			args{[]byte("Not a number")},
			true,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.c != coordinate(tt.want) {
				t.Errorf("UnmarshalJSON() = %v, want %v", tt.c, tt.want)
			}
		})
	}
}

func Test_client_RetrievePeople(t *testing.T) {
	t.Run("When request to RetrievePeople is successful then parse People are returned", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("[{\"id\": 1, \"first_name\": \"Maurise\", \"last_name\": \"Shieldon\", \"email\": \"mshieldon0@squidoo.com\", \"ip_address\": \"192.57.232.111\", \"latitude\": 34.003135, \"longitude\": -117.7228641 }, { \"id\": 2, \"first_name\": \"Bendix\", \"last_name\": \"Halgarth\", \"email\": \"bhalgarth1@timesonline.co.uk\", \"ip_address\": \"4.185.73.82\", \"latitude\": -2.9623869, \"longitude\": 104.7399789}]")) //nolint:errcheck
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: http.Client{},
		}

		p, err := c.RetrievePeople(context.Background())
		if err != nil {
			t.Errorf("RetrievePeople() error = %v", err)
		}

		e := People{
			{
				ID:        1,
				FirstName: "Maurise",
				LastName:  "Shieldon",
				Email:     "mshieldon0@squidoo.com",
				IPAddress: "192.57.232.111",
				Latitude:  coordinate(34.003135),
				Longitude: coordinate(-117.7228641),
			},
			{
				ID:        2,
				FirstName: "Bendix",
				LastName:  "Halgarth",
				Email:     "bhalgarth1@timesonline.co.uk",
				IPAddress: "4.185.73.82",
				Latitude:  coordinate(-2.9623869),
				Longitude: coordinate(104.7399789),
			},
		}

		if !reflect.DeepEqual(p, e) {
			t.Errorf("RetrievePeople() = %v, want %v", p, e)
		}
	})

	t.Run("When invalid client is instantiated with invalid baseURL then error is return", func(t *testing.T) {
		c := &client{
			baseURL:    string([]byte{0x7f}), // unicode character DELETE causes http.NewRequest err
			httpClient: http.Client{},
		}

		var e *url.Error
		if _, err := c.RetrievePeople(context.Background()); !errors.As(err, &e) {
			t.Errorf("RetrievePeople() error = %v, want = RetrievePeople: failed creating http request: invalid control character in URL", err)
		}
	})

	t.Run("When server returns HTTP status 500 then error is returned", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: http.Client{},
		}

		if _, err := c.RetrievePeople(context.Background()); err.Error() != "RetrievePeople: failed executing http request: request failed with status code 500 and body {}" {
			t.Errorf("RetrievePeople() error = %v, want = RetrievePeople: failed executing http request: request failed with status code 500 and body {}", err)
		}
	})
}

func Test_client_RetrievePeopleByCity(t *testing.T) {
	t.Run("When city is passed then correct URL path is created", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/city/london/users" {
				t.Errorf("RetrievePeopleByCity() %v, want = /city/london/users", r.URL.Path)
			}
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: http.Client{},
		}

		c.RetrievePeopleByCity(context.Background(), "london") //nolint:errcheck
	})

	t.Run("When request to RetrievePeopleByCity is successful then parse People are returned", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("[{\"id\": 1, \"first_name\": \"Maurise\", \"last_name\": \"Shieldon\", \"email\": \"mshieldon0@squidoo.com\", \"ip_address\": \"192.57.232.111\", \"latitude\": 34.003135, \"longitude\": -117.7228641 }, { \"id\": 2, \"first_name\": \"Bendix\", \"last_name\": \"Halgarth\", \"email\": \"bhalgarth1@timesonline.co.uk\", \"ip_address\": \"4.185.73.82\", \"latitude\": -2.9623869, \"longitude\": 104.7399789}]")) //nolint:errcheck
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: http.Client{},
		}

		p, err := c.RetrievePeopleByCity(context.Background(), "london")
		if err != nil {
			t.Errorf("RetrievePeopleByCity() error = %v", err)
		}

		e := People{
			{
				ID:        1,
				FirstName: "Maurise",
				LastName:  "Shieldon",
				Email:     "mshieldon0@squidoo.com",
				IPAddress: "192.57.232.111",
				Latitude:  coordinate(34.003135),
				Longitude: coordinate(-117.7228641),
			},
			{
				ID:        2,
				FirstName: "Bendix",
				LastName:  "Halgarth",
				Email:     "bhalgarth1@timesonline.co.uk",
				IPAddress: "4.185.73.82",
				Latitude:  coordinate(-2.9623869),
				Longitude: coordinate(104.7399789),
			},
		}

		if !reflect.DeepEqual(p, e) {
			t.Errorf("RetrievePeopleByCity() = %v, want %v", p, e)
		}
	})

	t.Run("When invalid client is instantiated with invalid baseURL then error is return", func(t *testing.T) {
		c := &client{
			baseURL:    string([]byte{0x7f}), // unicode character DELETE causes http.NewRequest err
			httpClient: http.Client{},
		}

		var e *url.Error
		if _, err := c.RetrievePeopleByCity(context.Background(), "london"); !errors.As(err, &e) {
			t.Errorf("RetrievePeopleByCity() error = %v, want = RetrievePeopleByCity: failed creating http request: invalid control character in URL", err)
		}
	})

	t.Run("When server returns HTTP status 500 then error is returned", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: http.Client{},
		}

		if _, err := c.RetrievePeopleByCity(context.Background(), "london"); err.Error() != "RetrievePeopleByCity: failed executing http request: request failed with status code 500 and body {}" {
			t.Errorf("RetrievePeopleByCity() error = %v, want = RetrievePeopleByCity: failed executing http request: request failed with status code 500 and body {}", err)
		}
	})
}
