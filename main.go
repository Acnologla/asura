package main

import (
	_ "asura/src/commands" // Initialize all commands and put them into an array
	"asura/src/database"
	"asura/src/handler"
	"asura/src/telemetry"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"time"
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
	// If it's not in production so it's good to read a ".env" file
	if os.Getenv("PRODUCTION") == "" {
		err := godotenv.Load()
		if err != nil {
			panic("Cannot read the motherfucking envfile")
		}
	}
	// Initialize datalog services for telemetry of the application
	telemetry.Init()
	database.Init()
	fmt.Println("Starting bot...")
	cache := disgord.NewCacheLFUImmutable(0, 0, 0, 0)
	client := disgord.New(disgord.Config{
		RejectEvents: []string{disgord.EvtPresenceUpdate, disgord.EvtTypingStart},
		BotToken:     os.Getenv("TOKEN"),
		Cache: &handler.Cache{
			CacheLFUImmutable: cache.(*disgord.CacheLFUImmutable),
			Messages:          map[disgord.Snowflake]*handler.Message{},
		},
	})
	defer client.Gateway().StayConnectedUntilInterrupted()
	handler.Client = client
	client.Gateway().MessageCreate(handler.OnMessage)
	client.Gateway().MessageReactionAdd(handler.OnReactionAdd)
	client.Gateway().MessageReactionRemove(handler.OnReactionRemove)
	client.Gateway().BotReady(onReady)

}
