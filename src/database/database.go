package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var Database *bun.DB

type DBConfig struct {
	User     string `json:"DB_USER"`
	Host     string `json:"DB_HOST"`
	Port     string `json:"DB_PORT"`
	Password string `json:"DB_PASS"`
	DbName   string `json:"DB_NAME"`
}

func GetEnvConfig() (config *DBConfig) {
	dbconfig := os.Getenv("DB_CONFIG")
	err := json.Unmarshal([]byte(dbconfig), &config)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Connect(config *DBConfig) (*bun.DB, error) {
	dbConfig := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", config.Host, config.Port)),
		pgdriver.WithUser(config.User),
		pgdriver.WithPassword(config.Password),
		pgdriver.WithDatabase(config.DbName),
		pgdriver.WithInsecure(true),
	)

	sqldb := sql.OpenDB(dbConfig)
	err := sqldb.Ping()

	if err != nil {
		return nil, err
	}
	Database = bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	Database.AddQueryHook(bundebug.NewQueryHook())

	return Database, nil
}
