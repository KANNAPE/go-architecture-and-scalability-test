package stream

import (
	"errors"
	"fmt"
	"time"
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

func (service *Service) GetStream(duration time.Duration) ([]Data, error) {
	_, err := service.repo.GetStream(duration)

	if err != nil {
		return nil, fmt.Errorf("failed to get stream from repo: %w", err)
	}

	return nil, errors.New("not implemented yet!")
}

type Data struct {
}
