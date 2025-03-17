﻿package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
)

func main() {
	m, err := migrate.New(
		"file://"+"migrations",
		"postgres://postgres:postgres@localhost:5432/grpcauthservice?sslmode=disable")
	if err != nil {
		slog.Error("migration connect err:", err)
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migration to apply")
			return
		}

		slog.Error("failed to up migrations:", err)
		panic(err)
	}

	slog.Info("migration applied successfully")
}
