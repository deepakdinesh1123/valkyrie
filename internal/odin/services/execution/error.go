package execution

import "fmt"

type ExecutionServiceError struct {
	Type    string
	Message string
}

func (e *ExecutionServiceError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

type TemplateError struct {
	Message string
}

func (e *TemplateError) Error() string {
	return e.Message
}
