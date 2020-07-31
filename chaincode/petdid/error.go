package main

type Code int

type baseError struct {
	msg  string
}

func (e *baseError) Error() string {
	return e.msg
}

func newBaseError(msg string) *baseError {
	return &baseError{msg}
}