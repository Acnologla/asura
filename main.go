package main

import (
	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
	"asura/src/logs"
	"asura/src/database"
	"os"
	_ "asura/src/commands"
	"asura/src/utils"
	"asura/src/handler"
	"context"
	"fmt"
)

func onReady(session disgord.Session, evt *disgord.Ready) {
	fmt.Println("Bot",evt.User.Username,"Logged!!")
}

func main(){
	if os.Getenv("production") == ""{
		err := godotenv.Load()
		if err != nil { panic("Cannot read the motherfucking envfile") }
	}

	logs.Init()
	database.Init()
	utils.MetricUpdate()

	fmt.Println("Iniciando bot....")

	client := disgord.New(disgord.Config{
        BotToken: os.Getenv("TOKEN"),
	}) 

	client.On(disgord.EvtMessageCreate, handler.OnMessage)
	client.On(disgord.EvtReady, onReady)

	defer client.StayConnectedUntilInterrupted(context.Background())
}