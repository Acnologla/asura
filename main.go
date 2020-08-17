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
	"strconv"
	"time"
)

func onReady(session disgord.Session, evt *disgord.Ready) {
	telemetry.Info(fmt.Sprintf("%s Started", evt.User.Username), map[string]string{
		"eventType": "ready",
	})
	go telemetry.MetricUpdate(handler.Client)

}

func onGuildDelete(session disgord.Session, evt *disgord.GuildDelete) {
	return
	guild, err := handler.Client.GetGuild(context.Background(), evt.UnavailableGuild.ID)
	if err != nil {
		fmt.Println(err)
	}
	telemetry.Warn(fmt.Sprintf("Leaved from %s", guild.Name), map[string]string{
		"id":        strconv.FormatUint(uint64(evt.UnavailableGuild.ID), 10),
		"eventType": "leave",
	})

}

func onGuildCreate(session disgord.Session, evt *disgord.GuildCreate) {
	return
	guild := evt.Guild
	bot, err := session.GetCurrentUser(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	member, err := handler.Client.GetMember(context.Background(), guild.ID, bot.ID, disgord.IgnoreCache)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.Until(member.JoinedAt.Time).Seconds())
	if 2 > time.Until(member.JoinedAt.Time).Seconds() {
		telemetry.Warn(fmt.Sprintf("Joined in  %s", guild.Name), map[string]string{
			"id":        strconv.FormatUint(uint64(guild.ID), 10),
			"eventType": "join",
		})
	}
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
	client.On(disgord.EvtGuildCreate, onGuildCreate)
	client.On(disgord.EvtGuildDelete, onGuildDelete)
	client.StayConnectedUntilInterrupted(context.Background())

	fmt.Println("Good bye!")
}
