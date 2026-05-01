package stream_test

import (
	"context"
	"errors"
	"testing"

	"kannape.com/upfluence-test/internal/services/stream"
)

// mockRepository is a fake implementation of stream.IRepository used for testing.
// It allows us to simulate the various states of the Upfluence stream connection.
type mockRepository struct {
	mockData []stream.Data
	mockErr  error
}

// GetStream returns the injected mock data and mock error.
func (m *mockRepository) GetStream(ctx context.Context) ([]stream.Data, error) {
	return m.mockData, m.mockErr
}

func TestService_GetStream(t *testing.T) {
	// Define test cases covering all branches of the GetStream function
	tests := []struct {
		name          string
		mockData      []stream.Data
		mockErr       error
		expectedError bool
		expectedLen   int
	}{
		{
			name: "Success: Repository returns data without any error",
			mockData: []stream.Data{
				{ID: 1, Timestamp: 1000},
				{ID: 2, Timestamp: 2000},
			},
			mockErr:       nil,
			expectedError: false,
			expectedLen:   2,
		},
		{
			name:          "Total Error: Repository returns an error and absolutely no data (nil array)",
			mockData:      nil,
			mockErr:       errors.New("connection timeout"),
			expectedError: true,
			expectedLen:   0,
		},
		{
			name:          "Total Error: Repository returns an error and an empty array",
			mockData:      []stream.Data{}, // len == 0
			mockErr:       errors.New("connection timeout"),
			expectedError: true,
			expectedLen:   0,
		},
		{
			name: "Partial Success: Repository returns an error but managed to fetch some data",
			mockData: []stream.Data{
				{ID: 1, Timestamp: 1000},
			}, // len == 1
			// This simulates a scenario where the stream reads a few items then breaks
			mockErr:       errors.New("connection closed prematurely"),
			expectedError: false, // The service should absorb the error and return the partial data
			expectedLen:   1,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			// 1. Setup the mock repository and the service
			repo := &mockRepository{
				mockData: testcase.mockData,
				mockErr:  testcase.mockErr,
			}
			service := stream.NewService(repo)

			// 2. Execute the method
			// The duration parameter doesn't matter here since we mock the repository
			data, err := service.GetStream(t.Context())

			// 3. Assert error expectations
			if testcase.expectedError && err == nil {
				t.Fatalf("expected an error but got nil")
			}
			if !testcase.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// 4. Assert data expectations
			if len(data) != testcase.expectedLen {
				t.Errorf("expected array length %d, got %d", testcase.expectedLen, len(data))
			}
		})
	}
}
