package cbr

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"app/internal/domain/cbr/storage"
	"app/pkg/client/postgresql"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Client postgresql.Client
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/currencies", h.GetCbrLast)
	router.HandlerFunc(http.MethodPost, "/api/currencies", h.GetCbrHistory)
}

func (h *Handler) GetCbrLast(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cbrStorage := storage.NewCbrStorage(h.Client)
	cbr, err := cbrStorage.SelectLast(context.Background())
	if err != nil {
		log.Println(err)
	}
	res := make(map[string]interface{})
	res["total"] = len(cbr)
	res["history"] = cbr

	r, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(r)
	return
}

func (h *Handler) GetCbrHistory(w http.ResponseWriter, req *http.Request) {

}
