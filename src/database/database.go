package database

import (
	"asura/src/adapter"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

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

var Client adapter.UserAdapter

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
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	Database.SetMaxOpenConns(maxOpenConns)
	Database.SetMaxIdleConns(maxOpenConns)
	Database.AddQueryHook(bundebug.NewQueryHook())
	Client = adapter.UserAdapterPsql{
		Db: Database,
	}
	return Database, nil
}
