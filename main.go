package main

import (
	"asura/src/cache"
	_ "asura/src/commands"
	"asura/src/database"
	"asura/src/firebase"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/telemetry"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
)

func initBot() {
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
	client.Gateway().MessageCreate(func(s disgord.Session, h *disgord.MessageCreate) {
		go func() {
			msg := h.Message
			if !msg.Author.Bot {
				if msg.GuildID != 0 {
					for _, user := range msg.Mentions {
						if user.ID.String() == appID {
							msg.Reply(context.Background(), s, "Use /help para ver meus comandos\nCaso meus comandos não aparecam me readicione no servidor com este link:\nhttps://discordapp.com/oauth2/authorize?client_id=470684281102925844&scope=applications.commands%%20bot&permissions=8")
							break
						}
					}
				}
			}
		}()
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
