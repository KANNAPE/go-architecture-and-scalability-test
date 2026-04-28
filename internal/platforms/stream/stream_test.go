package stream_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	platform_stream "kannape.com/upfluence-test/internal/platforms/stream"
)

func TestUpfluence_GetStream(t *testing.T) {
	// defining test cases
	tests := []struct {
		name        string
		mockHandler http.HandlerFunc
		duration    time.Duration
		expectedLen int
		expectedErr bool
	}{
		{
			name: "Real case",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")

				fmt.Fprintln(w, `data: {"instagram_media": {"id": 1, "timestamp": 1600000000, "likes": 100}}`)
				fmt.Fprintln(w, `data: {"tweet": {"id": 2, "timestamp": 1600000005, "retweets": 50}}`)

				// making sure the client has enough time to read data
				time.Sleep(50 * time.Millisecond)
			},
			duration:    100 * time.Millisecond,
			expectedLen: 2,
			expectedErr: false,
		},
		{
			name: "Malformed JSON case",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")

				// sending a mix of valid data, empty lines, and malformed JSON
				fmt.Fprintln(w, `data: {"instagram_media": {"id": 1, "timestamp": 1600000000, "likes": 100}}`)
				fmt.Fprintln(w, `\n`) // blank line
				fmt.Fprintln(w, `data: THIS IS NOT VALID JSON`)
				fmt.Fprintln(w, `data: {"tweet": {"id": 2, "timestamp": 1600000005, "retweets": 50}}`)

				time.Sleep(50 * time.Millisecond)
			},
			duration:    100 * time.Millisecond,
			expectedLen: 2, // should still be 2 because the invalid line and the blank line are skipped
			expectedErr: false,
		},
		{
			name: "Error 500 case",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			duration:    50 * time.Millisecond,
			expectedLen: 0,
			expectedErr: true,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			mockServer := httptest.NewServer(testcase.mockHandler)
			defer mockServer.Close()

			// instantiate the client with the mock server's dynamically generated URL
			client := platform_stream.NewUpfluenceStream(mockServer.URL)

			data, err := client.GetStream(testcase.duration)

			if testcase.expectedErr && err == nil {
				t.Fatalf("expected an error but got nil")
			}
			if !testcase.expectedErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(data) != testcase.expectedLen {
				t.Errorf("expected %d items in stream, got %d", testcase.expectedLen, len(data))
			}
		})
	}
}
