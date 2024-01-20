package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"khomeapi/internal/data"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

const version = "1.0.0"

type config struct {
	port int
	evn  string
	db   struct {
		dns          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4800, "API Server port")
	flag.StringVar(&cfg.evn, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dns, "db-dns", "postgres://admin:admin@localhost:5432/linux?sslmode=disable", "Postgres DNS")

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice that the default values we're using are the ones we discussed above?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.Parse()
	fmt.Println(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDb(cfg)
	if err != nil {
		fmt.Println("data base error")
	}
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModel(db),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/health-check", app.healthCheckHandler)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.port), Handler: app.routes(), IdleTimeout: time.Minute, ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second, ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError)}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.evn)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}

func openDb(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dns)

	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetConnMaxLifetime(cfg.db.maxIdleTime)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	fmt.Println(err, "error")
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
