package analysis

import "fmt"

// AnalysisRequest represents the expected query parameters
// for the GET /analysis endpoint.
type AnalysisRequest struct {
	Duration  string `query:"duration"`
	Dimension string `query:"dimension"`
}

// AnalysisResponse represents the successful JSON response payload.
// The percentile fields are dynamically mapped in the ToJSONMap function.
type AnalysisResponse struct {
	TotalPosts   int
	MinTimestamp int64
	MaxTimestamp int64
	Dimension    string
	P50          float32
	P90          float32
	P99          float32
}

// ToJSONMap dynamically generates the response map to ensure
// the percentile keys exactly match the requested dimension.
func (dto AnalysisResponse) ToJSONMap() map[string]interface{} {
	return map[string]interface{}{
		"total_posts":       dto.TotalPosts,
		"minimum_timestamp": dto.MinTimestamp,
		"maximum_timestamp": dto.MaxTimestamp,
		// dynamic stats fields, depending on the dimension
		fmt.Sprintf("%s_p50", dto.Dimension): dto.P50,
		fmt.Sprintf("%s_p90", dto.Dimension): dto.P90,
		fmt.Sprintf("%s_p99", dto.Dimension): dto.P99,
	}
}
