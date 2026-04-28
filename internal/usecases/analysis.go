package usecases

import (
	"context"
	"errors"

	"kannape.com/upfluence-test/internal/services/stream"
)

// fonction

type AnalysisUseCase struct {

}

func NewAnalysisUseCase() *AnalysisUseCase {
 return &AnalysisUseCase{}
}

func (uc *AnalysisUseCase) ComputePercentiles(ctx context.Context, data []stream.Data, percentiles... float32) error {
	return errors.New("error")
}