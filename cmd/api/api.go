package main

import (
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", app.healthCheckHandler)
	return mux
}

func (app *application) run() error {

	mux := app.mount()

	srv := http.Server{
		Addr:              app.config.addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
