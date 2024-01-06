package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/slog"
)

const version = "1.0.0"

type config struct {
	port int
	evn  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4800, "API Server port")
	flag.StringVar(&cfg.evn, "env", "development", "Environment (development|staging|production)")
	flag.Parse()
	fmt.Println(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/health-check", app.healthCheckHandler)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.port), Handler: app.routes(), IdleTimeout: time.Minute, ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second, ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError)}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.evn)
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}
