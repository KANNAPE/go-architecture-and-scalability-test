package http

import "fmt"

type AnalysisResponseDTO struct {
	TotalPosts   int
	MinTimestamp int64
	MaxTimestamp int64
	Dimension    string
	P50          uint32
	P90          uint32
	P99          uint32
}

// ToJSONMap is a func that will map the values to the json object we'll return to the user, with the dynamic stats fields 
func (dto AnalysisResponseDTO) ToJSONMap() map[string]interface{} {
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