package test

import (
	_ "asura/src/commands"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {

	if os.Getenv("PRODUCTION") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Cannot read the motherfucking envfile")
		}
	}
	m.Run()
}
