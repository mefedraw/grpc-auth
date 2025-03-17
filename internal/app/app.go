package app

import (
	grpcapp "grpcAuth/internal/app/grpc"
	"grpcAuth/internal/services/auth"
	"grpcAuth/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	connectionString string,
	tokenTTL time.Duration) *App {
	// TODO: init storage
	storage, err := postgres.New(connectionString)
	if err != nil {
		return nil
	}

	// TODO: init auth service
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{grpcApp}
}
