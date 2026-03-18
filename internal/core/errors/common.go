package core_errors

import "errors"

var (
	ErrNotFound        = errors.New("user not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrConflict        = errors.New("conflict")
)
