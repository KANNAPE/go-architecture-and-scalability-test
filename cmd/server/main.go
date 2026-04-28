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
- finir mes tests
- logger des trucs
- meilleure validation des inputs => créer des structs qui vont contenir mes champs (requête, duration, dimension) pour pouvoir mieux logger et pas juste découvrir ce qu'il manque au jour le jour comme là mtn
- créer un dossier pkg qui va contenir une struct pour chaque route (là en l'occurence juste une pour /analysis) /\ la struct en question
- pareil avec la struct qui va retourner les erreurs (et donc virer les http.Error)
- passer sur Echo pour les middlewares parce qu'on en veut
- commentaires et finir le README



type ErrorResponse struct {
	Title     string                 `json:"title"`
	Status    int                    `json:"status"`
	Detail    string                 `json:"detail,omitempty"`
	Instance  string                 `json:"instance,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Errors    map[string]interface{} `json:"errors,omitempty"`
}

*/
