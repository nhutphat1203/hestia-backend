package http_server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// HTTPServer bọc Gin engine
type HTTPServer struct {
	engine *gin.Engine
	port   string
}

func New(port string, logger *logger.Logger) *HTTPServer {
	r := gin.New()
	r.Use(LoggerMiddleware(logger))
	r.Use(gin.Recovery())

	return &HTTPServer{
		engine: r,
		port:   port,
	}
}

// Thêm route
func (s *HTTPServer) RegisterRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// Start server (blocking)
func (s *HTTPServer) Start() error {
	s.RegisterRoutes()
	return s.engine.Run(":" + s.port)
}
