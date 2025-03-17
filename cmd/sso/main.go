﻿package main

import (
	"fmt"
	"grpcAuth/internal/app"
	"grpcAuth/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config

	cfg := config.MustLoad()

	// TODO: init logger

	log := setupLogger(cfg.Env)

	log.Info("starting application",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
	)

	// TODO: init app

	// TODO: run GRPC-server

	// postgres://postgres:postgres@localhost:5432/grpcauthservice?sslmode=disable
	// postgres://postgres:postgres@localhost:5432/grpcauthservice?sslmode=disable
	postgresConString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresConfig.Username,
		cfg.PostgresConfig.Password,
		cfg.PostgresConfig.Host,
		cfg.PostgresConfig.Port,
		cfg.PostgresConfig.Database)

	application := app.New(log, cfg.GRPC.Port, postgresConString, cfg.TokenTTL)

	go func() {
		err := application.GRPCSrv.Run()
		if err != nil {
			panic(err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

// go run cmd/sso/main.go --config=./config/local.yaml
