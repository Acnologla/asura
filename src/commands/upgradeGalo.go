package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/translation"
	"context"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "upgradegalo",
		Description: "Melhore a raridade dos seus galos",
		Run:         runUpgradeGalo,
		Cooldown:    10,
		Category:    handler.Rinha,
	})
}

const SHARDS_TO_UPGRADE = 10

func getShardsByRarity(items []*entities.Item, rarity rinha.Rarity) *entities.Item {

	for _, item := range items {
		if item.Type == entities.ShardType && item.ItemID == int(rarity) {
			return item
		}
	}

	return &entities.Item{
		Type:     entities.ShardType,
		Quantity: 0,
	}
}

func genUpgradeRoosterOptions(user *entities.User) (opts []*disgord.SelectMenuOption) {
	for _, galo := range user.Galos {
		class := rinha.Classes[galo.Type]
		rarity := rinha.GetRarity(galo)
		if rarity != rinha.Mythic {
			opts = append(opts, &disgord.SelectMenuOption{
				Label:       class.Name,
				Value:       galo.ID.String(),
				Description: rarity.String(),
			})
		}
	}

	if len(opts) == 0 {
		opts = append(opts, &disgord.SelectMenuOption{
			Label: "Nenhum galo",
			Value: "nil",
		})
	}
	return
}

var rarityToShard = map[rinha.Rarity]rinha.Rarity{
	rinha.Common:    rinha.Rare,
	rinha.Rare:      rinha.Epic,
	rinha.Epic:      rinha.Legendary,
	rinha.Legendary: rinha.Mythic,
}

func runUpgradeGalo(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.UserID, "Items", "Galos")
	rareShard := getShardsByRarity(user.Items, rinha.Rare)
	epicShard := getShardsByRarity(user.Items, rinha.Epic)
	legendaryShard := getShardsByRarity(user.Items, rinha.Legendary)
	mythicShard := getShardsByRarity(user.Items, rinha.Mythic)
	opts := genUpgradeRoosterOptions(&user)

	r := entities.CreateMsg().
		Embed(&disgord.Embed{
			Title: "Upgrade rooster",
			Color: 65535,
			Footer: &disgord.EmbedFooter{
				Text: translation.T("upgradeGaloWarning", translation.GetLocale(itc)),
			},
			Description: translation.T("upgradeGaloDescription", translation.GetLocale(itc), map[string]int{
				"rareShards":      rareShard.Quantity,
				"epicShards":      epicShard.Quantity,
				"legendaryShards": legendaryShard.Quantity,
				"mythicShards":    mythicShard.Quantity,
				"maxShards":       SHARDS_TO_UPGRADE,
			}),
		}).
		Component(entities.CreateComponent().
			Select(
				translation.T("SelectGalo", translation.GetLocale(itc)),
				"galoUpgradeRarity",
				opts))

	err := handler.Client.SendInteractionResponse(ctx, itc, r.Res())

	if err != nil {
		return nil
	}

	handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
		userIC := ic.Member.User
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
		roosterID := uuid.MustParse(val)
		msg := ""
		var shard rinha.Rarity
		database.User.UpdateUser(ctx, userIC.ID, func(u entities.User) entities.User {
			rooster := rinha.GetRoosterByID(u.Galos, roosterID)
			if rooster == nil {
				return u
			}

			shardRarity, ok := rarityToShard[rinha.GetRarity(rooster)]

			if !ok {
				return u
			}

			shard = shardRarity
			shards := getShardsByRarity(u.Items, shard)

			if shards.Quantity < SHARDS_TO_UPGRADE {
				msg = "notEnoughShards"
				return u
			}

			database.User.RemoveItemQuantity(ctx, u.Items, shards.ID, SHARDS_TO_UPGRADE)

			database.User.UpdateRooster(ctx, &u, rooster.ID, func(r entities.Rooster) entities.Rooster {
				r.EvolvedRarity = int(shard)
				return r
			})

			msg = "upgradeGaloSuccess"
			return u
		}, "Galos", "Items")

		if msg != "" {
			ic.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T(msg, translation.GetLocale(ic), map[string]interface{}{
						"shard":     shard.String(),
						"maxShards": SHARDS_TO_UPGRADE,
					}),
				},
			})
		}
	}, 120)

	return nil
}
