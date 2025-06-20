package storageerrors

import "errors"

const (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrProductNotFound = errors.New("product not found")
)
