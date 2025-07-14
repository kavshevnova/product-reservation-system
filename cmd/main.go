package main

import (
	"fmt"
	"github.com/kavshevnova/product-reservation-system/pkg/app"
	"github.com/kavshevnova/product-reservation-system/pkg/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger := SetUpLogger(cfg.Env)
	logger.Info("Стартуем", slog.Any("Config", cfg))
	application := app.New(logger, cfg.GRPC.Port, cfg.StoragePath)
	go func() {
		application.GRPCsrv.MustStart()
		logger.Info("starting gRPC server")
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	logger.Info("starting graceful shutdown")
	application.GRPCsrv.Stop()
	logger.Info("graceful shutdown complete")
}

func SetUpLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
