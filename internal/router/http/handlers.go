package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"kannape.com/upfluence-test/internal/services/stream"
	"kannape.com/upfluence-test/internal/usecases"
)

type analysisHandler struct {
	streamService stream.IService
	useCase       *usecases.ComputePercentilesUseCase
}

func newAnalysisHandler(streamService stream.IService, useCase *usecases.ComputePercentilesUseCase) *analysisHandler {
	return &analysisHandler{
		streamService: streamService,
		useCase:       useCase,
	}
}

// AnalyseData handles the GET /analysis request.
func (h *analysisHandler) AnalyseData(c echo.Context) error {
	ctx := c.Request().Context()
	timestamp := time.Now().UTC().Format(time.RFC3339)
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	instance := c.Request().URL.Path

	var req AnalysisRequest

	// Binding parameters to the struct
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind query parameters", slog.String("request_id", reqID), slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Title:     "Bad Request",
			Status:    http.StatusBadRequest,
			Detail:    "failed to parse query parameters",
			Instance:  instance,
			RequestID: reqID,
			Timestamp: timestamp,
		})
	}

	// Validating the struct using our custom Echo validator
	if err := c.Validate(&req); err != nil {
		var valErr *ValidationError
		// If the error is our custom ValidationError, we inject its map into our RFC 7807 response
		if errors.As(err, &valErr) {
			slog.Info("validation failed for analysis request", slog.String("request_id", reqID))
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Title:     "Validation Error",
				Status:    http.StatusBadRequest,
				Detail:    "one or more query parameters are invalid or missing",
				Instance:  instance,
				RequestID: reqID,
				Timestamp: timestamp,
				Errors:    valErr.Errors, // The detailed map is dynamically injected here!
			})
		}
	}

	// Since validation passed, we safely parse the duration
	// (we know it won't fail because the validator already checked it)
	duration, _ := time.ParseDuration(req.Duration)
	dimension := strings.ToLower(req.Dimension)

	slog.Info("processing analysis request",
		slog.String("request_id", reqID),
		slog.String("duration", duration.String()),
		slog.String("dimension", dimension),
	)

	// Fetching data from the stream
	data, err := h.streamService.GetStream(duration)
	if err != nil {
		slog.Error("failed to fetch stream data", slog.String("request_id", reqID), slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Title:     "Internal Server Error",
			Status:    http.StatusInternalServerError,
			Detail:    "an error occurred while fetching data from the upstream service",
			Instance:  instance,
			RequestID: reqID,
			Timestamp: timestamp,
		})
	}

	// Executing Use Case
	result, err := h.useCase.Execute(ctx, data, dimension)
	if err != nil {
		slog.Error("failed to compute percentiles", slog.String("request_id", reqID), slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Title:     "Internal Server Error",
			Status:    http.StatusInternalServerError,
			Detail:    "an error occurred while computing the statistics",
			Instance:  instance,
			RequestID: reqID,
			Timestamp: timestamp,
		})
	}

	slog.Info("successfully computed analysis",
		slog.String("request_id", reqID),
		slog.Int("total_posts", result.TotalPosts),
	)

	// Return response
	responseDTO := AnalysisResponse{
		TotalPosts:   result.TotalPosts,
		MinTimestamp: result.MinTimestamp,
		MaxTimestamp: result.MaxTimestamp,
		Dimension:    result.Dimension,
		P50:          result.P50,
		P90:          result.P90,
		P99:          result.P99,
	}

	return c.JSON(http.StatusOK, responseDTO.ToJSONMap())
}
