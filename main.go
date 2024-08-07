package main

import (
	"asura/src/cache"
	_ "asura/src/commands"
	"asura/src/database"
	"asura/src/events"
	"asura/src/firebase"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/telemetry"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
)

func initBot() {
	rand.Seed(time.Now().UnixNano())
	client := disgord.New(disgord.Config{
		BotToken:     os.Getenv("TOKEN"),
		RejectEvents: []string{disgord.EvtPresenceUpdate, disgord.EvtTypingStart},
	})
	appID := os.Getenv("APP_ID")
	token := os.Getenv("TOKEN")
	handler.Init(appID, token, client)
	defer client.Gateway().StayConnectedUntilInterrupted()
	client.Gateway().BotReady(func() {
		client.UpdateStatusString("/help | Caso não apareça os comandos readicione o bot no servidor")
		go telemetry.MetricUpdate(client)
		fmt.Println("Logged in")
	})
	client.Gateway().MessageCreate(events.HandleMessage)
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
	rinha.SetTopToken(os.Getenv("TOP_TOKEN"))
	telemetry.Init()
	cache.Init()
	firebase.Init()
	_, err := database.Connect(database.GetEnvConfig())
	if err != nil {
		log.Fatal(err)
	}
	initBot()
}
