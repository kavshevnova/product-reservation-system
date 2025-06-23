package models

import "errors"

type Product struct {
	ProductID int64
	Name      string
	Price     float32
	Stock     int32
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrNotEnoughStock  = errors.New("not enough stock")
)
