package http

import (
	"encoding/json"
	"net/http"
	"time"

	"kannape.com/upfluence-test/internal/services/compute"
	"kannape.com/upfluence-test/internal/services/stream"
)

type analysisHandler struct {
	streamService  stream.IService
	computeService compute.IService
}

// newAnalysisHandler creates a new analysis handler for HTTP requests using the given services
func newAnalysisHandler(streamService stream.IService, computeService compute.IService) *analysisHandler {
	return &analysisHandler{
		streamService:  streamService,
		computeService: computeService,
	}
}

// analyseData will check for two query arguments, duration and dimension
// - duration is an integer that can either be follow by "s", "m", or "h", and represents the duration that the stream connection will be maintained
// - dimension is the metric we'll monitor, and can be either "likes", "favourites", "comments", or "retweets"
// If either of these query argument are missing, we'll throw an error 400
// When the duration is reached, the connection to the stream is cut and we'll return a JSON payload that'll contain the total number of posts analyzed,
// the timestamp range of the monitored, and the 50, 90, and 99 percentiles of the dimension that we chose to monitor
func (h *analysisHandler) analyseData(w http.ResponseWriter, r *http.Request) {
	data, err := h.streamService.GetStream(time.Second * 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
