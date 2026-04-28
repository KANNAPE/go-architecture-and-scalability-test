package main

import (
	"log/slog"
	"os"

	"kannape.com/upfluence-test/internal/platforms/stream"
	"kannape.com/upfluence-test/internal/router/http"
)

func main() {
	// Initialize a structured JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting Upfluence Analysis API server...")

	// Initialize the Upfluence stream platform
	streamRepo := stream.NewUpfluenceStream("https://stream.upfluence.co")

	// Initialize and start the HTTP server
	server := http.NewServer(streamRepo)

	slog.Info("Server is configured and ready to listen on port 8080")
	if err := server.Start(); err != nil {
		slog.Error("Server crashed", "error", err.Error())
		os.Exit(1)
	}
}
