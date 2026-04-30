package http

func registerRoutes(server *Server, analysisHdl *analysisHandler) {
	server.router.GET("/analysis", analysisHdl.AnalyseData)
}
