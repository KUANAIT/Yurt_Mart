package test

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/internal/adapters/handlers"
	"api-gateway/internal/core/services"
	"api-gateway/pkg/logger"
)

var (
	mockUserService    *MockUserService
	mockCache          *MockCache
	mockEventPublisher *MockEventPublisher
	handler            *handlers.HTTPHandler
	server             *httptest.Server
)

func setupTest(t *testing.T) {
	mockUserService = NewMockUserService()
	mockCache = NewMockCache()
	mockEventPublisher = NewMockEventPublisher()

	gatewayService := services.NewGatewayService(
		mockUserService,
		mockCache,
		mockEventPublisher,
	)

	log := logger.NewLogger()
	handler = handlers.NewHTTPHandler(gatewayService, log)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", handler.RegisterUser)
	mux.HandleFunc("/api/v1/user", handler.GetUser)
	mux.HandleFunc("/health", handler.HealthCheck)
	server = httptest.NewServer(mux)
}

func teardownTest() {
	if server != nil {
		server.Close()
	}
}

func TestHealthCheck(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make health check request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestRegisterUser(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	userData := map[string]string{
		"email":    "testno@example.com",
		"password": "testpassword555",
		"name":     "Test User",
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Failed to marshal user data: %v", err)
	}

	resp, err := http.Post(
		server.URL+"/api/v1/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("Failed to make register request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.UserID == "" {
		t.Error("Expected non-empty user ID in response")
	}

	t.Logf("Registered user ID: %s", response.UserID)

	testGetUser(t, response.UserID)
}

func testGetUser(t *testing.T, userID string) {
	resp, err := http.Get(server.URL + "/api/v1/user?id=" + userID)
	if err != nil {
		t.Fatalf("Failed to make get user request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var user struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if user.ID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, user.ID)
	}
}

func TestInvalidRegistration(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	invalidData := map[string]string{
		"email": "invalid-email",
		"name":  "Test User",
	}

	jsonData, err := json.Marshal(invalidData)
	if err != nil {
		t.Fatalf("Failed to marshal invalid data: %v", err)
	}

	resp, err := http.Post(
		server.URL+"/api/v1/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("Failed to make register request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetNonExistentUser(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	nonExistentID := "non-existent-id"

	resp, err := http.Get(server.URL + "/api/v1/user?id=" + nonExistentID)
	if err != nil {
		t.Fatalf("Failed to make get user request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}
