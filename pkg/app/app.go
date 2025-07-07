package app

import (
	grpcapp "github.com/kavshevnova/product-reservation-system/pkg/app/grpc"
	"github.com/kavshevnova/product-reservation-system/pkg/services/auth"
	"github.com/kavshevnova/product-reservation-system/pkg/services/shop"
	"github.com/kavshevnova/product-reservation-system/pkg/storage/authstorage"
	"github.com/kavshevnova/product-reservation-system/pkg/storage/shopstorage"
	"log/slog"
)

type App struct {
	GRPCsrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcport int,
	authStoragePath string,
	shopStoragePath string,
) *App {

	storageauth, err := authstorage.NewUsersStorage(authStoragePath)
	if err != nil {
		panic(err)
	}
	storageshop, err := shopstorage.NewShopStorage(shopStoragePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storageauth, storageauth)
	shopService := shop.New(log, storageshop, storageshop)

	grpcApp := grpcapp.New(log, authService, shopService, grpcport)

	return &App{
		GRPCsrv: grpcApp,
	}
}
