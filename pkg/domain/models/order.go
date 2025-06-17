package models

import "time"

type Order struct {
	ID        int64
	UserID    int64
	ProductID int64
	Quantity  int32
	Sum       float32
	Time      time.Time
}
