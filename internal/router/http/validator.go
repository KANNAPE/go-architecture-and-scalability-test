package http

import (
	"slices"
	"strings"
	"time"

	analysisAPI "kannape.com/upfluence-test/pkg/http/analysis"
)

// ValidationError is a custom error type that holds a map of invalid fields and their reasons.
type ValidationError struct {
	Errors map[string]interface{}
}

// Error implements the standard error interface.
func (e *ValidationError) Error() string {
	return "validation failed"
}

// RequestValidator implements the echo.Validator interface using the standard library.
type RequestValidator struct{}

// Validate checks the struct fields and populates the error map if any rules are broken.
func (v *RequestValidator) Validate(i interface{}) error {
	// We check if the interface is our AnalysisRequest
	if req, ok := i.(*analysisAPI.AnalysisRequest); ok {
		validationErrors := make(map[string]interface{})

		// Validate 'duration'
		if strings.TrimSpace(req.Duration) == "" {
			validationErrors["duration"] = "this parameter is required"
		} else if _, err := time.ParseDuration(req.Duration); err != nil {
			validationErrors["duration"] = "invalid format: expected a time duration like '5s', '10m', or '24h'"
		}

		// Validate 'dimension'
		allowedDimensions := []string{"likes", "comments", "favorites", "retweets"}
		if strings.TrimSpace(req.Dimension) == "" {
			validationErrors["dimension"] = "this parameter is required"
		} else if !slices.Contains(allowedDimensions, strings.ToLower(req.Dimension)) {
			validationErrors["dimension"] = "unsupported dimension: handled dimensions are 'likes', 'comments', 'favorites', and 'retweets'"
		}

		// If the map is not empty, it means we have validation errors
		if len(validationErrors) > 0 {
			return &ValidationError{Errors: validationErrors}
		}
	}

	return nil
}
