//go:build component

package test

import (
	"net/http"
	"testing"
)

func Test_GetHealth_204(t *testing.T) {
	r, err := HTTPClient.Get(baseURL + "/health")
	if err != nil {
		t.Errorf("GET /health error executing request = %v", err)
		return
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusNoContent {
		t.Errorf("GET /health HTTP status code = %v, want %v", r.StatusCode, http.StatusNoContent)
	}
}
