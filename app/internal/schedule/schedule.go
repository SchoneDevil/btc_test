package schedule

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"app/internal/clients/cbr"
	"app/internal/clients/kucoin"
	model2 "app/internal/domain/cbr/model"
	storage2 "app/internal/domain/cbr/storage"
	model3 "app/internal/domain/cbrtobtc/model"
	storage3 "app/internal/domain/cbrtobtc/storage"
	"app/internal/domain/courses/model"
	"app/internal/domain/courses/storage"
	"app/pkg/client/postgresql"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	Client postgresql.Client
}

func (sch Scheduler) Start() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(10).Second().Do(sch.GetBtcUsdt)
	s.Every(24).Hour().Do(sch.GetCbr)

	s.StartAsync()
	s.StartBlocking()
}

func (sch Scheduler) GetBtcUsdt() {
	ku := kucoin.IKucoin{}
	btcusdt := ku.GetBtcUsdt()

	courseStorage := storage.NewCourseStorage(sch.Client)

	buy, _ := strconv.ParseFloat(btcusdt.Data.Buy, 64)

	cbrStorage := storage2.NewCbrStorage(sch.Client)
	rur, err := cbrStorage.SelectUsd(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	rub := buy * rur.Value
	err = courseStorage.Insert(context.Background(), model.Course{
		Symbol:    btcusdt.Data.Symbol,
		Buy:       buy,
		Rub:       rub,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("%s", err)
	}

	ctbStorage := storage3.NewCbrToBtcStorage(sch.Client)

	cbrList, _ := cbrStorage.SelectLast(context.Background())
	for _, v := range cbrList {
		err := ctbStorage.InsertOrUpdate(context.Background(), model3.CbrToBtc{
			Name:      v.CharCode,
			Value:     v.Value * rub,
			CreatedAt: time.Now(),
		})
		if err != nil {
			log.Printf("%s", err)
		}
	}

}

func (sch Scheduler) GetCbr() {
	cbrInterface := cbr.ICbr{}
	cbrCourses := cbrInterface.GetCbr()

	cbrStorage := storage2.NewCbrStorage(sch.Client)

	for _, v := range cbrCourses {
		cbrStorage.Insert(context.Background(), model2.Cbr{
			CharCode:  v.CharCode,
			Name:      v.Name,
			Value:     v.Value,
			CreatedAt: time.Now(),
		})
	}

}
