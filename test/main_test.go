package test

import (
	"testing"
	"github.com/joho/godotenv"
	"os"
	"log"
)

func TestMain(m *testing.M){
	if os.Getenv("production") == ""{
		err := godotenv.Load("../.env")
		if err != nil { log.Println("Cannot read the motherfucking envfile") }
	}
	m.Run()
}