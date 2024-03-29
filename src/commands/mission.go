package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"fmt"
	"strconv"
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "mission",
		Description: translation.T("MissionHelp", "pt"),
		Run:         runMission,
		Cooldown:    10,
		Category:    handler.Rinha,
	})
}

func runMission(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := itc.Member.User
	galo := database.User.GetUser(ctx, user.ID, "Missions")
	galo.Missions = getMissions(ctx, &galo)
	text := rinha.MissionsToString(user.ID, &galo)
	embed := &disgord.Embed{
		Color:       65535,
		Title:       fmt.Sprintf("Missoes (%d/3)", len(galo.Missions)),
		Description: text,
	}
	if len(galo.Missions) != 3 {
		need := uint64(time.Now().Unix()) - galo.LastMission
		embed.Footer = &disgord.EmbedFooter{
			Text: translation.T("MissionTime", translation.GetLocale(itc), map[string]interface{}{
				"hours":   23 - (need / 60 / 60),
				"minutes": 59 - (need / 60 % 60),
			}),
		}
	}
	components := []*disgord.MessageComponent{}
	if (time.Now().Unix()-int64(galo.TradeMission))/60/60/24 >= 3 {
		for i := range galo.Missions {
			components = append(components, &disgord.MessageComponent{
				Type:     disgord.MessageComponentButton,
				Style:    disgord.Primary,
				Label:    translation.T("ChangeMission", translation.GetLocale(itc), i),
				CustomID: strconv.Itoa(i),
			})
		}
	}
	component := []*disgord.MessageComponent{{
		Type:       disgord.MessageComponentActionRow,
		Components: components,
	}}
	params := &disgord.CreateInteractionResponseData{}
	params.Embeds = []*disgord.Embed{embed}
	if len(components) > 0 {
		params.Components = component
	}
	handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: params,
	})
	if len(components) > 0 {
		handler.RegisterHandler(itc.ID, func(event *disgord.InteractionCreate) {
			if event.Member.User.ID == user.ID {
				i, _ := strconv.Atoi(event.Data.CustomID)
				database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
					if 0 > i || i >= len(u.Missions) || (time.Now().Unix()-int64(u.TradeMission))/60/60/24 < 3 {
						return u
					}
					newMission := rinha.CreateMission()
					newMission.ID = u.Missions[i].ID
					newMission.UserID = user.ID
					database.User.UpdateMissions(ctx, user.ID, &newMission, false)
					u.TradeMission = uint64(time.Now().Unix())
					return u
				}, "Missions")
				handler.DeleteHandler(itc.ID)
				handler.Client.SendInteractionResponse(ctx, event, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: translation.T("TradeMission", translation.GetLocale(event)),
					},
				})
			}
		}, 120)
	}
	return nil
}
