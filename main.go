package main

import (
	"asura/src/cache"
	"asura/src/handler"
	"asura/src/server"
	"context"
	"fmt"
	"os"

	_ "asura/src/commands"

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
	handler.Init(appID, token, client)
	defer client.Gateway().StayConnectedUntilInterrupted()
	client.Gateway().BotReady(func() {
		fmt.Println("Logged in")
	})
	client.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		if h.Type == disgord.InteractionApplicationCommand {
			response := server.ExecuteInteraction(h)
			s.SendInteractionResponse(context.Background(), h, response)
		}
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
