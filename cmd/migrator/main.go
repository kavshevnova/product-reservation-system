package main

import (
	"context"
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	_ "github.com/pressly/goose/v3"
	"log"
)

func main() {
	var storagePath, migrationsDir, command string
	flag.StringVar(&storagePath, "storage-path", "./storage/storage.db", "path to SQLitedatabase")
	flag.StringVar(&migrationsDir, "migrations-dir", "./migrations", "path to migrations directory")
	flag.StringVar(&command, "command", "up", "goose command (up, down, status)")
	flag.Parse()

	if storagePath == "" {
		panic("--storage-path is required")
	}
	if migrationsDir == "" {
		panic("--migrations-path is required")
	}
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		log.Fatalf("failed to open database:%v", err)
	}
	defer db.Close()

	goose.SetDialect("sqlite3")

	goose.SetVerbose(true)

	//выполняем команду
	if err := goose.RunContext(context.Background(), command, db, migrationsDir); err != nil {
		if err == goose.ErrNoNextVersion {
			log.Println("no migrations to be applied")
			return
		}
		log.Fatalf("failed to run goose:%v", err)
	}

	log.Println("migrations complete")
}
