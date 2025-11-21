package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/EmotionlessDev/avito-tech-internship/internal/application"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

func main() {
	// Init config
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("PR_POSTGRES_DSN"), "PostgreSQL DSN")
	flag.Parse()

	// Init logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Init DB
	db, err := openDB(cfg)
	if err != nil {
		logger.Error("cannot connect to db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	// Init application via constructor
	application := application.New(cfg, logger, db)

	// Create http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      application.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("starting server",
		slog.Int("port", cfg.port),
		slog.String("env", cfg.env),
	)

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("cannot start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
