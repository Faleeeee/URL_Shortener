package domain

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
	Error   *ErrorDetails `json:"error"`
	Meta    *Meta         `json:"meta"`
}

// ErrorDetails represents the structure for error details
type ErrorDetails struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}

// Meta represents metadata for pagination or other extra info
type Meta struct {
	Page  int   `json:"page,omitempty"`
	Limit int   `json:"limit,omitempty"`
	Total int64 `json:"total,omitempty"`
}
