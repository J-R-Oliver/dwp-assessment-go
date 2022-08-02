package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
)

var mockNext func(w http.ResponseWriter, r *http.Request)

type mock struct{}

func (m mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mockNext(w, r)
}

func TestPanicHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	mockNext = func(w http.ResponseWriter, r *http.Request) {
		if w != response {
			t.Errorf("LogRequestHandler() ResponseWriter = %v, want %v", w, response)
		}

		if r != request {
			t.Errorf("LogRequestHandler() Requetss = %v, want %v", r, request)
		}

		panic("test panic")
	}

	mockInternalServerErrorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		if err.Error() != "internal server error" {
			t.Errorf("PanicHandler() err = %v, want internal server error", err)
		}
	}

	h := PanicHandler(mock{}, mockInternalServerErrorHandler)

	h.ServeHTTP(response, request)
}

func TestLogRequestHandler(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	mockNext = func(w http.ResponseWriter, r *http.Request) {
		if w != response {
			t.Errorf("LogRequestHandler() ResponseWriter = %v, want %v", w, response)
		}

		if r != request {
			t.Errorf("LogRequestHandler() Requetss = %v, want %v", r, request)
		}
	}

	l := logging.New(logging.Info)
	h := LogRequestHandler(mock{}, l)

	h.ServeHTTP(response, request)

	w.Close()

	out, _ := ioutil.ReadAll(r)
	s := string(out)

	if !strings.Contains(s, "HTTP/1.1 GET /health") {
		t.Errorf("LogRequestHandler() = %v, want HTTP/1.1 GET /health", s)
	}
}
