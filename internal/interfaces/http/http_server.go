package http_server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nhutphat1203/hestia-backend/internal/config"
	"github.com/nhutphat1203/hestia-backend/internal/infrastructure/websocket"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
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
}

func New(cfg *config.Config, logger *logger.Logger, websocketHub *websocket.Hub) *HTTPServer {
	r := gin.New()
	r.Use(LoggerMiddleware(logger))
	r.Use(gin.Recovery())

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	return &HTTPServer{
		engine:       r,
		cfg:          cfg,
		logger:       logger,
		server:       srv,
		websocketHub: websocketHub,
	}
}

func (s *HTTPServer) RegisterRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	ws := s.engine.Group("/ws/v1/env")
	{
		ws.GET("", func(c *gin.Context) {
			s.websocketHub.ServeWS(c)
		})
	}
}

func (s *HTTPServer) Start() error {
	s.RegisterRoutes()
	s.logger.Info("HTTP server starting on port " + s.cfg.ServerPort)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	s.logger.Info("HTTP server shutting down...")
	return s.server.Shutdown(ctx)
}
