// internal/services/compute/compute_test.go
package compute_test

import (
	"math"
	"testing"

	"kannape.com/upfluence-test/internal/services/compute"
)

func floatNearlyEquals(a, b float64) bool {
	return math.Abs(a-b) <= 0.001
}

func TestComputePercentiles(t *testing.T) {
	tests := []struct {
		name        string
		input       []uint32
		expectedP50 float64
		expectedP90 float64
		expectedP99 float64
		expectedErr string
	}{
		{
			name:        "Easy case",
			input:       []uint32{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			expectedP50: 55,
			expectedP90: 91,
			expectedP99: 99.1,
		},
		{
			name:        "Empty dataset case",
			input:       []uint32{},
			expectedErr: "dataset is empty or doesn't contain enough values",
		},
		{
			name:        "Real case",
			input:       []uint32{46, 94, 128, 128, 564, 600, 678, 1025, 200238, 387734},
			expectedP50: 582,
			expectedP90: 218987.5,
			expectedP99: 370859.313,
		},
	}

	// Testing
	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			service := compute.NewService()

			p50, err := service.ComputePercentile(testcase.input, 0.5)
			if err != nil && err.Error() != testcase.expectedErr {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !floatNearlyEquals(float64(p50), testcase.expectedP50) {
				t.Errorf("expected p50=%.5f, got %.5f", testcase.expectedP50, p50)
				return
			}

			p90, err := service.ComputePercentile(testcase.input, 0.9)
			if err != nil && err.Error() != testcase.expectedErr {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !floatNearlyEquals(float64(p90), testcase.expectedP90) {
				t.Errorf("expected p90=%.5f, got %.5f", testcase.expectedP90, p90)
				return
			}

			p99, err := service.ComputePercentile(testcase.input, 0.99)
			if err != nil && err.Error() != testcase.expectedErr {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !floatNearlyEquals(float64(p99), testcase.expectedP99) {
				t.Errorf("expected p99=%.5f, got %.5f", testcase.expectedP99, p99)
				return
			}
		})
	}
}
