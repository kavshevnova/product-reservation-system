package shop

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrNotEnoughStock  = errors.New("not enough stock")
)
