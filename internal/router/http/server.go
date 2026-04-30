package http

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"kannape.com/upfluence-test/internal/services/compute"
	"kannape.com/upfluence-test/internal/services/stream"
	"kannape.com/upfluence-test/internal/usecases"
)

type Server struct {
	router     *echo.Echo
	streamRepo stream.IRepository
}

// NewServer initializes a new HTTP server using the Echo framework.
func NewServer(streamRepo stream.IRepository) *Server {
	e := echo.New()

	// Attach global middlewares
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// Register our custom validator
	e.Validator = &RequestValidator{}

	return &Server{
		router:     e,
		streamRepo: streamRepo,
	}
}

// Start configures the dependencies, registers routes, and starts listening.
func (server *Server) Start(ctx context.Context) error {
	// Initialize services
	streamService := stream.NewService(server.streamRepo)
	computeService := compute.NewService()

	// Initialize Use Cases
	computePercentilesUseCase := usecases.NewComputePercentilesUseCase(computeService)

	// Initialize Handlers
	analysisHdl := newAnalysisHandler(streamService, computePercentilesUseCase)

	// Register HTTP routes
	registerRoutes(server, analysisHdl)

	// Start the HTTP server on port 8080
	slog.InfoContext(ctx, "Listening on localhost:8080...")
	return server.router.Start(":8080")
}
