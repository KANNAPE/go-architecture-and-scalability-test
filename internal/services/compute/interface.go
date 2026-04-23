package compute

type IService interface {
	ComputePercentile(metrics []uint32, percentile float32) (float32, error)
}
