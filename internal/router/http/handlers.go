package http

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"kannape.com/upfluence-test/internal/services/stream"
)

type analysisHandler struct {
	streamService  stream.IService
	analysis
}

// newAnalysisHandler creates a new analysis handler for HTTP requests using the given services
func newAnalysisHandler(streamService stream.IService, ) *analysisHandler {
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
	ctx := r.Context()

	// fetching query parameters
	durationStr := r.URL.Query().Get("duration")
	dimensionStr := r.URL.Query().Get("dimension")

	if strings.TrimSpace(durationStr) == "" {
		http.Error(w, "'duration' parameter is missing!", http.StatusBadRequest)
		return
	}

	// parsing duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		http.Error(w, "duration format invalid! (value should either be in seconds 's', minutes 'm', or hours 'h')", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(dimensionStr) == "" {
		http.Error(w, "'dimension' parameter is missing!", http.StatusBadRequest)
		return
	}

	// checking for allowed dimension
	allowedDimensions := []string{"likes", "comments", "favorites", "retweets"} // passer ca dans des variables d'environement je pense

	isValidDimension := slices.Contains(allowedDimensions, strings.ToLower(dimensionStr))
	if !isValidDimension {
		http.Error(w, "dimension is not supported! (current handled dimension are 'likes', 'comments', 'favorites', and 'retweets')", http.StatusBadRequest)
		return
	}

	// reading data from stream
	data, err := h.streamService.GetStream(duration)
	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var metrics []uint32
	minTimestamp := ^int64(0)
	maxTimestamp := int64(0)

	for _, item := range data {
		// timestamps
		if item.Timestamp < minTimestamp {
			minTimestamp = item.Timestamp
		}
		if item.Timestamp > maxTimestamp {
			maxTimestamp = item.Timestamp
		}

		switch dimensionStr {
		case "likes":
			if item.Likes != nil {
				metrics = append(metrics, *item.Likes)
			}
		case "comments":
			if item.Comments != nil {
				metrics = append(metrics, *item.Comments)
			}
		case "favorites":
			if item.Favorites != nil {
				metrics = append(metrics, *item.Favorites)
			}
		case "retweets":
			if item.Retweets != nil {
				metrics = append(metrics, *item.Retweets)
			}
		}
	}

	// 50th percentile
	p50, err := h.computeService.ComputePercentile(ctx, metrics, 0.5)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 90th percentile
	p90, err := h.computeService.ComputePercentile(ctx, metrics, 0.9)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 99th percentile
	p99, err := h.computeService.ComputePercentile(ctx, metrics, 0.99)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// payload mapped to DTO
	responseDTO := AnalysisResponseDTO{
		TotalPosts:   len(data),
		MinTimestamp: minTimestamp,
		MaxTimestamp: maxTimestamp,
		Dimension:    dimensionStr,
		P50:          p50,
		P90:          p90,
		P99:          p99,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseDTO.ToJSONMap()); err != nil {
		panic(err)
	}
}
