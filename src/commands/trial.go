package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"asura/src/utils"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "trial",
		Category:    handler.Rinha,
		Description: "Faça desafios pro seu galo ficar mais forte",
		Run:         runTrial,
		Cooldown:    7,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "status",
				Description: "Veja seu status atual",
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "battle",
				Description: "Va para o proximo desafio",
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "all",
				Description: "Veja todos seus desafios completados",
			},
		),
	})
}

var trialDamageMultiplier = 4

const MAX_TRIALS = 5
const TRIAL_MIN_LEVEL = 15

func getTrial(trials []*entities.Trial, rooster int) *entities.Trial {
	for _, trial := range trials {
		if trial.Rooster == rooster {
			return trial
		}
	}
	return nil
}

func resMultiplier(resets int) int {
	if resets == 0 {
		return 0
	}
	return 1
}

func runTrial(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.UserID, "Items", "Galos", "Trials")
	command := itc.Data.Options[0].Name
	rooster := rinha.GetEquippedGalo(&user)
	trial := getTrial(user.Trials, rooster.Type)
	if trial == nil {
		database.User.InsertTrial(ctx, itc.Member.UserID, &entities.Trial{
			Rooster: rooster.Type,
		})
		user = database.User.GetUser(ctx, itc.Member.UserID, "Items", "Galos", "Trials")
		trial = getTrial(user.Trials, rooster.Type)
	}
	galoSprite := rinha.GetGaloImage(rooster, user.Items)
	switch command {
	case "status":
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Trials",
						Description: fmt.Sprintf("Seu galo esta na trial **%d/%d**\nBonus de dano **%d%%**\n\nNa trial final voce ganha uma lootbox de acordo com a raridade do seu galo\nUse `/trial battle` para batalhar", trial.Win, MAX_TRIALS, trialDamageMultiplier*trial.Win),
						Thumbnail: &disgord.EmbedThumbnail{
							URL: galoSprite,
						},
						Color: 65535,
					},
				},
			},
		}
	case "all":
		text := ""
		for _, trial := range user.Trials {
			rooster := rinha.Classes[trial.Rooster]
			text += fmt.Sprintf("**%s** - %d/%d\n", rooster.Name, trial.Win, MAX_TRIALS)
		}
		avatar, _ := itc.Member.User.AvatarURL(1024, true)
		if text == "" {
			text = "Nenhuma trial completa"
		}
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Trials",
						Description: text,
						Color:       65535,
						Footer: &disgord.EmbedFooter{
							Text:    "Use /trial battle para batalhar",
							IconURL: avatar,
						},
					},
				},
			},
		}
	case "battle":
		if trial.Win >= MAX_TRIALS {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Voce ja completou todas as trials",
				},
			}
		}
		roosterLevel := rinha.CalcLevel(rooster.Xp)
		if roosterLevel < TRIAL_MIN_LEVEL {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: fmt.Sprintf("Seu galo precisa ser no minimo level **%d** para participar das trials", TRIAL_MIN_LEVEL),
				},
			}
		}
		discordUser := itc.Member.User
		authorRinha := isInRinha(ctx, discordUser)
		if authorRinha != "" {
			return rinhaMessage(discordUser.Username, authorRinha)
		}
		lockEvent(ctx, discordUser.ID, "Desafio trial")
		defer unlockEvent(ctx, discordUser.ID)
		class := rinha.Classes[trial.Rooster]
		xpMultiplier := 1
		if class.Rarity == rinha.Epic {
			xpMultiplier = 3 + rooster.Resets
		}
		if class.Rarity > rinha.Epic {
			xpMultiplier = (int(class.Rarity) + rooster.Resets) * 4
		}
		xp := rinha.CalcXP(15+(trial.Win*7)) * xpMultiplier
		level := rinha.CalcLevel(xp)
		galoAdv := &entities.Rooster{
			Type:    rinha.GetRandByType(class.Rarity),
			Xp:      xp,
			Evolved: class.Rarity > rinha.Epic || level >= 40,
			Resets:  rooster.Resets + resMultiplier(rooster.Resets)*(trial.Win*2),
		}

		userAdv := entities.User{
			Galos: []*entities.Rooster{galoAdv},
		}

		if class.Rarity > rinha.Epic {
			atbs := (user.UserXp / 210) * ((trial.Win / 2) + 1)
			if class.Rarity == rinha.Legendary {
				atbs = atbs / 2
			}
			healthAtb := float64(atbs) * 0.55
			userAdv.Attributes = [6]int{int(healthAtb) + 500, atbs / 10, atbs, atbs / 9, atbs / 8, 50}
		}

		winner, _ := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
			GaloAuthor: &user,
			GaloAdv:    &userAdv,
			IDs:        [2]disgord.Snowflake{discordUser.ID},

			AuthorName:  rinha.GetName(discordUser.Username, *rooster),
			AdvName:     "Boss " + rinha.Classes[galoAdv.Type].Name,
			AuthorLevel: rinha.CalcLevel(rooster.Xp),
			AdvLevel:    rinha.CalcLevel(galoAdv.Xp),
			NoItems:     false,
		}, false)

		if winner == -1 {
			return nil
		}

		ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))

		go completeMission(ctx, &user, galoAdv, winner == 0, itc, "trial")
		if winner == 0 {
			database.User.AddTrialWin(ctx, trial)
			if trial.Win >= MAX_TRIALS {
				lootbox := rinha.GetTrialLootbox(class.Rarity)
				database.User.InsertItem(ctx, itc.Member.UserID, user.Items, lootbox, entities.LootboxType)
				lootboxName := rinha.LootNames[lootbox]
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{{
						Color:       65535,
						Title:       "Trial",
						Description: fmt.Sprintf("Parabens voce chegou no nivel maximo de trial para o galo **%s**, e ganhou uma lootbox **%s**", class.Name, lootboxName),
					}},
				})
				n := 0
				for _, t := range user.Trials {
					if t.Win >= MAX_TRIALS {
						n++
					}
				}
				if n >= 3 {
					completeAchievement(ctx, itc, 15)
				}
			} else {
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{{
						Color:       16776960,
						Title:       "Trial",
						Description: fmt.Sprintf("Voce venceu a trial **%d/%d**\nGanhou **%d%%** de dano a mais para seu galo", trial.Win, MAX_TRIALS, trialDamageMultiplier*trial.Win),
					}},
				})
			}
			return nil
		}
		ch.CreateMessage(&disgord.CreateMessage{
			Embeds: []*disgord.Embed{{
				Color:       16711680,
				Title:       "Trial",
				Description: "Parabens voce perdeu",
			}},
		})
	}

	return nil
}
