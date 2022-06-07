package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"

	cbrh "app/internal/handlers/cbr"
	"app/internal/handlers/cbrtobtc"
	"app/internal/handlers/courses"
	"app/internal/schedule"
	"app/pkg/client/postgresql"
	"app/pkg/shutdown"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

var wg sync.WaitGroup

func (a *App) Preload() {
	log.Println("Init preload")
	sch := schedule.Scheduler{Client: a.pgClient}

	go GetDataCbr(sch)
	go GetDataBtc(sch)

	wg.Add(2)
	wg.Wait()

}

func GetDataCbr(sch schedule.Scheduler) {
	defer wg.Done()
	sch.GetCbr()

}

func GetDataBtc(sch schedule.Scheduler) {
	defer wg.Done()
	sch.GetBtcUsdt()
}

type App struct {
	router     *httprouter.Router
	httpServer *http.Server
	pgClient   *pgxpool.Pool
}

func NewApp() (App, error) {
	router := httprouter.New()

	//router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	//router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	pgConfig := postgresql.NewPgConfig(
		os.Getenv("EDITOR"),
		os.Getenv("EDITOR"),
		os.Getenv("EDITOR"),
		os.Getenv("EDITOR"),
		os.Getenv("EDITOR"),
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		log.Fatal(err)
	}

	sch := schedule.Scheduler{Client: pgClient}
	go sch.Start()

	coursesHandler := courses.Handler{Client: pgClient}
	coursesHandler.Register(router)

	cbrHandler := cbrh.Handler{Client: pgClient}
	cbrHandler.Register(router)

	cbrtobtcHandler := cbrtobtc.Handler{Client: pgClient}
	cbrtobtcHandler.Register(router)

	return App{
		router:   router,
		pgClient: pgClient,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	var server *http.Server
	var listener net.Listener

	var err error

	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", "0.0.0.0", "8080"))
	if err != nil {
		log.Fatal(err)
	}

	server = &http.Server{
		Handler:      a.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	log.Println("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Println("server shutdown")
		default:
			log.Fatal(err)
		}
	}
}
