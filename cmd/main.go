package main

import (
	"fmt"
	"github.com/kavshevnova/product-reservation-system/pkg/config"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger := SetUpLogger(cfg.Env)
	logger.Info("Стартуем", slog.Any("Config", cfg))
	//TODO: сделать grpc файлы
	//TODO: сделать все остальное
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
