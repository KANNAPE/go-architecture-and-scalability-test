package compute

import "context"

// compute will make all the math needed to handle percentiles

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) FilterData(ctx context.Context) Data {
	return Data{}
}

type Data struct {
}
