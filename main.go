/*
    __   ____  _  _  ____   __
   / _\ / ___)/ )( \(  _ \ / _\
  /    \\___ \) \/ ( )   //    \
  \_/\_/(____/\____/(__\_)\_/\_/

  This is the source code of the bot
  "Asura" made in Golang and manteined by
  Chiyoku and Acnologia
  Copyright 2020

*/

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
	fmt.Println("Bot", evt.User.Username, "Logged!!")
}

func main() {

	// If it's not in production so it's good to read a ".env" file
	if os.Getenv("production") == "" {
		err := godotenv.Load()
		if err != nil {
			panic("Cannot read the motherfucking envfile")
		}
	}

	// Initialize datalog services for telemetry of the application
	telemetry.Init()
	database.Init()
	telemetry.MetricUpdate()

	fmt.Println("Starting bot...")

	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
	})

	handler.Client = client

	client.On(disgord.EvtMessageCreate, handler.OnMessage)
	client.On(disgord.EvtMessageUpdate, handler.OnMessageUpdate)
	client.On(disgord.EvtReady, onReady)

	client.StayConnectedUntilInterrupted(context.Background())

	fmt.Println("Good bye!")
}
