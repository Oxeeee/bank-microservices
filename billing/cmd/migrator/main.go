package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var dsn, migrationsPath string

	flag.StringVar(&dsn, "dsn", "", "conntection to database")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")

	flag.Parse()

	if dsn == "" {
		panic("dsn is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		migrationsPath,
		dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no changes to apply")
			return
		}

		log.Fatal(err)
	}

	log.Println("migrations apply successfully")
}
