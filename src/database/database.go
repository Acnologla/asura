package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DBConfig struct {
	User     string
	Host     string
	Port     string
	Password string
	DbName   string
}

func GetEnvConfig() *DBConfig {
	return &DBConfig{
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		DbName:   os.Getenv("DB_NAME"),
	}
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
	return bun.NewDB(sqldb, pgdialect.New()), nil
}
