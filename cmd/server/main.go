package main

import (
	"fmt"

	"kannape.com/upfluence-test/internal/platforms/stream"
	httpx "kannape.com/upfluence-test/internal/router/http"
)

func main() {
	fmt.Println("Hello World!")

	streamRepo := stream.NewUpfluenceStream("https://stream.upfluence.co")

	server := httpx.NewServer(streamRepo)

	if err := server.Start(); err != nil {
		panic(err)
	}
}
