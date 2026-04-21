package stream

import "time"

type IRepository interface {
	GetStream(duration time.Duration) ([]Data, error)
}

type IService interface {
	GetStream(duration time.Duration) ([]Data, error)
}
