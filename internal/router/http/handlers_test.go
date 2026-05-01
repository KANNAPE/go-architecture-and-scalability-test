package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"kannape.com/upfluence-test/internal/services/compute"
	"kannape.com/upfluence-test/internal/services/stream"
	"kannape.com/upfluence-test/internal/usecases"

	analysisAPI "kannape.com/upfluence-test/pkg/http/analysis"
)

// mockStreamService is a fake implementation of stream.IService used for testing.
// It allows us to control exactly what GetStream returns without making real HTTP calls.
type mockStreamService struct {
	mockData []stream.Data
	mockErr  error
}

func (m *mockStreamService) GetStream(ctx context.Context) ([]stream.Data, error) {
	return m.mockData, m.mockErr
}

// ptrUint32 is a small helper to create pointers to uint32 inline
func ptrUint32(v uint32) *uint32 {
	return &v
}

func TestAnalyseData(t *testing.T) {
	// Setup standard dependencies
	computeService := compute.NewService()
	useCase := usecases.NewComputePercentilesUseCase(computeService)

	// Define test cases
	tests := []struct {
		name           string
		queryParams    string
		mockStreamData []stream.Data
		mockStreamErr  error
		expectedStatus int
		// We will partially check the response body based on these flags
		expectErrorPayload bool
		expectedTotalPosts int
	}{
		{
			name:        "Valid request and successful computation",
			queryParams: "?duration=1s&dimension=likes",
			mockStreamData: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)},
				{ID: 2, Timestamp: 2000, Likes: ptrUint32(20)},
				{ID: 3, Timestamp: 3000, Likes: ptrUint32(30)},
			},
			mockStreamErr:      nil,
			expectedStatus:     http.StatusOK,
			expectErrorPayload: false,
			expectedTotalPosts: 3,
		},
		{
			name:               "Bind error: malformed URL encoding",
			queryParams:        "?duration=%zz", // %zz is an invalid URL escape sequence
			mockStreamData:     nil,
			mockStreamErr:      nil,
			expectedStatus:     http.StatusBadRequest,
			expectErrorPayload: true,
		},
		{
			name:               "Validation error: missing parameters",
			queryParams:        "",
			mockStreamData:     nil,
			mockStreamErr:      nil,
			expectedStatus:     http.StatusBadRequest,
			expectErrorPayload: true,
		},
		{
			name:               "Validation error: invalid duration format",
			queryParams:        "?duration=invalid&dimension=likes",
			mockStreamData:     nil,
			mockStreamErr:      nil,
			expectedStatus:     http.StatusBadRequest,
			expectErrorPayload: true,
		},
		{
			name:               "Validation error: unsupported dimension",
			queryParams:        "?duration=5s&dimension=unsupported",
			mockStreamData:     nil,
			mockStreamErr:      nil,
			expectedStatus:     http.StatusBadRequest,
			expectErrorPayload: true,
		},
		{
			name:               "Stream service returns an error",
			queryParams:        "?duration=1s&dimension=likes",
			mockStreamData:     nil,
			mockStreamErr:      errors.New("stream connection failed"),
			expectedStatus:     http.StatusInternalServerError,
			expectErrorPayload: true,
		},
		{
			name:        "Use case computation fails (not enough data for dimension)",
			queryParams: "?duration=1s&dimension=comments",
			mockStreamData: []stream.Data{
				{ID: 1, Timestamp: 1000, Likes: ptrUint32(10)}, // Data is fetched but no comments
			},
			mockStreamErr:      nil,
			expectedStatus:     http.StatusOK,
			expectErrorPayload: false,
			expectedTotalPosts: 1,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			// 1. Setup the Echo context and HTTP recorder
			e := echo.New()
			e.Validator = &RequestValidator{} // Inject our custom validator

			req := httptest.NewRequest(http.MethodGet, "/analysis"+testcase.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 2. Setup the handler with the mock stream service
			mockStream := &mockStreamService{
				mockData: testcase.mockStreamData,
				mockErr:  testcase.mockStreamErr,
			}
			handler := newAnalysisHandler(mockStream, useCase)

			// 3. Execute the handler
			err := handler.AnalyseData(c)

			var actualStatus int
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					actualStatus = he.Code
				}
			} else {
				actualStatus = rec.Code
			}

			// 4. Assert the HTTP Status Code
			if actualStatus != testcase.expectedStatus {
				t.Errorf("expected status %d, got %d", testcase.expectedStatus, actualStatus)
			}

			// 5. Assert the payload structure
			if testcase.expectErrorPayload {
				var errResp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if errResp.Title == "" {
					t.Errorf("expected error payload to have a Title, got empty")
				}
			} else {
				var successResp map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &successResp); err != nil {
					t.Fatalf("failed to unmarshal success response: %v", err)
				}

				actualTotalPosts, ok := successResp["total_posts"].(float64)
				if !ok {
					t.Fatalf("total_posts field is missing or not a number")
				}
				if int(actualTotalPosts) != testcase.expectedTotalPosts {
					t.Errorf("expected total_posts %d, got %d", testcase.expectedTotalPosts, int(actualTotalPosts))
				}
			}
		})
	}
}

func TestValidator(t *testing.T) {
	v := &RequestValidator{}

	// Case 1: Valid struct
	validReq := &analysisAPI.AnalysisRequest{
		Duration:  "5m",
		Dimension: "likes",
	}
	if err := v.Validate(validReq); err != nil {
		t.Errorf("expected no error for valid request, got %v", err)
	}

	// Case 2: Invalid interface type (should return nil as per our safe type assertion)
	if err := v.Validate("not a struct"); err != nil {
		t.Errorf("expected nil when passing unknown type, got %v", err)
	}

	// Case 3: Missing parameters to trigger ValidationError
	invalidReq := &analysisAPI.AnalysisRequest{
		Duration:  "",
		Dimension: "",
	}
	err := v.Validate(invalidReq)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	// Case 4: Verify ValidationError properties
	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected error to be of type *ValidationError")
	}
	if valErr.Error() != "validation failed" {
		t.Errorf("unexpected error message: %s", valErr.Error())
	}
	if len(valErr.Errors) != 2 {
		t.Errorf("expected 2 validation errors, got %d", len(valErr.Errors))
	}
}

func TestAnalysisResponseDTO(t *testing.T) {
	dto := analysisAPI.AnalysisResponse{
		TotalPosts:   42,
		MinTimestamp: 1000,
		MaxTimestamp: 5000,
		Dimension:    "retweets",
		P50:          10.5,
		P90:          20.5,
		P99:          30.5,
	}

	m := dto.ToJSONMap()

	if m["total_posts"] != 42 {
		t.Errorf("expected total_posts to be 42, got %v", m["total_posts"])
	}
	if m["retweets_p50"] != float32(10.5) {
		t.Errorf("expected retweets_p50 to be 10.5, got %v", m["retweets_p50"])
	}
	if m["retweets_p99"] != float32(30.5) {
		t.Errorf("expected retweets_p99 to be 30.5, got %v", m["retweets_p99"])
	}

	// Ensure other keys don't dynamically bleed
	if _, exists := m["likes_p50"]; exists {
		t.Errorf("likes_p50 should not exist when dimension is retweets")
	}
}

func TestServerInitialization(t *testing.T) {
	// Simple test to ensure NewServer doesn't panic and wires correctly
	mockStream := &mockStreamService{}
	server := NewServer(mockStream)

	if server == nil {
		t.Fatalf("expected server to be initialized")
	}
	if server.router == nil {
		t.Errorf("expected router to be initialized")
	}
}
