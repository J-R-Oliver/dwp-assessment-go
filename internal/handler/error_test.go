package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
)

func TestHandlers_notFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/path", nil)

	h := Handlers{
		Service:         nil,
		DefaultDistance: 0,
		Cities:          nil,
		Logger:          logging.New(logging.Info),
	}

	h.NotFound(w, r)

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("NotFound() = %v, want %v", resp.StatusCode, http.StatusNotFound)
	}

	if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
		t.Errorf("NotFound() = %v, want application/json", resp.Header.Get("Content-Type"))
	}

	b, _ := io.ReadAll(resp.Body)
	body := string(b)

	expectedBody := `"status":404,"message":"Route Not Found","path":"/path"`

	if !strings.Contains(body, expectedBody) {
		t.Errorf("NotFound() = %v, want %v", body, expectedBody)
	}
}

func TestHandlers_InternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/path", nil)

	h := Handlers{
		Service:         nil,
		DefaultDistance: 0,
		Cities:          nil,
		Logger:          logging.New(logging.Info),
	}

	h.InternalServerError(w, r, nil)

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("InternalServerError() = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}

	if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
		t.Errorf("InternalServerError() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
	}

	b, _ := io.ReadAll(resp.Body)
	body := string(b)

	expectedBody := `"status":500,"message":"Internal Server Error","path":"/path"`

	if !strings.Contains(body, expectedBody) {
		t.Errorf("NotFound() = %v, want %v", body, expectedBody)
	}
}

func TestHandlers_badRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/path", nil)

	h := Handlers{
		Service:         nil,
		DefaultDistance: 0,
		Cities:          nil,
		Logger:          logging.New(logging.Info),
	}

	h.badRequest(w, r, "test message")

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("badRequest() = %v, want %v", resp.StatusCode, http.StatusBadRequest)
	}

	if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
		t.Errorf("badRequest() = %v, want application/json", resp.Header.Get("Content-Type"))
	}

	b, _ := io.ReadAll(resp.Body)
	body := string(b)

	expectedBody := `"status":400,"message":"test message","path":"/path"`

	if !strings.Contains(body, expectedBody) {
		t.Errorf("badRequest() = %v, want %v", body, expectedBody)
	}
}

func TestHandlers_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/path", nil)

	h := Handlers{
		Service:         nil,
		DefaultDistance: 0,
		Cities:          nil,
		Logger:          logging.New(logging.Info),
	}

	h.notFound(w, r, "test message")

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("notFound() = %v, want %v", resp.StatusCode, http.StatusNotFound)
	}

	if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
		t.Errorf("notFound() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
	}

	b, _ := io.ReadAll(resp.Body)
	body := string(b)

	expectedBody := `"status":404,"message":"test message","path":"/path"`

	if !strings.Contains(body, expectedBody) {
		t.Errorf("notFound() = %v, want %v", body, expectedBody)
	}
}

func TestHandlers_methodNotAllow(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/path", nil)

	h := Handlers{
		Service:         nil,
		DefaultDistance: 0,
		Cities:          nil,
		Logger:          logging.New(logging.Info),
	}

	h.methodNotAllow(w, r)

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("methodNotAllow() = %v, want %v", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
		t.Errorf("methodNotAllow() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
	}

	b, _ := io.ReadAll(resp.Body)
	body := string(b)

	expectedBody := `"status":405,"message":"Method Not Allowed","path":"/path"`

	if !strings.Contains(body, expectedBody) {
		t.Errorf("methodNotAllow() = %v, want %v", body, expectedBody)
	}
}

func TestHandlers_errorHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/path", nil)

	h := Handlers{
		Service:         nil,
		DefaultDistance: 0,
		Cities:          nil,
		Logger:          logging.New(logging.Info),
	}

	h.errorHandler(w, r, 405, "Method Not Allowed")

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("errorHandler() = %v, want %v", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	if resp.Header.Get("Content-Type") != ContentTypeApplicationJSON {
		t.Errorf("errorHandler() = %v, want %v", resp.Header.Get("Content-Type"), ContentTypeApplicationJSON)
	}

	b, _ := io.ReadAll(resp.Body)
	body := string(b)

	expectedBody := `"status":405,"message":"Method Not Allowed","path":"/path"`

	if !strings.Contains(body, expectedBody) {
		t.Errorf("errorHandler() = %v, want %v", body, expectedBody)
	}
}
