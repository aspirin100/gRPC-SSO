package main

import (
	"errors"
	"flag"
	"fmt"
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
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("sqlite3://%s", storagePath),
	)
	if err != nil {
		panic(err)
	}

	err = mInstance.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no changes")
		} else {
			panic(err)
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
