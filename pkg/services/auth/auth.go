package auth

import "errors"

var (
	//Credentials - Реквизиты для входа
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)
