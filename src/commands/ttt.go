package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "ttt",
		Description: translation.T("TTTHelp", "pt"),
		Run:         runTTT,
		Cooldown:    10,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Type:        disgord.OptionTypeUser,
			Required:    true,
			Name:        "user",
			Description: translation.T("TTTDesc", "pt"),
		}),
		Category: handler.Games,
	})
}

func board(tiles [9]int) *disgord.InteractionApplicationCommandCallbackData {
	arrs := []*disgord.MessageComponent{
		{
			Type: disgord.MessageComponentActionRow,
		},
		{
			Type: disgord.MessageComponentActionRow,
		},
		{
			Type: disgord.MessageComponentActionRow,
		},
	}
	params := &disgord.InteractionApplicationCommandCallbackData{
		Content: "Clique nos botoes para jogar",
	}
	for i, tile := range tiles {
		var arrField int = i / 3
		button := &disgord.MessageComponent{
			Type:     disgord.MessageComponentButton,
			Style:    disgord.Secondary,
			CustomID: strconv.Itoa(i),
			Label:    "\u200f",
		}
		if tile == 1 {
			button.Style = disgord.Primary
			button.Emoji = &disgord.Emoji{
				Name: "❌",
			}
			button.Disabled = true
		} else if tile == 2 {
			button.Style = disgord.Danger
			button.Emoji = &disgord.Emoji{
				Name: "⭕",
			}
			button.Disabled = true
		}
		arrs[arrField].Components = append(arrs[arrField].Components, button)
	}
	params.Components = arrs
	return params
}

func win(tiles [9]int) int {
	for i := 1; i < 3; i++ {
		if (tiles[0] == i && tiles[1] == i && tiles[2] == i) || (tiles[3] == i && tiles[4] == i && tiles[5] == i) || (tiles[6] == i && tiles[7] == i && tiles[8] == i) || (tiles[0] == i && tiles[3] == i && tiles[6] == i) || (tiles[1] == i && tiles[4] == i && tiles[7] == i) || (tiles[2] == i && tiles[5] == i && tiles[8] == i) || (tiles[0] == i && tiles[4] == i && tiles[8] == i) || (tiles[2] == i && tiles[4] == i && tiles[6] == i) {
			return i
		}
	}
	if tiles[0] != 0 && tiles[1] != 0 && tiles[2] != 0 && tiles[3] != 0 && tiles[4] != 0 && tiles[5] != 0 && tiles[6] != 0 && tiles[7] != 0 && tiles[8] != 0 {
		return 0
	}
	return 3
}

func playTTT(tile int, tiles *[9]int, turn int) bool {
	if tiles[tile] != 0 {
		return false
	}
	tiles[tile] = turn
	return true
}
func runTTT(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	user := utils.GetUser(itc, 0)
	if user.ID == itc.Member.User.ID || user.Bot {
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: "Invalid user",
			},
		}
	}
	tiles := [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	turn := 1
	handler.Client.SendInteractionResponse(context.Background(), itc, &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: board(tiles),
	})
	defer handler.RegisterHandler(itc.ID, func(interaction *disgord.InteractionCreate) {
		turnUser := user
		letter := ":x:"
		if turn == 1 {
			turnUser = itc.Member.User
			letter = ":o:"
		}
		if turnUser.ID == interaction.Member.User.ID {
			tile, _ := strconv.Atoi(interaction.Data.CustomID)
			played := playTTT(tile, &tiles, turn)
			if played {
				if turn == 2 {
					turn--
				} else {
					turn++
				}
				winner := win(tiles)
				playType := ":hash:"
				if winner != 3 {
					playType = ":crown:"
					if letter == ":x:" {
						letter = ":o:"
					} else {
						letter = ":x:"
					}
					if winner == 0 {
						letter = ":no_good:"
					}
					handler.DeleteHandler(itc.ID)
				}
				handler.Client.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
					Type: disgord.InteractionCallbackUpdateMessage,
					Data: &disgord.InteractionApplicationCommandCallbackData{
						Content:    fmt.Sprintf("%s%s", playType, letter),
						Components: board(tiles).Components,
					},
				})
			}
		}
	}, 60*5)
	return nil
}
