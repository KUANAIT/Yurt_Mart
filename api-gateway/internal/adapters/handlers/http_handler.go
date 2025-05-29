package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"api-gateway/internal/core/ports"
	"api-gateway/pkg/logger"
)

type HTTPHandler struct {
	service ports.UserServicePort
	logger  logger.Logger
}

func NewHTTPHandler(service ports.UserServicePort, logger logger.Logger) *HTTPHandler {
	return &HTTPHandler{
		service: service,
		logger:  logger,
	}
}

func (h *HTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if request.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}
	if request.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Validate email format
	if !strings.Contains(request.Email, "@") {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password length
	if len(request.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	userID, err := h.service.RegisterUser(r.Context(), request.Email, request.Password, request.Name)
	if err != nil {
		h.logger.Error("Registration failed: %v", err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	response := struct {
		UserID string `json:"user_id"`
	}{
		UserID: userID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *HTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
