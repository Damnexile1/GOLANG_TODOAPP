package core_errors

import "errors"

var (
	ErrNotFound        = errors.New("model not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrConflict        = errors.New("conflict")
)
