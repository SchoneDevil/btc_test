package model

import "time"

type Cbr struct {
	ID        string    `json:"id"`
	CharCode  string    `json:"char_code"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}
