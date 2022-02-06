package main

import (
	"asura/src/cache"
	_ "asura/src/commands"
	"asura/src/database"
	"asura/src/firebase"
	"asura/src/handler"
	"asura/src/telemetry"
	"fmt"
	"log"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
)

var Port string

func initBot() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
	})
	appID := os.Getenv("APP_ID")
	token := os.Getenv("TOKEN")
	firebase.Init()
	handler.Init(appID, token, client)
	defer client.Gateway().StayConnectedUntilInterrupted()
	client.Gateway().BotReady(func() {
		go telemetry.MetricUpdate(client)
		fmt.Println("Logged in")
	})
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
	telemetry.Init()
	cache.Init()
	_, err := database.Connect(database.GetEnvConfig())
	if err != nil {
		log.Fatal(err)
	}
	initBot()
}
