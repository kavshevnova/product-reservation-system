package models

import "errors"

type User struct {
	UserID   int64
	Email    string
	Passhash []byte
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)
