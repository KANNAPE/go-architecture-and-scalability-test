package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"kannape.com/upfluence-test/internal/services/compute"
	"kannape.com/upfluence-test/internal/services/stream"
)

// ComputePercentilesResult holds the computed statistics for a specific dimension
// extracted from the stream data.
type ComputePercentilesResult struct {
	TotalPosts   int
	MinTimestamp int64
	MaxTimestamp int64
	Dimension    string
	P50          float32
	P90          float32
	P99          float32
}

// ComputePercentilesUseCase is responsible for parsing stream data,
// extracting the relevant metrics, and calculating the required percentiles.
type ComputePercentilesUseCase struct {
	computeService compute.IService
}

// NewComputePercentilesUseCase creates a new instance of the use case,
// injecting the required compute service for percentile calculations.
func NewComputePercentilesUseCase(computeService compute.IService) *ComputePercentilesUseCase {
	return &ComputePercentilesUseCase{
		computeService: computeService,
	}
}

// Execute parses the raw stream data, extracts the relevant metrics
// for the requested dimension, and computes the required percentiles.
func (uc *ComputePercentilesUseCase) Execute(ctx context.Context, data []stream.Data, dimension string) (*ComputePercentilesResult, error) {
	var metrics []uint32
	var err error

	// initialize minTimestamp to the maximum possible int64 value
	// initialize maxTimestamp to 0
	minTimestamp := ^int64(0)
	maxTimestamp := int64(0)

	// handle the case where the stream returned no data
	if len(data) == 0 {
		slog.WarnContext(ctx, "no data returned from stream")
		return &ComputePercentilesResult{
			TotalPosts: 0,
			Dimension:  dimension,
		}, nil
	}

	for _, item := range data {
		// update timestamp boundaries
		if item.Timestamp < minTimestamp {
			minTimestamp = item.Timestamp
		}
		if item.Timestamp > maxTimestamp {
			maxTimestamp = item.Timestamp
		}

		// extract the appropriate metric based on the requested dimension
		switch dimension {
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

	resultStruct := &ComputePercentilesResult{
		TotalPosts:   len(data),
		Dimension:    dimension,
		MinTimestamp: minTimestamp,
		MaxTimestamp: maxTimestamp,
	}

	// if we have not enough data to compute the percentiles, we early return
	if len(metrics) < compute.MinDatasetLength {
		slog.WarnContext(ctx, "not enough data to compute percentiles for dimension", slog.String("dimension", dimension))
		return resultStruct, nil
	}

	// compute the 50th percentile
	resultStruct.P50, err = uc.computeService.ComputePercentile(ctx, metrics, 0.5)
	if err != nil {
		slog.ErrorContext(ctx, "failed to compute p50")
		return nil, fmt.Errorf("failed to compute p50: %w", err)
	}

	// compute the 90th percentile
	resultStruct.P90, err = uc.computeService.ComputePercentile(ctx, metrics, 0.9)
	if err != nil {
		slog.ErrorContext(ctx, "failed to compute p90")
		return nil, fmt.Errorf("failed to compute p90: %w", err)
	}

	// compute the 99th percentile
	resultStruct.P99, err = uc.computeService.ComputePercentile(ctx, metrics, 0.99)
	if err != nil {
		slog.ErrorContext(ctx, "failed to compute p99")
		return nil, fmt.Errorf("failed to compute p99: %w", err)
	}

	// return the final aggregated result
	return resultStruct, nil
}
