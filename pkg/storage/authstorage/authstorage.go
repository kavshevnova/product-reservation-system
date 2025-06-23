package authstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/kavshevnova/product-reservation-system/pkg/domain/models"
	"github.com/lib/pq"
)

type StorageUsers struct {
	db *sql.DB
}

func NewUsersStorage(db *sql.DB) (*StorageUsers, error) {
	return &StorageUsers{db: db}, nil
}

func (s *StorageUsers) SaveUser(ctx context.Context, email string, passhash []byte) (uid int64, err error) {
	const op = "storage.authstorage.SaveUser"
	const query = "INSERT INTO users (email, passhash) VALUES ($1, $2) RETURNING id"
	var id int64
	err = s.db.QueryRowContext(ctx, query, email, passhash).Scan(&id)
	if err != nil {
		if isDuplicateKeyError(err) {
			return 0, fmt.Errorf("%s: %w", op, models.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *StorageUsers) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.authstorage.User"
	const query = "SELECT id, email, passhash FROM users WHERE email = $1"

	var user models.User
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.Passhash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

// Вспомогательная функция для проверки ошибки дублирования
func isDuplicateKeyError(err error) bool {
	// Для PostgreSQL
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}
