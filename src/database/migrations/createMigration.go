package main

import (
	"asura/src/database"
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

func main() {
	var migrations = migrate.NewMigrations()

	godotenv.Load()
	db, _ := database.Connect(database.GetEnvConfig())
	if len(os.Args) < 2 {
		return
	}
	if err := migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
	name := os.Args[1]
	migrator := migrate.NewMigrator(db, migrations)
	_, err := migrator.CreateSQLMigrations(context.Background(), name)
	if err != nil {
		fmt.Printf("Migration %s created", name)
	}
}
