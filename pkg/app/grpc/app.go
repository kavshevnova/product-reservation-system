package grpcapp

import (
	"fmt"
	"github.com/kavshevnova/product-reservation-system/pkg/grpc/authgrpc"
	"github.com/kavshevnova/product-reservation-system/pkg/grpc/shopgrpc"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	logger *slog.Logger
	grpc   *grpc.Server
	port   int
}

func New(
	logger *slog.Logger,
	authService authgrpc.Auth,
	shopService shopgrpc.Shop,
	port int) *App {
	grpcServer := grpc.NewServer()
	//регистрируем оба сервиса на одном сервере
	authgrpc.RegisterAuthServerAPI(grpcServer, authService)
	shopgrpc.RegisterShopServerAPI(grpcServer, shopService)

	return &App{
		logger: logger,
		grpc:   grpcServer,
		port:   port,
	}
}

func (a *App) MustStart() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcApp.Run"
	log := a.logger.With(
		slog.String("operation", op), slog.Int("port", a.port))

	log.Info("starting grpc app")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("starting grpc server on port", slog.String("address", lis.Addr().String()))

	if err := a.grpc.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcApp.Stop"
	log := a.logger.With(slog.String("operation", op), slog.Int("port", a.port))

	a.grpc.GracefulStop()
}
