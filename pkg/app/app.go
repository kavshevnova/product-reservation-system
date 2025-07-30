package app

import (
	grpcapp "github.com/kavshevnova/product-reservation-system/pkg/app/grpc"
	"github.com/kavshevnova/product-reservation-system/pkg/services/auth"
	"github.com/kavshevnova/product-reservation-system/pkg/services/shop"
	"github.com/kavshevnova/product-reservation-system/pkg/storages/authstorage"
	"github.com/kavshevnova/product-reservation-system/pkg/storages/shopstorage"
	"log/slog"
)

type App struct {
	GRPCsrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcport int,
	storagepath string,
) *App {

	addr := "redis:6380"
	password := ""
	db := 0

	storageAuth, err := authstorage.NewUsersStorage(addr, password, db)
	if err != nil {
		panic(err)
	}

	storageShop, err := shopstorage.NewShopStorage(storagepath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storageAuth, storageAuth)
	shopService := shop.New(log, storageShop, storageShop)

	grpcApp := grpcapp.New(log, authService, shopService, grpcport)

	return &App{
		GRPCsrv: grpcApp,
	}
}
