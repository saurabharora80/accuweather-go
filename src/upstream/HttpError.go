package upstream

import (
	"fmt"
)

type HttpError struct {
	StatusCode int
	Path       string
	Method     string
	Details    []byte
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%s %s failed with %d => %s", e.Method, e.Path, e.StatusCode, e.Details)
}
