package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestRouter() http.Handler {
	logger := log.New(io.Discard, "", 0)

	app := NewApplication(logger)

	return app.Router()
}

func TestHealtHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router := newTestRouter()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}

	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf(
			"expected Content-Type application/json, got %q",
			contentType,
		)
	}

	var body healthResponse

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", body.Status)
	}
}

func TestHealthHanderRejectsUnsupportedMethods(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/health", nil)
	recorder := httptest.NewRecorder()

	router := newTestRouter()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			response.StatusCode,
		)
	}

	allow := response.Header.Get("Allow")

	if allow != http.MethodGet {
		t.Errorf(
			"expected Allow header %q, got %q",
			http.MethodGet,
			allow,
		)
	}
}

func TestRouterReturnsJSONNotFound(t *testing.T) {

	router := newTestRouter()

	request := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.StatusCode,
		)
	}

	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf(
			"expected Content-Type application/json, got %q",
			contentType,
		)
	}

	var body struct {
		Error string `json:"error"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Error != "route not found" {
		t.Errorf(
			"expected error %q, got %q",
			"route not found",
			body.Error,
		)
	}
}

func TestRequestLoggerLogsCompletedRequest(t *testing.T) {

	var logOutput bytes.Buffer

	logger := log.New(&logOutput, "", 0)

	app := NewApplication(logger)
	router := app.Router()

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	expectedLog := "request method=GET path=/health status=200 duration="

	if !strings.Contains(logOutput.String(), expectedLog) {
		t.Errorf(
			"expected log output to contain %q, got %q",
			expectedLog,
			logOutput.String(),
		)
	}
}

func TestRecoverPanicReturnsInternalServerError(t *testing.T) {
	var logOutput bytes.Buffer

	logger := log.New(&logOutput, "", 0)

	app := NewApplication(logger)

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected test failure")
	},
	)

	handler := app.requestLogger(
		app.recoverPanic(panicHandler),
	)

	request := httptest.NewRequest(http.MethodGet, "/panic-test", nil)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			response.StatusCode,
		)
	}

	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf(
			"expected Content-Type application/json, got %q",
			contentType,
		)
	}

	var body errorResponse

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Error != "internal server error" {
		t.Errorf(
			"expected error %q, got %q",
			"internal server error",
			body.Error,
		)
	}

	if !strings.Contains(logOutput.String(), "panic recovered: unexpected test failure") {
		t.Errorf(
			"expected log output to contain panic recovery message, got %q",
			logOutput.String(),
		)
	}

	if !strings.Contains(logOutput.String(), "status=500") {
		t.Errorf(
			"expected log output to contain status=500, got %q",
			logOutput.String(),
		)
	}
}
