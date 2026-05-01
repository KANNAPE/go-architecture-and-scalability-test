// internal/services/compute/compute_test.go
package compute_test

import (
	"fmt"
	"math"
	"testing"

	"kannape.com/upfluence-test/internal/services/compute"
)

func floatNearlyEquals(a, b float64) bool {
	return math.Abs(a-b) <= 0.001
}

func TestComputePercentiles(t *testing.T) {
	// defining test cases
	tests := []struct {
		name           string
		input          []uint32
		percentile     float32
		expectedOutput float64
		expectedErr    string
	}{
		{
			name:           "Easy case",
			input:          []uint32{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			percentile:     0.5,
			expectedOutput: 55,
		},
		{
			name:        "Empty dataset case",
			input:       []uint32{},
			expectedErr: fmt.Sprintf("dataset is empty or doesn't contain enough values (minimum expected values: %d)", compute.MinDatasetLength),
		},
		{
			name:        "Nil dataset case",
			input:       nil,
			expectedErr: fmt.Sprintf("dataset is empty or doesn't contain enough values (minimum expected values: %d)", compute.MinDatasetLength),
		},
		{
			name:           "Real case (50th percentile)",
			input:          []uint32{46, 94, 128, 128, 564, 600, 678, 1025, 200238, 387734},
			percentile:     0.5,
			expectedOutput: 582,
		},
		{
			name:           "Real case (90th percentile)",
			input:          []uint32{46, 94, 128, 128, 564, 600, 678, 1025, 200238, 387734},
			percentile:     0.9,
			expectedOutput: 218987.5,
		},
		{
			name:           "Real case (99th percentile)",
			input:          []uint32{46, 94, 128, 128, 564, 600, 678, 1025, 200238, 387734},
			percentile:     0.99,
			expectedOutput: 370859.313,
		},
		{
			name:        "Invalid percentile case",
			input:       []uint32{10, 20, 30},
			percentile:  1.5,
			expectedErr: "percentile value 1.500000 is not valid",
		},
		{
			name:           "Unsorted dataset case",
			input:          []uint32{30, 10, 20}, // the service should sort it automatically
			percentile:     0.5,
			expectedOutput: 20,
		},
	}

	// Testing
	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			service := compute.NewService()

			percentile, err := service.ComputePercentile(t.Context(), testcase.input, testcase.percentile)
			if err != nil && err.Error() != testcase.expectedErr {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !floatNearlyEquals(float64(percentile), float64(testcase.expectedOutput)) {
				t.Errorf("expected p50=%.5f, got %.5f", testcase.expectedOutput, percentile)
				return
			}
		})
	}
}
