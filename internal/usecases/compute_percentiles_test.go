package usecases_test

import (
	"context"
	"errors"
	"testing"

	"kannape.com/upfluence-test/internal/services/compute"
	"kannape.com/upfluence-test/internal/services/stream"
	"kannape.com/upfluence-test/internal/usecases"
)

// ptrUint32 is a small helper to easily create pointers to uint32 inline.
func ptrUint32(v uint32) *uint32 {
	return &v
}

// mockComputeService implements compute.IService to allow us to force errors
// and test the error handling branches of the Use Case.
type mockComputeService struct {
	p50Err      error
	p90Err      error
	p99Err      error
	realService compute.IService // Fallback to the real implementation if no error is forced
}

// ComputePercentile overrides the real calculation to optionally return a forced error.
func (m *mockComputeService) ComputePercentile(ctx context.Context, dataset []uint32, percentile float32) (float32, error) {
	if percentile == 0.5 && m.p50Err != nil {
		return 0, m.p50Err
	}
	if percentile == 0.9 && m.p90Err != nil {
		return 0, m.p90Err
	}
	if percentile == 0.99 && m.p99Err != nil {
		return 0, m.p99Err
	}
	// If no error is configured for this percentile, fallback to the actual math computation
	return m.realService.ComputePercentile(ctx, dataset, percentile)
}

func TestComputePercentilesUseCase_Execute(t *testing.T) {
	// We instantiate the real compute service to use it as a fallback in our mock
	realCompute := compute.NewService()
	ctx := context.Background()

	tests := []struct {
		name           string
		data           []stream.Data
		dimension      string
		mockService    *mockComputeService
		expectedError  bool
		expectedErrMsg string
		expectedLen    int
	}{
		// 1. Success paths for all 4 dimensions (Covers the entire switch statement)
		{
			name: "Valid case with likes",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)},
				{ID: 2, Timestamp: 2000, Likes: ptrUint32(20)},
			},
			dimension:     "likes",
			mockService:   &mockComputeService{realService: realCompute},
			expectedError: false,
			expectedLen:   2,
		},
		{
			name: "Valid case with comments",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Comments: ptrUint32(10)},
				{ID: 2, Timestamp: 2000, Comments: ptrUint32(20)},
			},
			dimension:     "comments",
			mockService:   &mockComputeService{realService: realCompute},
			expectedError: false,
			expectedLen:   2,
		},
		{
			name: "Valid case with favorites",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Favorites: ptrUint32(10)},
				{ID: 2, Timestamp: 2000, Favorites: ptrUint32(20)},
			},
			dimension:     "favorites",
			mockService:   &mockComputeService{realService: realCompute},
			expectedError: false,
			expectedLen:   2,
		},
		{
			name: "Valid case with retweets",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Retweets: ptrUint32(10)},
				{ID: 2, Timestamp: 2000, Retweets: ptrUint32(20)},
			},
			dimension:     "retweets",
			mockService:   &mockComputeService{realService: realCompute},
			expectedError: false,
			expectedLen:   2,
		},

		// 2. Edge case: Empty data (Covers the early exit block)
		{
			name:          "Empty stream data",
			data:          []stream.Data{},
			dimension:     "likes",
			mockService:   &mockComputeService{realService: realCompute},
			expectedError: false,
			expectedLen:   0,
		},

		// 3. Error case: Not enough data for the specified metric (Covers the MinDatasetLength block)
		{
			name: "Not enough data for the requested dimension",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)}, // Only 1 valid metric, minimum required is 2
			},
			dimension:     "likes",
			mockService:   &mockComputeService{realService: realCompute},
			expectedError: false,
			expectedLen:   1,
		},

		// 4. Error cases: Mocking compute errors (Covers the 3 error return branches of p50, p90, p99)
		{
			name: "Error triggered during P50 computation",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)},
				{ID: 2, Timestamp: 2000, Likes: ptrUint32(20)},
			},
			dimension:      "likes",
			mockService:    &mockComputeService{realService: realCompute, p50Err: errors.New("simulated p50 error")},
			expectedError:  true,
			expectedErrMsg: "failed to compute p50: simulated p50 error",
		},
		{
			name: "Error triggered during P90 computation",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)},
				{ID: 2, Timestamp: 2000, Likes: ptrUint32(20)},
			},
			dimension:      "likes",
			mockService:    &mockComputeService{realService: realCompute, p90Err: errors.New("simulated p90 error")},
			expectedError:  true,
			expectedErrMsg: "failed to compute p90: simulated p90 error",
		},
		{
			name: "Error triggered during P99 computation",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)},
				{ID: 2, Timestamp: 2000, Likes: ptrUint32(20)},
			},
			dimension:      "likes",
			mockService:    &mockComputeService{realService: realCompute, p99Err: errors.New("simulated p99 error")},
			expectedError:  true,
			expectedErrMsg: "failed to compute p99: simulated p99 error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Instantiate the Use Case with the mocked service for this specific test
			uc := usecases.NewComputePercentilesUseCase(tc.mockService)

			result, err := uc.Execute(ctx, tc.data, tc.dimension)

			// Validate error behaviors
			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected an error, but got nil")
				}
				if err.Error() != tc.expectedErrMsg {
					t.Errorf("expected error message %q, got %q", tc.expectedErrMsg, err.Error())
				}
			} else {
				// Validate success behaviors
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.TotalPosts != tc.expectedLen {
					t.Errorf("expected total posts %d, got %d", tc.expectedLen, result.TotalPosts)
				}
			}
		})
	}
}
