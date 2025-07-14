package authgrpc

import (
	"context"
	"errors"
	authv1 "github.com/kavshevnova/product-reservation-system/gen/go/auth"
	"github.com/kavshevnova/product-reservation-system/pkg/domain/models"
	"github.com/kavshevnova/product-reservation-system/pkg/services/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error)
	LoginUser(ctx context.Context, email, password string) (success bool, err error)
}

type AuthServerAPI struct {
	authv1.UnimplementedAuthServiceServer
	auth Auth
}

// функция RegisterAuthServerAPI регистрирует обработчик для того чтобы он обрабатывал те запросы которые поступают в grpc сервер
func RegisterAuthServerAPI(grpcServer *grpc.Server, auth Auth) {
	authv1.RegisterAuthServiceServer(grpcServer, &AuthServerAPI{auth: auth})
}

func (a *AuthServerAPI) Register(ctx context.Context, request *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := ValidateRegister(request); err != nil {
		return nil, err
	}
	userID, err := a.auth.RegisterNewUser(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		if errors.Is(err, models.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &authv1.RegisterResponse{UserId: userID}, nil
}

func (a *AuthServerAPI) Login(ctx context.Context, request *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := ValidateLogin(request); err != nil {
		return nil, err
	}
	success, err := a.auth.LoginUser(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &authv1.LoginResponse{Success: success}, nil
}

func (a *AuthServerAPI) mustEmbedUnimplementedAuthServiceServer() {}

func ValidateRegister(request *authv1.RegisterRequest) error {
	if request.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "missing email")
	}
	if request.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "missing password")
	}
	return nil
}

func ValidateLogin(request *authv1.LoginRequest) error {
	if request.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "missing email")
	}
	if request.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "missing password")
	}
	return nil
}
