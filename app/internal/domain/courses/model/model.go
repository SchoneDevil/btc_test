package model

import "time"

type Course struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Buy       float64   `json:"buy"`
	Rub       float64   `json:"rub"`
	CreatedAt time.Time `json:"created_at"`
}
