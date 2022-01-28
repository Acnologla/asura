package main

import (
	"asura/src/cache"
	_ "asura/src/commands"
	"asura/src/firebase"
	"asura/src/handler"
	"asura/src/server"
	"fmt"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

var Port string

func runTestVersion() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
	})
	appID := os.Getenv("APP_ID")
	token := os.Getenv("TOKEN")
	firebase.Init()
	handler.Init(appID, token, client)
	defer client.Gateway().StayConnectedUntilInterrupted()
	client.Gateway().BotReady(func() {
		fmt.Println("Logged in")
	})
	//TODO worker pool
	client.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		handler.InteractionChannel <- h
	})
}

func main() {
	if os.Getenv("PRODUCTION") == "" {
		err := godotenv.Load()
		if err != nil {
			panic("Cannot read the motherfucking envfile")
		}
	}
	cache.Init()
	if os.Getenv("PRODUCTION") == "" {
		runTestVersion()
		return
	}

	server.Init(os.Getenv("PUBLIC_KEY"))
	Port = os.Getenv("PORT")
	fasthttp.ListenAndServe(":"+Port, server.Handler)
	fmt.Println("server started")

}
