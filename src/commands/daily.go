package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "daily",
		Description: translation.T("DailyHelp", "pt"),
		Run:         runDaily,
		Cooldown:    9,
		Category:    handler.Rinha,
	})
}

func runDaily(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	galo := database.User.GetUser(itc.Member.UserID)
	topGGCalc := (uint64(time.Now().Unix()) - galo.Daily) / 60 / 60 / 12
	voted := rinha.HasVoted(itc.Member.UserID)
	if voted && topGGCalc >= 1 {
		strike := 0
		money := 0
		xp := 0
		database.User.UpdateUser(itc.Member.UserID, func(u entities.User) entities.User {
			u.Daily = uint64(time.Now().Unix())
			if topGGCalc >= 2 {
				u.DailyStrikes = 0
			}
			money = 40 + u.DailyStrikes/4
			xp = 60 + u.DailyStrikes
			u.DailyStrikes++
			u.Money += money
			database.User.UpdateEquippedRooster(u, func(r entities.Rooster) entities.Rooster {
				r.Xp += xp
				return r
			})
			strike = u.DailyStrikes
			return u
		}, "Galos")
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Color: 65535,
						Title: "Daily",
						Description: translation.T("DailyMessage", translation.GetLocale(itc), map[string]interface{}{
							"money":   money,
							"xp":      xp,
							"strikes": strike,
						}),
					},
				},
			},
		}
	} else {
		need := uint64(time.Now().Unix()) - galo.Daily
		if topGGCalc >= 1 && !voted {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("VoteMessage", translation.GetLocale(itc)),
				},
			}
		} else {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T("TimeMessage", translation.GetLocale(itc), map[string]interface{}{
						"minutes": 59 - (need / 60 % 60),
						"hours":   11 - (need / 60 / 60),
					}),
				},
			}
		}
	}
}
