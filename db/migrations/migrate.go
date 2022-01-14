package main

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

func main() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
