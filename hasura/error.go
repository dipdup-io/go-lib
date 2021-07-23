package hasura

import "strings"

// APIError -
type APIError struct {
	Path string `json:"path"`
	Text string `json:"error"`
	Code string `json:"code"`
}

// Error -
func (e APIError) Error() string {
	var builder strings.Builder
	builder.WriteString("hasura api error")
	if e.Path != "" {
		builder.WriteString(" path=")
		builder.WriteString(e.Path)
	}
	if e.Text != "" {
		builder.WriteString(" error=")
		builder.WriteString(e.Text)
	}
	if e.Code != "" {
		builder.WriteString(" code=")
		builder.WriteString(e.Code)
	}
	return builder.String()
}

// AlreadyExists -
func (e APIError) AlreadyExists() bool {
	return e.Code == "already-exists"
}
