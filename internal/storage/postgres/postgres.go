package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"grpcAuth/internal/domain/models"
	"grpcAuth/internal/storage"
	"log/slog"
)

type Storage struct {
	db *sql.DB
}

func New(connectionString string) (*Storage, error) {
	const op = "postgresql.New"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error(op, "\t", err, "ConnectionString:", connectionString)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		slog.Error(op, "failed to ping db", "\t", err, "ConnectionString:", connectionString)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context,
	email string,
	passHash []byte) (uid int64, err error) {
	const op = "postgresql.SaveUser"

	log := slog.With(slog.String("op", op))

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRowContext(ctx, email, passHash).Scan(&uid)
	if err != nil {
		log.Error("failed to save user", "email", email, "error", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user saved", "id", uid, "email", email)
	return uid, nil
}

func (s *Storage) User(ctx context.Context, email string) (user models.User, err error) {
	const op = "postgresql.User"

	log := slog.With(slog.String("op", op))

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email=$1")
	if err != nil {
		log.Error("failed to prepare query", "stmt", stmt, "error", err)
		return user, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, email)
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successfully extracted", "id", user.ID, "email", user.Email)
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {
	const op = "postgresql.IsAdmin"
	log := slog.With(slog.String("op", op))

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id=$1")
	if err != nil {
		log.Error("failed to prepare query", "stmt", stmt, "error", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, err, storage.ErrAppNotFound)
		}
		log.Error("failed to query user", "user_id", userID, "error", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (app models.App, err error) {
	const op = "postgresql.App"
	log := slog.With(slog.String("op", op))
	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id=$1")
	if err != nil {
		log.Error("failed to prepare query", "stmt", stmt, "error", err)
	}
	row := stmt.QueryRowContext(ctx, appID)
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app, fmt.Errorf("%s: %w", op, err, storage.ErrAppNotFound)
		}
		log.Error("failed to query app", "app_id", appID, "error", err)
	}

	return app, nil
}
