package http

// ErrorResponse represents the standardized JSON error payload
// sent to the client when an API request fails.
type ErrorResponse struct {
	Title     string                 `json:"title"`
	Status    int                    `json:"status"`
	Detail    string                 `json:"detail,omitempty"`
	Instance  string                 `json:"instance,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Errors    map[string]interface{} `json:"errors,omitempty"`
}
