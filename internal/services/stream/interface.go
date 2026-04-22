package stream

import "time"

type IRepository interface {
	GetStream(duration time.Duration) ([]Data, error)
}

type IService interface {
	GetStream(duration time.Duration) ([]Data, error)
}

type Data struct {
	ID        int64   `json:"id"`
	Timestamp uint64  `json:"timestamp"`
	Likes     *uint32 `json:"likes"`
	Comments  *uint32 `json:"comments"`
	Favorites *uint32 `json:"favorites"`
	Retweets  *uint32 `json:"retweets"`
}
