package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/kavshevnova/product-reservation-system/pkg/domain/models"
	storageerrors "github.com/kavshevnova/product-reservation-system/pkg/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type Auth struct {
	log         *slog.Logger
	usrsaver    UserSaver
	usrprovider UserProvider
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passhash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

var (
	//Credentials - Реквизиты для входа
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

func New(
	log *slog.Logger,
	usrsaver UserSaver,
	usprovider UserProvider,
) *Auth {
	return &Auth{
		log:         log,
		usrsaver:    usrsaver,
		usrprovider: usprovider,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("operation", op), slog.String("email", email))

	log.Info("Register new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", slog.StringValue(err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrsaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			a.log.Warn("User already exists", slog.StringValue(err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("Failed to save user", slog.StringValue(err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user registered")
	return id, nil
}

func (a *Auth) LoginUser(ctx context.Context, email, password string) (success bool, err error) {
	const op = "auth.LoginUser"

	log := a.log.With(slog.String("operation", op), slog.String("email", email))

	log.Info("attempting to login user")

	usr, err := a.usrprovider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storageerrors.ErrUserNotFound) {
			a.log.Warn("User not found", slog.StringValue(err.Error()))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("Failed to get user", slog.StringValue(err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(usr.Passhash, []byte(password)); err != nil {
		a.log.Warn("Invalid credentials", slog.StringValue(err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	log.Info("user logged in")
	return true, nil
}
