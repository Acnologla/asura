package test

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("production") == "" {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Println("Cannot read the motherfucking envfile")
		}
	}
	m.Run()
}
