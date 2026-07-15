package hasura

import "strings"

// APIError -
type APIError struct {
	Path string `json:"path"`
	Text string `json:"error"`
	Code string `json:"code"`
}

// Error implements the error interface, formatting the non-empty Hasura API path,
// message and code carried by e as "hasura api error path=... error=... code=...".
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

// AlreadyExists reports whether e represents Hasura's "already-exists" error, e.g.
// when tracking a table or creating a REST endpoint that is already tracked/created.
func (e APIError) AlreadyExists() bool {
	return e.Code == "already-exists"
}

// PermissionDenied reports whether e represents Hasura's "permission-denied" error,
// e.g. when dropping a select permission that does not exist.
func (e APIError) PermissionDenied() bool {
	return e.Code == "permission-denied"
}
