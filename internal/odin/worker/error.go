package worker

import "fmt"

type WorkerError struct {
	Type    string
	Message string
}

func (e *WorkerError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *WorkerError) Unwrap() error {
	return e
}

type WorkerInfoNotFoundError struct {
}

func (e *WorkerInfoNotFoundError) Error() string {
	return "worker info not found"
}
