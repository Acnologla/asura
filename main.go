package main

import (
	"asura/server"
	"fmt"
	"os"

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
	server.PublicKey = os.Getenv("PUBLIC_KEY")
	Port = os.Getenv("PORT")
	fasthttp.ListenAndServe(":"+Port, server.Handler)
	fmt.Println("server started")
}
