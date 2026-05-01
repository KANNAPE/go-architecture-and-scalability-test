package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"kannape.com/upfluence-test/internal/platforms/stream"
	"kannape.com/upfluence-test/internal/router/http"
)

func main() {
	// We create a global context which will serve to pass additional information between processes and requests, that will be
	// gracefully closed when receiving os system calls (interrupt and terminated)
	appCtx, stopAppCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopAppCtx()

	// Initialize a structured JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.InfoContext(appCtx, "starting Upfluence Analysis API server...")

	// Initialize the Upfluence stream platform
	streamRepo := stream.NewUpfluenceStream("https://stream.upfluence.co")

	// Initialize and start the HTTP server
	server := http.NewServer(streamRepo)

	if err := server.Start(appCtx); err != nil {
		slog.ErrorContext(appCtx, "server crashed", "error", err.Error())
		panic(fmt.Errorf("http server crashed: %w", err))
	}
}

/*
	todo:
	- rajouter un makefile pour faire genre je m'y connais de zinzin
	- rajouter un .env avec les variables globales dedans
	- faire ce ptn de readme
	- envoyer et être embauché (fuck la prison)
*/
