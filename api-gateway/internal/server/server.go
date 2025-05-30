package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/internal/auth"
	"api-gateway/internal/config/order"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	router        *gin.Engine
	config        *order.Config
	orderClient   *proxy.OrderServiceClient
	paymentClient *proxy.PaymentServiceClient
	rateLimiter   *middleware.RateLimiter
	jwtAuth       *auth.JWTAuth
}

func NewServer(config *order.Config) (*Server, error) {
	// Initialize gRPC clients
	orderClient, err := proxy.NewOrderServiceClient(config.Services.OrderServiceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create order service client: %w", err)
	}

	paymentClient, err := proxy.NewPaymentServiceClient(config.Services.PaymentServiceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment service client: %w", err)
	}

	// Initialize middleware
	rateLimiter := middleware.NewRateLimiter(&config.Redis, &config.RateLimiting)
	jwtAuth := auth.NewJWTAuth(&config.Auth)

	server := &Server{
		router:        gin.Default(),
		config:        config,
		orderClient:   orderClient,
		paymentClient: paymentClient,
		rateLimiter:   rateLimiter,
		jwtAuth:       jwtAuth,
	}

	server.setupRoutes()
	return server, nil
}

func (s *Server) setupRoutes() {
	// Create handlers
	orderHandler := handler.NewOrderHandler(s.orderClient, s.paymentClient)
	paymentHandler := handler.NewPaymentHandler(s.paymentClient)

	// Middleware
	s.router.Use(gin.Recovery())
	s.router.Use(s.rateLimiter.Middleware())

	// Health check and metrics
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", s.handleLogin)
			auth.POST("/register", s.handleRegister)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(s.jwtAuth.Middleware())
		{
			// Order routes
			orders := protected.Group("/orders")
			{
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
				orders.GET("/delivery-slots", orderHandler.GetAvailableDeliverySlots)
			}

			// Delivery address routes
			addresses := protected.Group("/addresses")
			{
				addresses.POST("", orderHandler.AddDeliveryAddress)
				addresses.GET("", orderHandler.ListDeliveryAddresses)
			}

			// Payment routes
			payments := protected.Group("/payments")
			{
				payments.POST("", paymentHandler.InitiatePayment)
				payments.POST("/credit-card", paymentHandler.ProcessCreditCardPayment)
				payments.POST("/metamask/initiate", paymentHandler.InitiateMetaMaskPayment)
				payments.POST("/metamask/confirm", paymentHandler.ConfirmMetaMaskPayment)
				payments.GET("/:id", paymentHandler.GetPayment)
				payments.GET("/order/:order_id", paymentHandler.GetPaymentsByOrder)
				payments.GET("/pending", paymentHandler.GetPendingPayments)
			}
		}
	}
}

func (s *Server) handleLogin(c *gin.Context) {
	// TODO: Implement user authentication
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (s *Server) handleRegister(c *gin.Context) {
	// TODO: Implement user registration
	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         ":" + s.config.Server.Port,
		Handler:      s.router,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Server.ShutdownTimeout)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	// Close clients
	s.orderClient.Close()
	s.paymentClient.Close()
	s.rateLimiter.Close()

	return nil
}
