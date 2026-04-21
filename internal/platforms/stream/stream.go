package stream

// This package will fetch data from the Upfluence stream depending on what the user wants

import (
	"errors"
	"time"

	"kannape.com/upfluence-test/internal/services/stream"
)

type Upfluence struct {
	baseURL string
}

func NewUpfluenceStream(baseURL string) *Upfluence {
	return &Upfluence{
		baseURL: baseURL,
	}
}

func (stream *Upfluence) GetStream(duration time.Duration) ([]stream.Data, error) {
	// concatenate baseURL with api endpoint /stream

	// call api

	// return data

	return nil, errors.New("not implemented yet!")
}
