package compute

import (
	"errors"
	"fmt"
	"slices"
)

// compute will make all the math needed to handle percentiles

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

// returns an error?
func (s *Service) ComputePercentile(dataset []uint32, percentile float32) (float32, error) {
	// first, we check if the percentile value ranges between 0 and 1 (inclusive), if not => error
	if percentile < 0 || percentile > 1 {
		return 0, fmt.Errorf("percentile value %f is not valid", percentile)
	}

	// then we check if the dataset is empty
	if len(dataset) == 0 {
		return 0, errors.New("dataset is empty")
	}

	// we make sure our dataset is sorted 
	slices.Sort(dataset)

	// we first need to retrieve the rank, if it's a round value with no decimal part, then we return the value dataset[value]
	rank := percentile * float32((len(dataset) + 1))
	if rank == float32(int32(rank)) {
		//return float32(dataset[int32(rank)]), nil
	}

	// return data
	return 0, errors.New("not implemented yet!")
}
