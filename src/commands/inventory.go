package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "inventory",
		Description: "Veja seu inventario",
		Run:         runInventory,
		Cooldown:    7,
		Category:    handler.Profile,
		Aliases:     []string{"inventario"},
	})
}

func roostersToString(roosters []*entities.Rooster) (str string) {
	for _, rooster := range roosters {
		class := rinha.Classes[rooster.Type]
		str += fmt.Sprintf("%s (%s)\n", class.Name, class.Rarity.String())
	}
	return
}

func itensToString(itens []*entities.Item) (str string) {
	for _, item := range itens {
		if item.Type == entities.NormalType {
			i := rinha.Items[item.ItemID]
			str += fmt.Sprintf("[%d] %s (%s)\n", item.Quantity, i.Name, rinha.LevelToString(i.Level))
		}
	}
	return
}

func skinsToString(itens []*entities.Item) (str string) {
	for _, item := range itens {
		if item.Type == entities.CosmeticType {
			i := rinha.Cosmetics[item.ItemID]
			if i.Type == rinha.Skin {
				str += fmt.Sprintf("[%d] %s (%s)\n", item.Quantity, i.Name, i.Rarity.String())
			}
		}
	}
	return
}

func cosmeticsToString(itens []*entities.Item) (str string) {
	for _, item := range itens {
		if item.Type == entities.CosmeticType {
			i := rinha.Cosmetics[item.ItemID]
			if i.Type != rinha.Skin {
				str += fmt.Sprintf("[%d] %s (%s)\n", item.Quantity, i.Name, i.Rarity.String())
			}
		}
	}
	return
}

func keysToString(itens []*entities.Item) (str string) {
	for _, item := range itens {
		if item.Type == entities.KeyType {
			str += fmt.Sprintf("[%d] Chave (%s) \n", item.Quantity, rinha.Rarity(item.ItemID).String())
		}
	}
	return
}

func runInventory(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.UserID, "Items", "Galos")
	galos := roostersToString(user.Galos)
	items := itensToString(user.Items)
	skins := skinsToString(user.Items)
	cosmetics := cosmeticsToString(user.Items)
	keys := keysToString(user.Items)
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Title: "Inventario",
					Color: 65535,
					Fields: []*disgord.EmbedField{
						{
							Name:   "Galos",
							Value:  galos,
							Inline: true,
						},
						{
							Name:   "Itens",
							Value:  items,
							Inline: false,
						},
						{
							Name:   "Skins",
							Value:  skins,
							Inline: false,
						},
						{
							Name:   "Cosmeticos",
							Value:  cosmetics,
							Inline: false,
						},
						{
							Name:   "Chaves",
							Value:  keys,
							Inline: false,
						},
					},
					//	Description: fmt.Sprintf("Galos:\n%s\nItens:\n%s\nCosmeticos:\n%s", galos, items, cosmetics),
				},
			},
		},
	}
}
