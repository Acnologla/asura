package commands

import (
	"asura/src/handler"
	"asura/src/translation"
	"asura/src/utils"
	"context"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

var emojis = []string{"1⃣", "2⃣", "3⃣", "4⃣", "5⃣", "6⃣", "7⃣", "8⃣", "9⃣"}

var connect4Emojis = map[int]string{
	0: "⚪",
	1: "🟡",
	2: "🔴",
}

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "connect4",
		Description: translation.T("Connect4Help", "pt"),
		Run:         runConnect4,
		Cooldown:    20,
		Category:    handler.Games,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Type:        disgord.OptionTypeUser,
			Required:    true,
			Name:        "user",
			Description: translation.T("Connect4Desc", "pt"),
		}),
	})
}

func drawConnect4Board(board []([]int)) (text string) {
	text = strings.Join(emojis[:len(board)], "") + "\n"
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			text += connect4Emojis[board[i][j]]
		}
		text += "\n"
	}
	return
}

func checkConnect4Win(board []([]int)) int {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board)-3; j++ {
			actual := board[i][j]
			if actual == 0 {
				continue
			}
			if board[i][j+1] == actual && board[i][j+2] == actual && board[i][j+3] == actual {
				return 1
			}
		}
		for j := 0; j < len(board)-3; j++ {
			actual := board[j][i]
			if actual == 0 {
				continue
			}
			if board[j+1][i] == actual && board[j+2][i] == actual && board[j+3][i] == actual {
				return 1
			}
		}
	}
	for i := 0; i < len(board)-3; i++ {
		for j := 0; j < len(board)-3; j++ {
			actual := board[i][j]
			if actual == 0 {
				continue
			}
			if board[i+1][j+1] == actual && board[i+2][j+2] == actual && board[i+3][j+3] == actual {
				return 1
			}
		}
		for j := 3; j < len(board); j++ {
			actual := board[i][j]
			if actual == 0 {
				continue
			}
			if board[i+1][j-1] == actual && board[i+2][j-2] == actual && board[i+3][j-3] == actual {
				return 1
			}
		}
	}
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			if board[i][j] == 0 {
				return 0
			}
		}
	}
	return 2
}

func generateConnect4() []*disgord.MessageComponent {
	arrs := []*disgord.MessageComponent{
		{
			Type: disgord.MessageComponentActionRow,
		},
		{
			Type: disgord.MessageComponentActionRow,
		},
	}
	for i := 0; i < 8; i++ {
		var arrField int = i / 4
		button := &disgord.MessageComponent{
			Type:     disgord.MessageComponentButton,
			Style:    disgord.Primary,
			CustomID: strconv.Itoa(i + 1),
			Label:    strconv.Itoa(i + 1),
		}
		arrs[arrField].Components = append(arrs[arrField].Components, button)
	}
	return arrs
}

func runConnect4(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	board := utils.MakeBoard(8, 8)
	user := utils.GetUser(itc, 0)
	if user.Bot || user.ID == itc.Member.UserID {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Invalid user",
			},
		}
	}

	itcID, err := handler.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content:    connect4Emojis[1] + "\n\n" + drawConnect4Board(board),
			Components: generateConnect4(),
		},
	})
	if err != nil {
		return nil
	}
	turn := 1
	handler.RegisterHandler(itcID, func(interaction *disgord.InteractionCreate) {
		turnUser := user
		num := 2
		u := interaction.Member.UserID
		emoji := interaction.Data.CustomID
		if turn == 1 {
			turnUser = itc.Member.User
			num = 1
		}
		if u == turnUser.ID {
			tile, _ := strconv.Atoi(emoji)
			tile--
			played := false
			for j := len(board) - 1; 0 <= j; j-- {
				if board[j][tile] == 0 {
					board[j][tile] = num
					played = true
					break
				}
			}
			if played {
				winned := checkConnect4Win(board)
				if winned != 0 {
					emoji := "❌"
					if winned == 1 {
						emoji = ":crown:" + connect4Emojis[num]
					}
					handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
						Type: disgord.InteractionCallbackUpdateMessage,
						Data: &disgord.CreateInteractionResponseData{
							Content:    emoji + "\n\n" + drawConnect4Board(board),
							Components: generateConnect4(),
						},
					})
				} else {
					handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
						Type: disgord.InteractionCallbackUpdateMessage,
						Data: &disgord.CreateInteractionResponseData{
							Content:    connect4Emojis[turn+1] + "\n\n" + drawConnect4Board(board),
							Components: generateConnect4(),
						},
					})
					if turn == 1 {
						turn = 0
					} else {
						turn = 1
					}

				}
			}

		}
	}, 60*20)
	return nil
}
