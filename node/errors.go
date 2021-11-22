package node

import (
	"fmt"
)

// RequestError -
type RequestError struct {
	Code int
	Body string
	Err  error
}

// Error -
func (e RequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s | %s | status code: %d", e.Err.Error(), e.Body, e.Code)
	}
	return fmt.Sprintf("%s | status code: %d", e.Body, e.Code)
}
