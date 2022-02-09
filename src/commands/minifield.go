package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"strconv"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

var miniFieldEmojis = []string{"0️⃣", "1⃣", "2⃣", "3⃣", "4⃣", "5⃣", "6⃣", "7⃣", "8⃣", "9⃣"}

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "minifield",
		Description: translation.T("MinifieldHelp", "pt"),
		Run:         runMinifield,
		Cooldown:    15,
		Category:    handler.Games,
	})
}

type MinifieldTile struct {
	Type   int
	Hidden bool
}

func makeMinifieldBoard(x, y int) []([]*MinifieldTile) {
	board := []([]*MinifieldTile){}
	for i := 0; i < x; i++ {
		row := []*MinifieldTile{}
		for j := 0; j < y; j++ {
			ty := 0
			if 2 >= utils.RandInt(10) {
				ty = 1
			}
			row = append(row, &MinifieldTile{
				Hidden: true,
				Type:   ty,
			})
		}
		board = append(board, row)
	}
	return board
}

func calcBombs(board []([]*MinifieldTile), i, j int) (numBombs int) {
	for k, row := range board {
		for a, tile := range row {
			if tile.Type == 1 {
				if k == i-1 {
					if a == j+1 || a == j || j-1 == a {
						numBombs++
					}
				}
				if k == i {
					if a == j+1 || j-1 == a {
						numBombs++
					}
				}
				if k == i+1 {
					if a == j+1 || a == j || j-1 == a {
						numBombs++
					}
				}
			}

		}
	}
	return
}

func generateMinifieldBoard(board []([]*MinifieldTile)) []*disgord.MessageComponent {
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
		{
			Type: disgord.MessageComponentActionRow,
		},
		{
			Type: disgord.MessageComponentActionRow,
		},
	}
	for i, row := range board {
		for j, tile := range row {
			button := &disgord.MessageComponent{
				Type:     disgord.MessageComponentButton,
				Style:    disgord.Secondary,
				CustomID: strconv.Itoa(i) + strconv.Itoa(j),
				Label:    "\u200f",
			}
			if !tile.Hidden {
				if tile.Type == 0 {
					button.Style = disgord.Success
					button.Emoji = &disgord.Emoji{
						Name: miniFieldEmojis[calcBombs(board, i, j)],
					}
				} else {
					button.Style = disgord.Danger
					button.Emoji = &disgord.Emoji{
						Name: "💥",
					}
				}
			}
			arrs[i].Components = append(arrs[i].Components, button)

		}
	}
	return arrs
}
func endMinifield(board []([]*MinifieldTile)) {
	for _, row := range board {
		for _, tile := range row {
			tile.Hidden = false
		}
	}
}

func totalCells(board []([]*MinifieldTile)) (total int) {
	for _, row := range board {
		for _, tile := range row {
			if tile.Type == 0 {
				total++
			}
		}
	}
	return
}

func runMinifield(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	board := makeMinifieldBoard(5, 5)
	total := totalCells(board)
	data := generateMinifieldBoard(board)
	handler.Client.SendInteractionResponse(context.Background(), itc, &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.InteractionApplicationCommandCallbackData{
			Content:    translation.T("MinifieldPlay", "pt"),
			Components: data,
		},
	})
	handler.RegisterHandler(itc.ID, func(interaction *disgord.InteractionCreate) {
		if interaction.Member.User.ID != itc.Member.UserID {
			return
		}
		i, _ := strconv.Atoi(string(interaction.Data.CustomID[0]))
		j, _ := strconv.Atoi(string(interaction.Data.CustomID[1]))
		cell := board[i][j]
		if !cell.Hidden {
			return
		}
		cell.Hidden = false
		content := "Clique nos botoes para jogar"
		totalUsed := 0
		for _, row := range board {
			for _, tile := range row {
				if tile.Type == 0 && !tile.Hidden {
					totalUsed++
				}
			}
		}
		if totalUsed == total {
			endMinifield(board)
			content = "Você ganhou"
			handler.DeleteHandler(itc.ID)
		} else if cell.Type == 1 {
			endMinifield(board)
			content = "Você perdeu"
			handler.DeleteHandler(itc.ID)
		}
		newMsg := generateMinifieldBoard(board)
		handler.Client.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackUpdateMessage,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content:    content,
				Components: newMsg,
			},
		})
	}, 60*2)
	return nil
}
