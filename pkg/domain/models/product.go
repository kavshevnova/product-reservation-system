package models

import "errors"

type Product struct {
	ProductID int64   `db:"product_id"`
	Name      string  `db:"name"`
	Price     float32 `db:"price"`
	Stock     int32   `db:"stock"`
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrNotEnoughStock  = errors.New("not enough stock")
)
