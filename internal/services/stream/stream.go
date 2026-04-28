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
		if len(dataArray) == 0 {
			return nil, fmt.Errorf("failed to get stream from repo: %w", err)
		}
		
		// errors from reading some data, but the array is not empty, we log then discard
		fmt.Printf("warning: some data couldn't be parsed: %s", err.Error())
	}

	return dataArray, nil
}
