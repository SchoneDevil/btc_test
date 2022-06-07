package model

import "time"

type CbrToBtc struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}
