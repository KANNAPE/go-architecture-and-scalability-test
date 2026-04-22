package stream

import (
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
	dataArray, err := service.repo.GetStream(duration)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream from repo: %w", err)
	}

	return dataArray, nil
}

type Data struct {
	ID        int64   `json:"id"`
	Timestamp uint64  `json:"timestamp"`
	Likes     *uint32 `json:"likes"`
	Comments  *uint32 `json:"comments"`
	Favorites *uint32 `json:"favorites"`
	Retweets  *uint32 `json:"retweets"`
}
