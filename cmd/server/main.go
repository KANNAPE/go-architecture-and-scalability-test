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

/*****************************
TODO:

- logger des trucs

- créer un dossier pkg qui va contenir une struct pour chaque route (là en l'occurence juste une pour /analysis) /\ la struct en question


- commentaires et finir le README

*/
