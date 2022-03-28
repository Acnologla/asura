package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"fmt"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "equip",
		Description: translation.T("RunEquip", "pt"),
		Run:         runEquip,
		Cooldown:    15,
		Category:    handler.Profile,
	})
}

func genEquipOptions(user *entities.User) (opts []*disgord.SelectMenuOption) {
	for _, galo := range user.Galos {
		if !galo.Equip {
			class := rinha.Classes[galo.Type]
			opts = append(opts, &disgord.SelectMenuOption{
				Label:       class.Name,
				Value:       galo.ID.String(),
				Description: fmt.Sprintf("Equipar galo %s (level %d)", class.Name, rinha.CalcLevel(galo.Xp)),
			})
		}
	}

	if len(opts) == 0 {
		opts = append(opts, &disgord.SelectMenuOption{
			Label:       "Nenhum galo",
			Value:       "nil",
			Description: "Nenhum galo para equipar",
		})
	}
	return
}

func runEquip(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	galo := database.User.GetUser(ctx, itc.Member.UserID, "Galos")
	optsGalos := genEquipOptions(&galo)
	r := entities.CreateMsg().
		Embed(&disgord.Embed{
			Title: "Equip",
			Color: 65535,
		}).
		Component(entities.CreateComponent().
			Select(
				translation.T("EquipGaloPlaceholder", translation.GetLocale(itc)),
				"galoEquip",
				optsGalos))
	handler.Client.SendInteractionResponse(ctx, itc, r.Res())
	handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
		userIC := ic.Member.User
		name := ic.Data.CustomID
		if userIC.ID != itc.Member.UserID {
			return
		}
		if len(ic.Data.Values) == 0 {
			return
		}
		val := ic.Data.Values[0]
		if val == "nil" {
			return
		}
		itemID := uuid.MustParse(val)
		msg := ""
		database.User.UpdateUser(ctx, userIC.ID, func(u entities.User) entities.User {
			if isInRinha(ctx, userIC) != "" {
				msg = "IsInRinha"
				return u
			}
			if name == "galoEquip" {
				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					database.User.UpdateRooster(ctx, &u, itemID, func(r2 entities.Rooster) entities.Rooster {
						if r2.Type != 0 {
							r.Equip = false
							r2.Equip = true
						}
						return r2
					})
					return r
				})
				msg = "EquipGalo"
			}
			return u
		}, "Galos")
		if msg != "" {
			handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T(msg, translation.GetLocale(ic)),
				},
			})
		}
	}, 120)
	return nil
}
