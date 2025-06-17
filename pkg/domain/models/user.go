package models

type User struct {
	UserID   int64
	Email    string
	Passhash []byte
}
