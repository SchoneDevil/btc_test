package handlers

import "github.com/julienschmidt/httprouter"

type Handler interface {
	Register(router *httprouter.Router)
}

type PostFilter struct {
	DateStart  string `json:"date_start,omitempty"`
	DateFinish string `json:"date_finish,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}
