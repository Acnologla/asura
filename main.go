package main

import (
	_ "asura/src/commands" // Initialize all commands and put them into an array
	"asura/src/database"
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils/rinha"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
)

func onReady() {
	myself, err := handler.Client.Cache().GetCurrentUser()
	handler.Client.UpdateStatusString("Use j!comandos para ver meus comandos")
	if err == nil {
		telemetry.Info(fmt.Sprintf("%s Started", myself.Username), map[string]string{
			"eventType": "ready",
		})
	}
	go telemetry.MetricUpdate(handler.Client)

}
func main() {
	rand.Seed(time.Now().UnixNano())
	//If it's not in production so it's good to read a ".env" file
	if os.Getenv("PRODUCTION") == "" {
		err := godotenv.Load()
		if err != nil {
			panic("Cannot read the motherfucking envfile")
		}
	}
	rinha.SetTopToken(os.Getenv("TOP_TOKEN"))
	//Initialize datalog services for telemetry of the application
	telemetry.Init()
	database.Init()
	/*rinha.UpdateGaloDB(365948625676795904, func(galo rinha.Galo) (rinha.Galo, error) {
		galo.Cosmetics = append(galo.Cosmetics, 21)
		return galo, nil
	})*/
	fmt.Println("Starting bot...")
	cache := disgord.NewBasicCache()
	client := disgord.New(disgord.Config{
		RejectEvents: []string{disgord.EvtPresenceUpdate, disgord.EvtTypingStart},
		BotToken:     os.Getenv("TOKEN"),
		Cache: &handler.Cache{
			BasicCache: cache,
			Messages:   map[disgord.Snowflake]*handler.Message{},
		},
	})
	defer client.Gateway().StayConnectedUntilInterrupted()
	handler.Client = client
	client.Gateway().MessageCreate(handler.OnMessage)
	client.Gateway().MessageReactionAdd(handler.OnReactionAdd)
	client.Gateway().MessageReactionRemove(handler.OnReactionRemove)
	client.Gateway().InteractionCreate(handler.Interaction)
	client.Gateway().BotReady(onReady)

}
