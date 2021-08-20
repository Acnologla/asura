package commands

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"missoes", "daily", "mission", "missão"},
		Run:       runMission,
		Available: true,
		Cooldown:  3,
		Usage:     "j!mission",
		Help:      "Veja suas missoes",
		Category:  1,
	})
}

func runMission(session disgord.Session, msg *disgord.Message, args []string) {
	user := msg.Author
	galo, _ := rinha.GetGaloDB(user.ID)
	text := rinha.MissionsToString(user.ID, galo)
	galo, _ = rinha.GetGaloDB(user.ID)
	embed := &disgord.Embed{
		Color:       65535,
		Title:       fmt.Sprintf("Missoes (%d/3)", len(galo.Missions)),
		Description: text,
	}
	if len(galo.Missions) != 3 {
		need := uint64(time.Now().Unix()) - galo.LastMission
		embed.Footer = &disgord.EmbedFooter{
			Text: fmt.Sprintf("Faltam %d horas e %d minutos para voce receber uma nova missão", 23-(need/60/60), 59-(need/60%60)),
		}
	}
	components := []*disgord.MessageComponent{}
	if (time.Now().Unix()-int64(galo.MissionTrade))/60/60/24 >= 3 {
		for i := range galo.Missions {
			components = append(components, &disgord.MessageComponent{
				Type:     disgord.MessageComponentButton,
				Style:    disgord.Primary,
				Label:    fmt.Sprintf("Alterar missão %d", i),
				CustomID: strconv.Itoa(i),
			})
		}
	}
	component := []*disgord.MessageComponent{{
		Type:       disgord.MessageComponentActionRow,
		Components: components,
	}}
	params := disgord.CreateMessageParams{}
	params.Embed = embed
	if len(components) > 0 {
		params.Components = component
	}
	message, err := msg.Reply(context.Background(), session, params)
	if err == nil {
		if len(components) > 0 {
			handler.RegisterBHandler(message, func(event *disgord.InteractionCreate) {
				if event.Member.User.ID == msg.Author.ID {
					i, _ := strconv.Atoi(event.Data.CustomID)
					rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
						if 0 > i || i >= len(galo.Missions) {
							return galo, nil
						}
						if (time.Now().Unix()-int64(galo.MissionTrade))/60/60/24 < 3 {
							return galo, nil
						}
						galo.Missions[i] = rinha.CreateMission()
						galo.MissionTrade = uint64(time.Now().Unix())
						return galo, nil
					})
					handler.DeleteBHandler(message)
					session.SendInteractionResponse(context.Background(), event, &disgord.InteractionResponse{
						Type: disgord.ChannelMessageWithSource,
						Data: &disgord.InteractionApplicationCommandCallbackData{
							Content: "Você trocou sua missão com sucesso!",
						},
					})
				}
			}, 120)
		}
	}
}
