package models

import (
	"errors"
	"time"
)

type Order struct {
	ID         int64
	UserID     int64
	ProductID  int64
	Quantity   int32
	Sum        float32
	Status     string
	Time       time.Time
	PaymentURL string
}

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFound      = errors.New("order not found")
)
