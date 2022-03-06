package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "2048",
		Description: translation.T("2048Help", "pt"),
		Run:         run2048,
		Cooldown:    30,
		Category:    handler.Games,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "size",
			Type:        disgord.OptionTypeNumber,
			MinValue:    3,
			MaxValue:    6,
			Description: translation.T("2048Size", "pt"),
		}),
	})
}

var arrows = []string{"➡", "⬇", "⬆", "⬅"}
var numberEmojis = [10]string{":zero:", ":one:", ":two:", ":three:", ":four:", ":five:", ":six:", ":seven:", ":eight:", ":nine:"}
var emojis2048 = map[int]string{
	0:    "<:0:744682229023768737>",
	2:    "<:2:744682213047795794>",
	4:    "<:4:744682221452918905>",
	8:    "<:8:744682233855475713>",
	16:   "<:16:744682222145110147>",
	32:   "<:32:744682243276013650>",
	64:   "<:64:744682218538008686>",
	128:  "<:128:744682220559794287>",
	256:  "<:256:744682236544155699>",
	512:  "<:512:744682226482020482>",
	1024: "<:1024:744682235164229653>",
	2048: "<:2048:744682245113249853>",
	4096: "<:4096:750621450255204354>",
	8192: "<:8192:750622586899005461>",
}

func slideLeft(board []([]int), points *int) {
	for i := range board {
		s := 0
		for j := 1; j < len(board[i]); j++ {
			if board[i][j] != 0 {
				for k := j; k > s; k-- {
					if board[i][k-1] == 0 {
						board[i][k-1] = board[i][k]
						board[i][k] = 0
					} else if board[i][k-1] == board[i][k] {
						*points += board[i][k] * 2
						board[i][k-1] += board[i][k]
						board[i][k] = 0
						s = k
						break
					} else {
						break
					}
				}
			}
		}
	}
}

func rotateBoard(board []([]int), c bool) []([]int) {
	var rotatedBoard = utils.MakeBoard(len(board), len(board))
	for i, row := range board {
		for j := range row {
			if c {
				rotatedBoard[i][j] = board[j][len(board)-i-1]
			} else {
				rotatedBoard[i][j] = board[len(board)-j-1][i]
			}
		}
	}
	return rotatedBoard
}

func draw2048Board(board []([]int)) (text string) {
	for _, row := range board {
		for _, tile := range row {
			text += emojis2048[tile]
		}
		text += "\n"
	}
	return
}

func drawPoints(points int) (text string) {
	for _, number := range strconv.Itoa(points) {
		index, _ := strconv.Atoi(string(number))
		text += numberEmojis[index]
	}
	return
}

func generate2048Buttons() (buttons []*disgord.MessageComponent) {
	buttons = []*disgord.MessageComponent{
		{
			Type: disgord.MessageComponentActionRow,
		},
	}
	for _, arrow := range arrows {
		buttons[0].Components = append(buttons[0].Components, &disgord.MessageComponent{
			Type:  disgord.MessageComponentButton,
			Style: disgord.Primary,
			Label: "\u200f",
			Emoji: &disgord.Emoji{
				Name: arrow,
			},
			CustomID: arrow,
		})
	}
	return
}

func run2048(itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	size := 4
	if len(itc.Data.Options) > 0 {
		size = int(itc.Data.Options[0].Value.(float64))
	}
	board := utils.MakeBoard(size, size)
	for i := 0; i < 2; i++ {
		board[utils.RandInt(size)][utils.RandInt(size)] = 2
	}
	lastPlay := time.Now()
	points := 0
	handler.Client.SendInteractionResponse(context.Background(), itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content:    ":zero:\n\n" + draw2048Board(board),
			Components: generate2048Buttons(),
		},
	})
	go func() {
		for {
			time.Sleep(time.Second)
			if time.Since(lastPlay).Seconds()/60 >= 2 {
				handler.DeleteHandler(itc.ID)

				handler.Client.EditInteractionResponse(context.Background(), itc, &disgord.Message{
					Content:    fmt.Sprintf(":skull:%s\n\n%s", drawPoints(points), draw2048Board(board)),
					Components: []*disgord.MessageComponent{},
				})
				return
			}
		}
	}()

	handler.RegisterHandler(itc.ID, func(interaction *disgord.InteractionCreate) {
		u := interaction.Member.UserID
		if u == itc.Member.UserID {
			oldBoard := utils.DeepClone(board)
			if interaction.Data.CustomID == arrows[0] {
				board = rotateBoard(board, true)
				board = rotateBoard(board, true)
				slideLeft(board, &points)
				board = rotateBoard(board, false)
				board = rotateBoard(board, false)
			} else if interaction.Data.CustomID == arrows[1] {
				board = rotateBoard(board, false)
				slideLeft(board, &points)
				board = rotateBoard(board, true)
			} else if interaction.Data.CustomID == arrows[2] {
				board = rotateBoard(board, true)
				slideLeft(board, &points)
				board = rotateBoard(board, false)
			} else if interaction.Data.CustomID == arrows[3] {
				slideLeft(board, &points)
			}
			if !utils.IsEqual(oldBoard, board) {
				lastPlay = time.Now()
				empty := []int{}
				for i, row := range board {
					for j, tile := range row {
						if tile == 0 {
							empty = append(empty, j+(i*len(board)))
						}
					}
				}
				if len(empty) > 0 {
					n := empty[utils.RandInt(len(empty))]
					board[n/len(board)][n%len(board)] = 2
				}
				handler.Client.SendInteractionResponse(context.Background(), interaction, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackUpdateMessage,
					Data: &disgord.CreateInteractionResponseData{
						Content: fmt.Sprintf("%s\n\n%s", drawPoints(points), draw2048Board(board)),
					},
				})
			}

		}

	}, 0)
	return nil
}
