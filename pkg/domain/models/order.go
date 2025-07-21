package models

import (
	"errors"
	"time"
)

type Order struct {
	ID         int64     `db:"order_id"`
	UserID     int64     `db:"user_id"`
	ProductID  int64     `db:"product_id"`
	Quantity   int32     `db:"quantity"`
	Sum        float32   `db:"sum"`
	Status     string    `db:"status"`
	Time       time.Time `db:"time"`
	PaymentURL string
}

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFound      = errors.New("order not found")
)
