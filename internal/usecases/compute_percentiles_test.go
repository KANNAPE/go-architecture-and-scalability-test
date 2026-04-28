package usecases_test

import (
	"context"
	"testing"

	"kannape.com/upfluence-test/internal/services/compute"
	"kannape.com/upfluence-test/internal/services/stream"
	"kannape.com/upfluence-test/internal/usecases"
)

// ptrUint32 is a small helper to easily create pointers to uint32 inline
func ptrUint32(v uint32) *uint32 {
	return &v
}

func TestProcessStreamData(t *testing.T) {
	computeService := compute.NewService()
	uc := usecases.NewComputePercentilesUseCase(computeService)
	ctx := context.Background()

	// define test cases
	tests := []struct {
		name          string
		data          []stream.Data
		dimension     string
		expectedError bool
		expectedLen   int
		expectedP50   float32
	}{
		{
			name: "Valid case with likes",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)},
				{ID: 2, Timestamp: 2000, Likes: ptrUint32(20)},
				{ID: 3, Timestamp: 3000, Likes: ptrUint32(30)},
			},
			dimension:     "likes",
			expectedError: false,
			expectedLen:   3,
			expectedP50:   20, // 50th percentile of [10, 20, 30]
		},
		{
			name: "Valid case with missing metrics for the dimension",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)}, // has no comments
				{ID: 2, Timestamp: 2000, Comments: ptrUint32(5)},
				{ID: 3, Timestamp: 3000, Comments: ptrUint32(15)},
			},
			dimension:     "comments",
			expectedError: false,
			expectedLen:   3,
			expectedP50:   10, // 50th percentile of [5, 15]
		},
		{
			name: "Not enough data for the requested dimension",
			data: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)}, // only likes, no retweets
			},
			dimension:     "retweets",
			expectedError: true,
		},
		{
			name:          "Empty stream data",
			data:          []stream.Data{},
			dimension:     "likes",
			expectedError: false,
			expectedLen:   0,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			result, err := uc.Execute(ctx, testcase.data, testcase.dimension)

			// check error expectations
			if testcase.expectedError && err == nil {
				t.Fatalf("expected an error, but got nil")
			}
			if !testcase.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// validate results if no error was expected
			if !testcase.expectedError {
				if result.TotalPosts != testcase.expectedLen {
					t.Errorf("expected total posts %d, got %d", testcase.expectedLen, result.TotalPosts)
				}
				// check P50 only if there is data
				if testcase.expectedLen > 0 && result.P50 != testcase.expectedP50 {
					t.Errorf("expected P50 %.2f, got %.2f", testcase.expectedP50, result.P50)
				}
			}
		})
	}
}
