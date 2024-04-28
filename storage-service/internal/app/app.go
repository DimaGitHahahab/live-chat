package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"storage-service/internal/repository"
	"storage-service/internal/server"
	"storage-service/internal/service"
	"storage-service/pkg/config"
	"storage-service/pkg/signal"

	pgxLogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type App struct {
	config  *config.Config
	sigQuit chan os.Signal
	server  *server.Server
	ctx     context.Context
}

func New(cfg *config.Config) *App {
	sigQuit := signal.GetShutdownChannel()

	ctx := context.Background()

	repo := setupRepo(ctx, cfg)

	cache := setupCache(cfg)

	storageService := service.NewStorage(ctx, repo, cache, strings.Split(cfg.KafkaBrokers, ","), cfg.KafkaTopic, cfg.KafkaGroup)

	srv := server.New(storageService, cfg.HttpPort)

	return &App{
		config:  cfg,
		sigQuit: sigQuit,
		server:  srv,
		ctx:     ctx,
	}
}

func (a *App) Run() {
	go func() {
		log.Infoln("Starting server on port ", a.config.HttpPort)
		if err := a.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("Failed to start server: ", err)
		}
	}()

	<-a.sigQuit
	log.Infoln("Gracefully shutting down server")

	ctx, cancel := context.WithTimeout(a.ctx, 2*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalln("Failed to shutdown the server gracefully: ", err)
	}

	log.Infoln("Server shutdown is successful")
}

func setupRepo(ctx context.Context, cfg *config.Config) repository.Repository {
	pool, err := setupPgxPool(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return repository.New(pool)
}

func setupPgxPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(cfg.PostgresUrl)
	if err != nil {
		log.Errorln("failed to parse pgxpool pgxConfig: ", err)
		return nil, err
	}

	pgxConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxLogrus.NewLogger(log.StandardLogger()),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Errorln("failed to create new pool: ", err)
		return nil, err
	}

	log.Infoln("Pgx pool initialization successful")
	return pool, nil
}

func setupCache(cfg *config.Config) repository.Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return repository.NewCache(client)
}
