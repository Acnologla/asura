package main

import (
	"asura/src/database"
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

func main() {
	var migrations = migrate.NewMigrations()

	godotenv.Load()
	db, err := database.Connect(database.GetEnvConfig())
	fmt.Println(err)
	if err != nil {
		return
	}
	if err := migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}

	migrator := migrate.NewMigrator(db, migrations)
	err = migrator.Init(context.Background())
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Created migrations tables")
	}
}
