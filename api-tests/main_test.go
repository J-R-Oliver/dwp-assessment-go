package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"
)

var baseURL string
var HTTPClient http.Client

func TestMain(m *testing.M) {
	fmt.Println("Running test setup...")

	if e, b := os.LookupEnv("BASE_URL"); b {
		baseURL = e
	} else {
		baseURL = "http://localhost:8080"
	}

	fmt.Printf("Setting base URL: %s\n", baseURL)

	HTTPClient = http.Client{}

	os.Exit(m.Run())
}

func test404(t *testing.T, path string, errorMessage string) {
	t.Helper()

	r, err := HTTPClient.Get(baseURL + path)
	if err != nil {
		t.Errorf("GET %s error executing request = %v", path, err)
		return
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusNotFound {
		t.Errorf("GET %s HTTP status code = %v, want %v", path, r.StatusCode, http.StatusNotFound)
	}

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("GET %s HTTP Content-Type = %v, want application/json", path, r.Header.Get("Content-Type"))
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("GET %s error reading body = %v", path, err)
		return
	}

	var j map[string]interface{}
	if err = json.Unmarshal(b, &j); err != nil {
		t.Errorf("GET %s error unmarshalling body = %v", path, err)
		return
	}

	if m, _ := regexp.MatchString("\\d{4}-[01]\\d-[0-3]\\dT[0-2]\\d:[0-5]\\d:[0-5]\\d\\.\\d+([+-][0-2]\\d:[0-5]\\d|Z)", j["timestamp"].(string)); !m {
		t.Errorf("GET %s timestamp = %v, want ISO timestamp", path, j["timestamp"])
	}

	if j["status"].(float64) != http.StatusNotFound {
		t.Errorf("GET %s status = %v, want %v", path, j["status"], http.StatusNotFound)
	}

	if j["message"].(string) != errorMessage {
		t.Errorf("GET %s message = %v, want %s", path, j["message"], errorMessage)
	}

	if j["path"].(string) != path {
		t.Errorf("GET %s path = %v, want %s", path, j["path"], path)
	}
}

func test405(t *testing.T, method string, path string) {
	t.Helper()

	req, err := http.NewRequest(method, baseURL+path, nil)
	if err != nil {
		t.Errorf("GET %s error creating request = %v", path, err)
		return
	}

	r, err := HTTPClient.Do(req)
	if err != nil {
		t.Errorf("GET %s error executing request = %v", path, err)
		return
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("GET %s HTTP status code = %v, want %v", path, r.StatusCode, http.StatusMethodNotAllowed)
	}

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("GET /api/people HTTP Content-Type = %v, want application/json", r.Header.Get("Content-Type"))
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("GET %s error reading body = %v", path, err)
		return
	}

	if method == http.MethodHead {
		return
	}

	var j map[string]interface{}
	if err = json.Unmarshal(b, &j); err != nil {
		t.Errorf("GET %s error unmarshalling body = %v", path, err)
		return
	}

	if m, _ := regexp.MatchString("\\d{4}-[01]\\d-[0-3]\\dT[0-2]\\d:[0-5]\\d:[0-5]\\d\\.\\d+([+-][0-2]\\d:[0-5]\\d|Z)", j["timestamp"].(string)); !m {
		t.Errorf("GET %s timestamp = %v, want ISO timestamp", path, j["timestamp"])
	}

	if j["status"].(float64) != http.StatusMethodNotAllowed {
		t.Errorf("GET %s status = %v, want %v", path, j["status"], http.StatusMethodNotAllowed)
	}

	if j["message"].(string) != "Method Not Allowed" {
		t.Errorf("GET %s message = %v, want Method Not Allowed", path, j["message"])
	}

	if j["path"].(string) != "/api/people/london" {
		t.Errorf("GET %s path = %v, want /api/people/london", path, j["path"])
	}
}
