package http_server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nhutphat1203/hestia-backend/internal/config"
	"github.com/nhutphat1203/hestia-backend/internal/infrastructure/auth"
	"github.com/nhutphat1203/hestia-backend/internal/infrastructure/websocket"
	"github.com/nhutphat1203/hestia-backend/internal/interfaces/http/handlers"
	"github.com/nhutphat1203/hestia-backend/internal/interfaces/http/middlewares"
	service "github.com/nhutphat1203/hestia-backend/internal/services"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
	"github.com/nhutphat1203/hestia-backend/pkg/response"
)

func LoggerMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Incoming request: " + c.Request.Method + " " + c.Request.URL.Path)
		c.Next()
		status := c.Writer.Status()
		logger.Info("Response status: " + strconv.Itoa(status))
	}
}

type HTTPServer struct {
	engine       *gin.Engine
	cfg          *config.Config
	logger       *logger.Logger
	server       *http.Server
	websocketHub *websocket.Hub
	authService  *service.AuthService
}

func New(cfg *config.Config, logger *logger.Logger, websocketHub *websocket.Hub, authService *service.AuthService) *HTTPServer {
	r := gin.New()
	r.Use(LoggerMiddleware(logger))
	r.Use(gin.Recovery())

	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      r,
		ReadTimeout:  cfg.ServerReadTimeout * time.Second,
		WriteTimeout: cfg.ServerWriteTimeout * time.Second,
		IdleTimeout:  cfg.ServerIdleTimeout * time.Second,
	}

	return &HTTPServer{
		engine:       r,
		cfg:          cfg,
		logger:       logger,
		server:       srv,
		websocketHub: websocketHub,
		authService:  authService,
	}
}

func (s *HTTPServer) RegisterRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		response.SendSuccess(c, http.StatusOK, "ok", nil)
	})

	// --- JWT Authentication Setup ---
	jwtService := auth.NewJWTService(s.cfg.JWTSecret, "hestia-api", int(s.cfg.JWTExpiration.Minutes()))
	authHandler := handlers.NewAuthHandler(s.authService)

	// --- Public Routes ---
	authRoutes := s.engine.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh_token", authHandler.RefreshToken)
		authRoutes.POST("/logout", authHandler.Logout)
	}

	// --- Protected API Routes ---
	authGroup := s.engine.Group("/api/v1")
	{
		authGroup.Use(middlewares.AuthMiddleware(jwtService))
		// Add other protected API routes here
		authGroup.GET("/protected", func(c *gin.Context) {
			response.SendSuccess(c, http.StatusOK, "This is a protected route", nil)
		})
	}

	// --- Protected WebSocket Routes ---
	ws := s.engine.Group("/ws/v1/env")
	{
		ws.Use(middlewares.AuthMiddleware(jwtService))
		ws.GET("", func(c *gin.Context) {
			s.websocketHub.ServeWS(c)
		})
	}

}

func (s *HTTPServer) Start() error {
	s.RegisterRoutes()
	s.logger.Info("HTTP server starting on address " + s.cfg.ServerAddress)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	s.logger.Info("HTTP server shutting down...")
	return s.server.Shutdown(ctx)
}
