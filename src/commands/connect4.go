package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"github.com/andersfylling/disgord"
	"strings"
)

var connect4Emojis = map[int]string{
	0: "⚪",
	1: "🟡",
	2: "🔴",
}

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"connect4", "c4"},
		Run:       runConnect4,
		Available: true,
		Cooldown:  10,
		Usage:     "j!connect4 @user",
		Help:      "Jogue connect4",
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
			actual = board[j][i]
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
func runConnect4(session disgord.Session, msg *disgord.Message, args []string) {
	ctx := context.Background()
	board := utils.MakeBoard(8)
	if len(msg.Mentions) == 0 {
		msg.Reply(ctx, session, msg.Author.Mention()+", Voce precisa mencionar alguem para jogar connect4")
		return
	}
	user := msg.Mentions[0]
	if user.Bot || user.ID == msg.Author.ID {
		msg.Reply(ctx, session, msg.Author.Mention()+", Usuario invalido")
		return
	}
	message, err := msg.Reply(ctx, session, connect4Emojis[1]+"\n\n"+drawConnect4Board(board))
	if err == nil {
		for _, emoji := range emojis[:len(board)] {
			utils.Try(func() error {
				return message.React(ctx, session, emoji)
			}, 5)
		}
		turn := 1
		mes :=  session.Channel(message.ChannelID).Message(message.ID)
		handler.RegisterHandler(message, func(removed bool, emoji disgord.Emoji, u disgord.Snowflake) {
			turnUser := user
			num := 2
			if turn == 1 {
				turnUser = msg.Author
				num = 1
			}
			if !removed && u == turnUser.ID {
				if utils.Includes(emojis[:len(board)], emoji.Name) {
					tile := utils.IndexOf(emojis, emoji.Name)
					played := false
					for j := len(board) - 1; 0 <= j; j-- {
						if board[j][tile] == 0 {
							board[j][tile] = num
							played = true
							break
						}
					}
					if played {
						msgUpdater := mes.Update()
						winned := checkConnect4Win(board)
						if winned != 0 {
							emoji := "❌"
							if winned == 1 {
								emoji = ":crown:" + connect4Emojis[num]
							}
							msgUpdater.SetContent(emoji + "\n\n" + drawConnect4Board(board))
							msgUpdater.Execute()
							handler.DeleteHandler(message)
							mes.DeleteAllReactions()
						} else {
							msgUpdater.SetContent(connect4Emojis[turn+1] + "\n\n" + drawConnect4Board(board))
							if turn == 1 {
								turn = 0
							} else {
								turn = 1
							}
							msgUpdater.Execute()
							mes.Reaction(emoji.Name).DeleteUser(u)
						}
					}
				}
			} else if u != message.Author.ID  {
				mes.Reaction(emoji.Name).DeleteUser(u)
			}
		}, 60*20)
	}
}
