package test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

func Test_GetApiPeople_200(t *testing.T) {
	r, err := HTTPClient.Get(baseURL + "/api/people")
	if err != nil {
		t.Errorf("GET /api/people error executing request = %v", err)
		return
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Errorf("GET /api/people HTTP status code = %v, want %v", r.StatusCode, http.StatusOK)
	}

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("GET /api/people HTTP Content-Type = %v, want application/json", r.Header.Get("Content-Type"))
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("GET /api/people error reading body = %v", err)
		return
	}

	f, err := os.ReadFile("testdata/all-people.txt")
	if err != nil {
		t.Errorf("GET /api/people error reading testdata/all-people.txt = %v", err)
		return
	}

	if !bytes.Equal(b, f) {
		t.Errorf("GET /api/people resopnse body %s, want %s", b, f)
	}
}

func Test_GetApiPeople_405(t *testing.T) {
	tests := []string{
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
		http.MethodPatch,
	}

	for _, method := range tests {
		t.Run(method, func(t *testing.T) {
			test405(t, method, "/api/people")
		})
	}
}
