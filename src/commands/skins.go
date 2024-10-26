package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"context"
	"fmt"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "skins",
		Category:    handler.Profile,
		Description: translation.T("SkinsHelp", "pt"),
		Run:         runSkins,
		Cooldown:    10,
	})
}

func genSkinOptions(skins []*rinha.Cosmetic, items []*entities.Item) []*disgord.SelectMenuOption {
	opts := []*disgord.SelectMenuOption{}
	for i, cosmetic := range items {
		skin := skins[i]
		opts = append(opts, &disgord.SelectMenuOption{
			Label:       skin.Name,
			Value:       cosmetic.ID.String(),
			Description: fmt.Sprintf("[%s] - %s (Galo %s)", skin.Rarity.String(), skin.Name, rinha.Classes[skin.Extra].Name),
		})
	}
	if len(opts) >= 25 {
		opts = opts[:25]
	}
	return opts
}

const SKINS_CHUNK = 25

func genSkinRows(skins []*rinha.Cosmetic, items []*entities.Item) []*disgord.MessageComponent {
	chunkedSkins := utils.Chunk(skins, SKINS_CHUNK)
	chunkedItems := utils.Chunk(items, SKINS_CHUNK)
	rows := []*disgord.MessageComponent{}

	if len(chunkedItems) > 5 {
		chunkedItems = chunkedItems[:5]
	}

	for i, chunk := range chunkedItems {
		opts := genSkinOptions(chunkedSkins[i], chunk)
		rows = append(rows, &disgord.MessageComponent{
			Type: disgord.MessageComponentActionRow,
			Components: []*disgord.MessageComponent{
				{
					MaxValues:   1,
					Type:        disgord.MessageComponentSelectMenu,
					Options:     opts,
					CustomID:    fmt.Sprintf("skins_%d", i),
					Placeholder: "Selecione a skin",
				},
			},
		})
	}
	return rows
}

func skinsToText(skins []*rinha.Cosmetic, items []*entities.Item) string {
	text := ""
	for i, cosmetic := range items {
		if cosmetic.Equip {
			skin := skins[i]
			galo := rinha.Classes[skin.Extra]
			text += fmt.Sprintf("[**%s**] - %s (Galo %s)\n", skin.Rarity.String(), skin.Name, galo.Name)
		}
	}
	return text
}

func runSkins(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	galo := database.User.GetUser(ctx, itc.Member.UserID, "Items")
	skins, items := rinha.GetCosmeticsByTypes(galo.Items, rinha.Skin)
	if len(skins) == 0 {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("SkinsNoSkins", translation.GetLocale(itc)),
			},
		}
	}
	components := genSkinRows(skins, items)
	itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{{
				Color:       65535,
				Title:       "Skins",
				Description: skinsToText(skins, items),
			}},
			Components: components,
		},
	})
	handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
		u := ic.Member.User.ID
		if u != itc.Member.UserID {
			return
		}
		if len(ic.Data.Values) == 0 {
			return
		}
		id := ic.Data.Values[0]
		database.User.UpdateUser(ctx, u, func(u entities.User) entities.User {
			database.User.UpdateItem(ctx, &u, uuid.MustParse(id), func(i entities.Item) entities.Item {
				i.Equip = !i.Equip
				return i
			})
			return u
		}, "Items")
		galo := database.User.GetUser(ctx, itc.Member.UserID, "Items")
		skins, items := rinha.GetCosmeticsByTypes(galo.Items, rinha.Skin)

		ic.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackUpdateMessage,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{{
					Color:       65535,
					Title:       "Skins",
					Description: skinsToText(skins, items),
				}},
				Components: components,
			},
		})

	}, 120)
	return nil
}
