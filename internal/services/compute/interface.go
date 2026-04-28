package compute

import "context"

type IService interface {
	ComputePercentile(ctx context.Context, dataset []uint32, percentile float32) (float32, error)
}
