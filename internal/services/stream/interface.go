package stream

import (
	"context"
)

type IRepository interface {
	GetStream(ctx context.Context) ([]Data, error)
}

type IService interface {
	GetStream(ctx context.Context) ([]Data, error)
}

type Data struct {
	ID        int64   `json:"id"`
	Timestamp int64   `json:"timestamp"`
	Likes     *uint32 `json:"likes"`
	Comments  *uint32 `json:"comments"`
	Favorites *uint32 `json:"favorites"`
	Retweets  *uint32 `json:"retweets"`
}
