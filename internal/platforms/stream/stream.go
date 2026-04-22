package stream

// This package will fetch data from the Upfluence stream depending on what the user wants

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

func (s *Upfluence) GetStream(duration time.Duration) ([]stream.Data, error) {
	// concatenate baseURL with api endpoint /stream
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	url := s.baseURL + "/stream"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating the request: %w", err)
	}

	// calling api
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error when calling api endpoint: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http response code %d", response.StatusCode)
	}

	// scanning for data
	var results []stream.Data
	var errs []error

	scanner := bufio.NewScanner(response.Body)

	for scanner.Scan() {
		line := scanner.Text()

		// since each data from the stream starts with "data: "
		line = strings.TrimPrefix(line, "data: ")

		// simple security, avoiding white lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		data := struct {
			Payload stream.Data `json:"instagram_media"` // trouver comme faire mieux et récup les autres types de data
		}{}

		if err := json.Unmarshal([]byte(line), &data); err != nil {
			errs = append(errs, fmt.Errorf("Failed to unmarshal data \"%s\": %w", line, err))
			continue
		}

		results = append(results, data.Payload)
	}

	// making sure the scan stopped because the duration was reached, not anything else
	if err := scanner.Err(); err != nil {
		if ctx.Err() != context.DeadlineExceeded {
			return results, fmt.Errorf("error while reading the stream: %w", err)
		}
	}

	// return data
	return results, errors.Join(errs...)
}
