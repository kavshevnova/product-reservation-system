package main

import (
	"fmt"
	"github.com/kavshevnova/product-reservation-system/pkg/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	//TODO: сделать grpc файлы
	//TODO: сделать логгер
	//TODO: сделать все остальное
}
