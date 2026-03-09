package entity

import "time"

type Item struct {
	ID        string
	Name      string
	Price     float64
	CreatedAt time.Time
}
