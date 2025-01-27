package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath,
		migrationsPath string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")

	flag.Parse()
	validateFlags(storagePath, migrationsPath)

	mInstance, err := migrate.New(
		"file://"+migrationsPath,
		"sqlite3://"+storagePath,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = mInstance.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no changes")
		} else {
			log.Fatal(err)
		}
	}
}

func validateFlags(storagePath, migrationsPath string) {
	if storagePath == "" {
		panic("storage path should be not empty")
	}

	if migrationsPath == "" {
		panic("migrations path should be not empty")
	}
}
