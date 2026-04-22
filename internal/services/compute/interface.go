package compute

type IService interface {
	ComputePercentiles(metrics []uint32) Data
}

type Data struct {
	P50 uint32
	P90 uint32
	P99 uint32
}
