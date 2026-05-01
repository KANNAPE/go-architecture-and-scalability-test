package stream

import (
	"context"
	"fmt"
	"log/slog"
)

// stream package will get all the data from the stream in the given time window

type Service struct {
	repo IRepository
}

func NewService(repo IRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (service *Service) GetStream(ctx context.Context) ([]Data, error) {
	dataArray, err := service.repo.GetStream(ctx)
	if err != nil {
		if len(dataArray) == 0 {
			slog.ErrorContext(ctx, "failed to get stream")
			return nil, fmt.Errorf("failed to get stream from repo: %w", err)
		}

		// Errors occurred while reading some data, but the array is not empty.
		// We log the warning and return the partial data safely.
		slog.Warn("partial data fetched with errors",
			slog.String("error", err.Error()),
			slog.Int("retrieved_items", len(dataArray)),
		)
	}

	return dataArray, nil
}
