﻿package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	authgrpc "grpcAuth/internal/grpc/auth"
	"log/slog"
	"net"
)

type App struct {
	log         *slog.Logger
	authService authgrpc.Auth
	gRPCServer  *grpc.Server
	port        int
}

func New(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authService)

	return &App{log, authService, gRPCServer, port}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(slog.String("op", op),
		slog.Int("port", a.port))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting gRPC server", "err", a.port)

	err = a.gRPCServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	log := a.log.With(slog.String("op", op))
	log.Info("stopping gRPC server", "port", a.port)
	a.gRPCServer.GracefulStop()
}
