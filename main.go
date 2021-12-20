package main

import (
	"asura/src/server"
	"fmt"
	"os"

	_ "asura/src/commands"
	"asura/src/handler"

	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

var Port string

func main() {
	if os.Getenv("PRODUCTION") == "" {
		err := godotenv.Load()
		if err != nil {
			panic("Cannot read the motherfucking envfile")
		}
	}
	server.Init(os.Getenv("PUBLIC_KEY"))
	Port = os.Getenv("PORT")
	appID := os.Getenv("APP_ID")
	token := os.Getenv("TOKEN")
	handler.Init(appID, token)
	fasthttp.ListenAndServe(":"+Port, server.Handler)

	fmt.Println("server started")
}
