package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"strconv"

	"github.com/andersfylling/disgord"
)

var miniFieldEmojis = []string{"0️⃣", "1⃣", "2⃣", "3⃣", "4⃣", "5⃣", "6⃣", "7⃣", "8⃣", "9⃣"}

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"campominado", "minifield"},
		Run:       runMinifield,
		Available: true,
		Cooldown:  10,
		Usage:     "j!minifield",
		Help:      "Jogue campo minado",
		Category:  3,
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
			if 2 >= utils.RandInt(9) {
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

func generateMinifieldBoard(board []([]*MinifieldTile)) *disgord.CreateMessageParams {
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
	params := &disgord.CreateMessageParams{}
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
	params.Components = arrs
	return params
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

func runMinifield(session disgord.Session, msg *disgord.Message, args []string) {
	board := makeMinifieldBoard(5, 5)
	total := totalCells(board)
	text := generateMinifieldBoard(board)
	text.Content = "Clique nos botoes para jogar"
	message, err := msg.Reply(context.Background(), session, text)
	if err != nil {
		return
	}
	handler.RegisterBHandler(message, func(interaction *disgord.InteractionCreate) {
		if interaction.Member.User.ID != msg.Author.ID {
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
			handler.DeleteBHandler(message)
		} else if cell.Type == 1 {
			endMinifield(board)
			content = "Você perdeu"
			handler.DeleteBHandler(message)
		}
		newMsg := generateMinifieldBoard(board)
		session.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
			Type: disgord.UpdateMessage,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content:    content,
				Components: newMsg.Components,
			},
		})
	}, 60*2)
}
