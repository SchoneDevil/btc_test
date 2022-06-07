package courses

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"app/internal/domain/courses/storage"
	"app/internal/handlers"
	"app/pkg/client/postgresql"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Client postgresql.Client
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/btcusdt", h.GetBtc)
	router.HandlerFunc(http.MethodPost, "/api/btcusdt", h.PostBtc)
}

func (h *Handler) GetBtc(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courseStorage := storage.NewCourseStorage(h.Client)
	course, err := courseStorage.SelectLast(context.Background())
	if err != nil {
		log.Println(err)
	}

	r, _ := json.Marshal(course)

	w.WriteHeader(http.StatusOK)
	w.Write(r)
	return
}

//POST-запрос - историю с фильтрами по дате и времени и пагинацией

func (h *Handler) PostBtc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	var f handlers.PostFilter
	err := decoder.Decode(&f)
	if err != nil {
		fmt.Println(err)
	}
	courseStorage := storage.NewCourseStorage(h.Client)
	courses, err := courseStorage.Select(context.Background(), f)

	res := make(map[string]interface{})
	res["total"] = len(courses)
	res["history"] = courses

	r, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(r)

	w.WriteHeader(http.StatusOK)

}
