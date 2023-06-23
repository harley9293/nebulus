package errors

import "fmt"

type NebulusError struct {
	message string
	param   []any
}

func New(message string) *NebulusError {
	return &NebulusError{message: message}
}

func (e *NebulusError) Error() string {
	return fmt.Sprintf(e.message, e.param...)
}

func (e *NebulusError) Fill(param ...any) *NebulusError {
	newError := New(e.message)
	newError.param = param
	return newError
}
