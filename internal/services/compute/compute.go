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

	// then we check if the dataset is empty or has exactly one value
	if len(dataset) < 2 {
		return 0, errors.New("dataset is empty or doesn't contain enough values")
	}

	// we make sure our dataset is sorted
	slices.Sort(dataset)

	// we first need to retrieve the rank and the "integer" part of the rank
	rank := percentile * float32((len(dataset) - 1))
	integerRank := int32(rank)

	// if the rank and its integer part are equal, it means the data at the index rank is our percentile value
	if rank == float32(integerRank) {
		return float32(dataset[int32(rank)]), nil
	}

	// else, we need to compute the fractional interpolated value of rank, as follows
	fractionalRank := rank - float32(integerRank)

	// we then need to fetch the two numbers that are at rank "integerRank" and "integerRank + 1"
	borderMinValue := dataset[integerRank]
	borderMaxValue := dataset[integerRank+1] // at this point, integerRank strictly cannot be equal to len()-1 so we don't even need to check

	// we can retrieve the interpolated rank
	rank = fractionalRank*(float32(borderMaxValue-borderMinValue)) + float32(borderMinValue)

	// return data
	return rank, nil
}
