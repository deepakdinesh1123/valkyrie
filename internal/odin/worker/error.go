package worker

import "fmt"

type WorkerError struct {
	Type    string
	Message string
}

func (e *WorkerError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
