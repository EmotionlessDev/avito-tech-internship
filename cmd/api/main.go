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

	"github.com/EmotionlessDev/avito-tech-internship/internal/config"
	prHttp "github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest/delivery/http"
	prCreate "github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest/service/create"
	prStorage "github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest/storage/pullrequest"
	prTeamStorage "github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest/storage/team"
	prUserStorage "github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest/storage/user"
	teamHttp "github.com/EmotionlessDev/avito-tech-internship/internal/domain/team/delivery/http"
	teamAdd "github.com/EmotionlessDev/avito-tech-internship/internal/domain/team/service/add"
	teamGet "github.com/EmotionlessDev/avito-tech-internship/internal/domain/team/service/get"
	teamStorage "github.com/EmotionlessDev/avito-tech-internship/internal/domain/team/storage/team"
	teamUserStorage "github.com/EmotionlessDev/avito-tech-internship/internal/domain/team/storage/user"
	userHttp "github.com/EmotionlessDev/avito-tech-internship/internal/domain/users/delivery/http"
	userUpdate "github.com/EmotionlessDev/avito-tech-internship/internal/domain/users/service/update"
	userStorage "github.com/EmotionlessDev/avito-tech-internship/internal/domain/users/storage/user"

	_ "github.com/lib/pq"
)

func main() {
	// Init config
	cfg := config.New(0, "", "")

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.DSN, "dsn", os.Getenv("PR_POSTGRES_DSN"), "PostgreSQL DSN")
	flag.Parse()

	// Init logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Init DB
	db, err := openDB(cfg)
	if err != nil {
		logger.Error("cannot connect to db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Error("error closing db", slog.String("error", err.Error()))
		}
	}()

	// Init storage
	teamStorage := teamStorage.NewStorage()
	teamUserStorage := teamUserStorage.NewStorage()

	userStorage := userStorage.NewStorage()

	prStorage := prStorage.NewStorage()
	prUserStorage := prUserStorage.NewStorage()
	prTeamStorage := prTeamStorage.NewStorage()

	// Init services
	teamAddService := teamAdd.NewService(teamStorage, teamUserStorage, db)
	teamGetService := teamGet.NewService(teamStorage, teamUserStorage, db)

	userUpdateService := userUpdate.NewService(userStorage, db)

	prCreateService := prCreate.NewService(db, prStorage, prUserStorage, prTeamStorage)

	// Init Handlers
	teamHandler := teamHttp.NewHandler(teamAddService, teamGetService)
	userHandler := userHttp.NewHandler(userUpdateService)
	prHandler := prHttp.NewHandler(prCreateService)

	// Init serveMux
	mux := http.NewServeMux()

	// Map Routes

	// Team routes
	mux.HandleFunc("/team/add", teamHandler.AddTeam)
	mux.HandleFunc("/team/get", teamHandler.GetTeam)

	// User routes
	mux.HandleFunc("/users/setIsActive", userHandler.SetUserActive)

	// PR routes
	mux.HandleFunc("/pullrequest/create", prHandler.CreatePR)

	// Create http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("starting server",
		slog.Int("port", cfg.Port),
		slog.String("env", cfg.Env),
	)

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("cannot start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func openDB(cfg config.ConfigProvider) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDBDSN())
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
