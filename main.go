package main

import (
	_ "asura/src/commands" // Initialize all commands and put them into an array
	"asura/src/database"
	"asura/src/handler"
	"asura/src/telemetry"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
	"os"
)

func onReady(session disgord.Session, evt *disgord.Ready) {
	handler.Client.UpdateStatusString("Use j!comandos para ver meus comandos")
	telemetry.Info(fmt.Sprintf("%s Started", evt.User.Username), map[string]string{
		"eventType": "ready",
	})
	go telemetry.MetricUpdate(handler.Client)

}

func main() {

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

	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
	})

	handler.Client = client

	client.On(disgord.EvtMessageCreate, handler.OnMessage)
	client.On(disgord.EvtMessageUpdate, handler.OnMessageUpdate)
	client.On(disgord.EvtMessageReactionAdd, handler.OnReactionAdd)
	client.On(disgord.EvtMessageReactionRemove, handler.OnReactionRemove)
	client.On(disgord.EvtReady, onReady)
	client.StayConnectedUntilInterrupted(context.Background())

	fmt.Println("Good bye!")
}
