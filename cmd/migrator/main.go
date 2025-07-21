package main

import (
	"context"
	"database/sql"
	"flag"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	_ "github.com/pressly/goose/v3"
	"log"
)

func main() {
	var dsn, migrationsDir, command string
	flag.StringVar(&dsn, "DataSoursName", "host=localhost port=5433 user=postgres password=mysecretpassword dbname=postgres sslmode=disable", "path to postgres database")
	flag.StringVar(&migrationsDir, "migrations-dir", "./migrations", "path to migrations directory")
	flag.StringVar(&command, "command", "up", "goose command (up, down, status)")
	flag.Parse()

	if dsn == "" {
		panic("--storage-path is required")
	}
	if migrationsDir == "" {
		panic("--migrations-path is required")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database:%v", err)
	}
	defer db.Close()

	goose.SetDialect("postgres")

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
