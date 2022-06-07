package cbrtobtc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"app/internal/domain/cbrtobtc/storage"
	"app/pkg/client/postgresql"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Client postgresql.Client
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/latest", h.GetListBtcToCbr)
}

func (h *Handler) GetListBtcToCbr(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cbrToBtcStorage := storage.NewCbrToBtcStorage(h.Client)
	cbr, err := cbrToBtcStorage.Latest(context.Background())
	if err != nil {
		log.Println(err)
	}
	res := make(map[string]float64)
	fmt.Println(cbr)
	for _, v := range cbr {
		res[v.Name] = v.Value
	}
	fmt.Println(res)
	r, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(r)
	return
}
