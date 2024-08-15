package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"context"
	"fmt"
	"strings"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "upgrades",
		Description: translation.T("UpgradesHelp", "pt"),
		Run:         runUpgrades,
		Cooldown:    12,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Type:        disgord.OptionTypeNumber,
			Name:        "upgrade_id",
			Description: translation.T("UpgradesId", "pt"),
			MinValue:    0,
			MaxValue:    4,
		}),
		Category: handler.Rinha,
	})
}

func runUpgrades(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := itc.Member.User
	galo := database.User.GetUser(ctx, user.ID)
	if len(itc.Data.Options) == 0 {
		desc := fmt.Sprintf("Upgrade Xp: **%d/%d**", galo.UserXp, rinha.CalcUserXp(&galo))
		if len(galo.Upgrades) > 0 {
			desc += "\n\nUpgrades:\n"
			upgrades := rinha.Upgrade{
				Childs: rinha.Upgrades,
			}
			for i, upgrade := range galo.Upgrades {
				upgrades = upgrades.Childs[upgrade]
				v := strings.Repeat("-", i*5)
				desc += fmt.Sprintf("%s%s\n%s%s\n", v, upgrades.Name, v, upgrades.Value)
			}
		}
		if rinha.HavePoint(&galo) {
			desc += translation.T("UpgradeDesc", translation.GetLocale(itc), rinha.UpgradesToString(&galo))
		}
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{{
					Title:       "Upgrades",
					Color:       65535,
					Description: desc,
				}},
			},
		}
	}
	i := int(itc.Data.Options[0].Value.(float64))

	if !rinha.HavePoint(&galo) {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("NoUpgrade", translation.GetLocale(itc)),
			},
		}
	}
	upgrades := rinha.GetCurrentUpgrade(&galo)
	if 0 > i || i >= len(upgrades.Childs) {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("InvalidUpgradeID", translation.GetLocale(itc)),
			},
		}
	}
	database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
		u.Upgrades = append(u.Upgrades, i)
		return u
	})
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: translation.T("NewUpgrade", translation.GetLocale(itc), upgrades.Childs[i].Name),
		},
	}

}
