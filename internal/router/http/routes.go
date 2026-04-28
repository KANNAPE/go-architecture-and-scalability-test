package http

import (
	"kannape.com/upfluence-test/internal/services/compute"
	"kannape.com/upfluence-test/internal/services/stream"
	"kannape.com/upfluence-test/internal/usecases"
)

type services struct {
	streamService  stream.IService
	computeService compute.IService
}

func registerRoutes(server *Server, services services) {
	analysisUseCases := usecases.NewAnalysisUseCase()

	analysisHandler := newAnalysisHandler(services.streamService, analysisUseCase)


	server.router.HandleFunc("GET /analysis", analysisHandler.analyseData)
}
