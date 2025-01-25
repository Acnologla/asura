package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
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

func shardsToString(itens []*entities.Item) (str string) {
	for _, item := range itens {
		if item.Type == entities.ShardType {
			rarity := rinha.Rarity(item.ItemID)
			str += fmt.Sprintf("[%d] Shard %s (%s)\n", item.Quantity, rarity.String(), rarity.String())
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

type inventoryButton struct {
	items string
	label string
	emoji string
}

func getEmbedText(button inventoryButton) string {
	if len(button.items) == 0 {
		return fmt.Sprintf("**%s**\n\nNenhum item", button.label)
	}
	return fmt.Sprintf("**%s**\n\n%s", button.label, button.items)
}

var btnOrder = []string{"rooster", "items", "cosmetics", "skins", "keys", "shards"}

func runInventory(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.UserID, "Items", "Galos")
	galos := roostersToString(user.Galos)
	items := itensToString(user.Items)
	skins := skinsToString(user.Items)
	cosmetics := cosmeticsToString(user.Items)
	keys := keysToString(user.Items)
	shards := shardsToString(user.Items)
	buttons := map[string]inventoryButton{
		"rooster":   {galos, "Galos", "🐓"},
		"items":     {items, "Itens", "🎒"},
		"skins":     {skins, "Skins", "👗"},
		"cosmetics": {cosmetics, "Cosmeticos", "🖼️"},
		"keys":      {keys, "Chaves", "🔑"},
		"shards":    {shards, "Shards", "⚪"},
	}
	buttonsComponent := []*disgord.MessageComponent{}
	for _, key := range btnOrder {
		button := buttons[key]
		buttonsComponent = append(buttonsComponent, &disgord.MessageComponent{
			Type:  disgord.MessageComponentButton,
			Label: button.label,
			Emoji: &disgord.Emoji{
				Name: button.emoji,
			},
			Style:    disgord.Primary,
			CustomID: key,
		})
	}
	divided := utils.Chunk(buttonsComponent, 3)
	rows := []*disgord.MessageComponent{
		{
			Type:       disgord.MessageComponentActionRow,
			Components: divided[0],
		},
		{
			Type:       disgord.MessageComponentActionRow,
			Components: divided[1],
		},
	}
	r := &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Components: rows,
			Embeds: []*disgord.Embed{
				{
					Title:       "Inventario",
					Color:       65535,
					Description: getEmbedText(buttons["rooster"]),
				},
			},
		},
	}
	itcID, err := handler.SendInteractionResponse(ctx, itc, r)
	if err != nil {
		return nil
	}
	handler.RegisterHandler(itcID, func(ic *disgord.InteractionCreate) {
		if ic.Member.UserID == user.ID {
			customID := ic.Data.CustomID
			if button, ok := buttons[customID]; ok {
				r.Data.Embeds[0].Description = getEmbedText(button)
				handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackUpdateMessage,
					Data: r.Data,
				})
			}
		}
	}, 200)
	return nil
}
