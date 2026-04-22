package compute

import "slices"

// compute will make all the math needed to handle percentiles

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

// returns an error?
func (s *Service) ComputePercentiles(metrics []uint32) Data {
	slices.Sort(metrics)

	// compute indexes for percentiles 50, 90, 99

	// return data

	return Data{}
}
