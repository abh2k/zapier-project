package handlers_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "zapier-project/internal/http"
	"zapier-project/internal/http/handlers"

	"zapier-project/internal/deployments"
)

func setupRouter() http.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	store := deployments.NewStore(deployments.SeedData())
	deploymentsHandler := handlers.NewDeploymentsHandler(store, logger)
	return httpserver.NewRouter(deploymentsHandler)
}

func TestListDeploymentsFiltersByServiceAndStatus(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/deployments?service=billing-api&status=failed", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response struct {
		Data []deployments.Deployment `json:"data"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Data) == 0 {
		t.Fatalf("expected filtered results, got none")
	}

	for _, d := range response.Data {
		if d.Service != "billing-api" || d.Status != "failed" {
			t.Fatalf("unexpected deployment in results: %+v", d)
		}
	}
}

func TestListDeploymentsInvalidStatus(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/deployments?status=in_progress", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	var response struct {
		Error handlers.APIError `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error.Code != "invalid_status" {
		t.Fatalf("expected error code invalid_status, got %q", response.Error.Code)
	}
}

func TestGetDeploymentByIDNotFound(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/deployments/deploy_missing", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}

	var response struct {
		Error handlers.APIError `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error.Code != "not_found" {
		t.Fatalf("expected error code not_found, got %q", response.Error.Code)
	}
}

func TestGetDeploymentByIDSuccess(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/deployments/deploy_001", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response struct {
		Data deployments.Deployment `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Data.ID != "deploy_001" {
		t.Fatalf("expected deployment id deploy_001, got %q", response.Data.ID)
	}
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Fatalf("expected health status ok, got %q", response["status"])
	}
}
