package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"api-gateway/internal/adapters/handlers"
	"api-gateway/pkg/logger"
)

type Server struct {
	handler *handlers.HTTPHandler
	logger  logger.Logger
	port    int
}

func NewServer(handler *handlers.HTTPHandler, logger logger.Logger, port int) *Server {
	return &Server{
		handler: handler,
		logger:  logger,
		port:    port,
	}
}

func (s *Server) Start() error {
	router := http.NewServeMux()

	// API routes
	router.HandleFunc("/api/v1/register", s.handler.RegisterUser)
	router.HandleFunc("/api/v1/user", s.handler.GetUser)
	router.HandleFunc("/health", s.handler.HealthCheck)

	// Middleware chain
	handler := s.loggingMiddleware(router)
	handler = s.recoveryMiddleware(handler)

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(s.port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logger.Info("Starting HTTP server on port %d", s.port)
	return server.ListenAndServe()
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Info("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Error("Recovered from panic: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	router := http.NewServeMux()

	// API routes
	router.HandleFunc("/api/v1/register", s.handler.RegisterUser)
	router.HandleFunc("/api/v1/user", s.handler.GetUser)
	router.HandleFunc("/health", s.handler.HealthCheck)

	// Middleware chain
	handler := s.loggingMiddleware(router)
	handler = s.recoveryMiddleware(handler)

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(s.port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return server.Shutdown(ctx)
}
