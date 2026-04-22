package http

import (
	"kannape.com/upfluence-test/internal/services/compute"
	"kannape.com/upfluence-test/internal/services/stream"
)

type services struct {
	streamService  stream.IService
	computeService compute.IService
}

func registerRoutes(server *Server, services services) {
	analysisHandler := newAnalysisHandler(services.streamService, services.computeService)

	server.router.HandleFunc("GET /analysis", analysisHandler.analyseData)
}
