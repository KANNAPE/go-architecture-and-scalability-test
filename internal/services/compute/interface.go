package compute

import "context"

type IService interface {
	FilterData(ctx context.Context) Data
}
