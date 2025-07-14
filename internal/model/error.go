package models

import "fmt"

type HTTPError struct {
	StatusCode int
	Msg        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error %d: %s", e.StatusCode, e.Msg)
}
