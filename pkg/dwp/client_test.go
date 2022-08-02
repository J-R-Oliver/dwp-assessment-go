package dwp

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://test-domain:1020"
	httpClient := http.Client{Timeout: time.Minute}

	want := &client{
		baseURL:    "http://test-domain:1020",
		httpClient: http.Client{Timeout: time.Minute},
	}

	if got := NewClient(baseURL, httpClient); !reflect.DeepEqual(got, want) {
		t.Errorf("NewClient() = %v, want %v", got, want)
	}
}

func Test_client_makeRequest(t *testing.T) {
	t.Run("When called with http.Request with Get HTTP method then server is called with Get HTTP method", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("makeRequest() request method = %s, want: GET", r.Method)
			}
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: *server.Client(),
		}

		r, err := http.NewRequest(http.MethodGet, server.URL+"/test-path", nil)
		if err != nil {
			t.Errorf("makeRequest() error building test http.Request = %v", err)
		}

		c.makeRequest(r, nil) //nolint:errcheck
	})

	t.Run("When called with http.Request with path /test-path then server is called with path /test-path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/test-path" {
				t.Errorf("makeRequest() request path = %s, want: /test-path", r.URL.Path)
			}
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: *server.Client(),
		}

		r, err := http.NewRequest(http.MethodGet, server.URL+"/test-path", nil)
		if err != nil {
			t.Errorf("makeRequest() error building test http.Request = %v", err)
		}

		c.makeRequest(r, nil) //nolint:errcheck
	})

	t.Run("When called with http.Request then server is called with request with Accept-Encoding header", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Accept-Encoding") != "application/json" {
				t.Errorf("makeRequest() request Accept-Encoding header = %s, want: application/json", r.Header.Get("Accept-Encoding"))
			}
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: *server.Client(),
		}

		r, err := http.NewRequest(http.MethodGet, server.URL+"/test-path", nil)
		if err != nil {
			t.Errorf("makeRequest() error building test http.Request = %v", err)
		}

		c.makeRequest(r, nil) //nolint:errcheck
	})

	t.Run("When server responds with JSON HTTP body then body is successful unmarshalled", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"test":"response"}`)) //nolint:errcheck
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: *server.Client(),
		}

		r, err := http.NewRequest(http.MethodGet, server.URL+"/test-path", nil)
		if err != nil {
			t.Errorf("makeRequest() error building test http.Request = %v", err)
		}

		type reponseBody struct {
			Test string
		}
		v := reponseBody{}

		if err := c.makeRequest(r, &v); err != nil {
			t.Errorf("makeRequest() error = %v", err)
		}

		w := reponseBody{"response"}

		if v != w {
			t.Errorf("makeRequest() response body = %s, want: %s", v, w)
		}
	})

	t.Run("When request fails then error is returned", func(t *testing.T) {
		c := &client{
			baseURL:    "http://server-down",
			httpClient: http.Client{},
		}

		r, err := http.NewRequest(http.MethodGet, "http://server-down/test-path", nil)
		if err != nil {
			t.Errorf("makeRequest() error building test http.Request = %v", err)
		}

		var d *net.DNSError
		if err := c.makeRequest(r, nil); !errors.As(err, &d) {
			t.Errorf("makeRequest() error = %v, want = Get \"http://server-down/test-path\": dial tcp: lookup server-down: no such host", err)
		}
	})

	t.Run("When server responds with HTTP error code then error is returned", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: *server.Client(),
		}

		r, err := http.NewRequest(http.MethodGet, server.URL+"/test-path", nil)
		if err != nil {
			t.Errorf("makeRequest() error building test http.Request = %v", err)
		}

		if err := c.makeRequest(r, nil); err.Error() != "request failed with status code 500 and body {}" {
			t.Errorf("makeRequest() error = %v, want = request failed with status code 500 and body {}", err)
		}
	})

	t.Run("When server responds with HTTP error code then error is returned", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// setting a Content-Length longer than body will cause an "unexpected EOF" from io.ReadAll
			w.Header().Add("Content-Length", "50")
			w.Write([]byte("a")) //nolint:errcheck
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: *server.Client(),
		}

		r, err := http.NewRequest(http.MethodGet, server.URL+"/test-path", nil)
		if err != nil {
			t.Errorf("makeRequest() error building test http.Request = %v", err)
		}

		if err := c.makeRequest(r, nil); err.Error() != "unexpected EOF" {
			t.Errorf("makeRequest() error = %v, want = unexpected EOF", err)
		}
	})

	t.Run("When interface{} argument is invalid for body then error is returned", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Not JSON")) //nolint:errcheck
		}))
		defer server.Close()

		c := &client{
			baseURL:    server.URL,
			httpClient: *server.Client(),
		}

		r, err := http.NewRequest(http.MethodGet, server.URL+"/test-path", nil)
		if err != nil {
			t.Errorf("makeRequest() error building test http.Request = %v", err)
		}

		v := struct {
			key string
		}{}
		var e *json.SyntaxError
		if err := c.makeRequest(r, &v); !errors.As(err, &e) {
			t.Errorf("makeRequest() error = %v, want = invalid character 'N' looking for beginning of value", err)
		}
	})
}
