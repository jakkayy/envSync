package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jakkayy/envSync/internal/database"
	"github.com/jakkayy/envSync/internal/server/handlers"
	"github.com/jakkayy/envSync/internal/server/middleware"
)

type Server struct {
	router *gin.Engine
	port   string
}

func NewServer(port string) *Server {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	dbType := os.Getenv("DB_TYPE")
	dbDSN := os.Getenv("DB_DSN")
	db, err := database.InitDB(dbType, dbDSN)
	if err != nil {
		fmt.Printf("Warning: Database initialization failed: %v\n", err)
	} else {
		if err := database.AutoMigrate(db); err != nil {
			fmt.Printf("Warning: Database auto migration failed: %v\n", err)
		}
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerAndRequestID())

	s := &Server{
		router: router,
		port:   port,
	}

	s.setupRoutes()
	return s
}

func (s *Server) Engine() *gin.Engine {
	return s.router
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	apiV1 := s.router.Group("/api/v1")
	apiV1.Use(middleware.AuthRequired())
	{
		apiV1.POST("/projects", handlers.CreateProject)
		apiV1.GET("/projects/:id/history", handlers.GetProjectHistory)
		apiV1.POST("/sync/push", handlers.PushSync)
		apiV1.GET("/sync/pull", handlers.PullSync)
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		fmt.Printf("🚀 envSync Central API Server running on http://localhost:%s\n", s.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down API Server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
