package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"math/rand"
	"strconv"
	"time"
)

var arrows = []string{"➡", "⬇", "⬆", "⬅"}
var numberEmojis = [10]string{":zero:", ":one:", ":two:", ":three:", ":four:", ":five:", ":six:", ":seven:", ":eight:", ":nine"}
var emojis2048 = map[int]string{
	0:    "<:2048_0:744682229023768737>",
	2:    "<:284_2:744682213047795794>",
	4:    "<:248_4:744682221452918905>",
	8:    "<:248_8:744682233855475713>",
	16:   "<:248_16:744682222145110147>",
	32:   "<:248_32:744682243276013650>",
	64:   "<:248_64:744682218538008686>",
	128:  "<:248_128:744682220559794287>",
	256:  "<:248_256:744682236544155699>",
	512:  "<:258_512:744682226482020482>",
	1024: "<:248_1024:744682235164229653>",
	2048: "<:248:744682245113249853>",
}

func slideLeft(board *[4][4]int, points *int) {
	for i, _ := range board {
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

func rotateBoard(board *[4][4]int, c bool) [4][4]int {
	var rotatedBoard = [4][4]int{}
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

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"2048", "jogar2048"},
		Run:       run2048,
		Available: true,
		Cooldown:  20,
		Usage:     "j!2048",
		Help:      "Jogue 2048",
	})
}

func drawBoard(board [4][4]int) (text string) {
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

func run2048(session disgord.Session, msg *disgord.Message, args []string) {
	board := [4][4]int{}
	for i := 0; i < 2; i++ {
		board[rand.Intn(4)][rand.Intn(4)] = 2
	}
	lastPlay := time.Now()
	points := 0
	message, err := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Title:       "2048",
			Description: ":zero:\n\n" + drawBoard(board),
		},
		Content: msg.Author.Mention(),
	})
	if err == nil {
		for i := 0; i < 4; i++ {
			err := message.React(context.Background(), session, arrows[i])
			if err != nil {
				time.Sleep(200 * time.Millisecond)
				i--
			}
		}
		go func() {
			for {
				time.Sleep(time.Second)

				if time.Since(lastPlay).Seconds()/60 >= 2 {
					handler.DeleteHandler(message)
					handler.Client.DeleteAllReactions(context.Background(), msg.ChannelID, message.ID)
					msgUpdater := handler.Client.UpdateMessage(context.Background(), msg.ChannelID, message.ID)
					msgUpdater.SetEmbed(&disgord.Embed{
						Title:       "2048",
						Color:       16711680,
						Description: fmt.Sprintf(":skull:%s\n\n%s", drawPoints(points), drawBoard(board)),
						Footer: &disgord.EmbedFooter{
							Text: "Voce não jogou durante 2 minutos",
						},
					})
					for {
						_, err := msgUpdater.Execute()
						if err == nil {
							break
						} else {
							time.Sleep(200 * time.Millisecond)
						}
					}
				}
			}
		}()
		handler.RegisterHandler(message, func(removed bool, emoji disgord.Emoji, u disgord.Snowflake) {
			if !removed {
				if u == msg.Author.ID {
					if utils.Includes(arrows, emoji.Name) {
						var oldPoints = points
						if emoji.Name == arrows[0] {
							board = rotateBoard(&board, true)
							board = rotateBoard(&board, true)
							slideLeft(&board, &points)
							board = rotateBoard(&board, false)
							board = rotateBoard(&board, false)
						} else if emoji.Name == arrows[1] {
							board = rotateBoard(&board, false)
							slideLeft(&board, &points)
							board = rotateBoard(&board, true)
						} else if emoji.Name == arrows[2] {
							board = rotateBoard(&board, true)
							slideLeft(&board, &points)
							board = rotateBoard(&board, false)
						} else if emoji.Name == arrows[3] {
							slideLeft(&board, &points)
						}
						if oldPoints != points {
							lastPlay = time.Now()
						}
						empty := []int{}
						for i, row := range board {
							for j, tile := range row {
								if tile == 0 {
									empty = append(empty, j+(i*len(board)))
								}
							}
						}
						n := empty[rand.Intn(len(empty))]
						board[n/len(board)][n%len(board)] = 2
						msgUpdater := handler.Client.UpdateMessage(context.Background(), msg.ChannelID, message.ID)
						msgUpdater.SetEmbed(&disgord.Embed{
							Title:       "2048",
							Description: fmt.Sprintf("%s\n\n%s", drawPoints(points), drawBoard(board)),
						})
						for {
							_, err := msgUpdater.Execute()
							if err == nil {
								break
							} else {
								time.Sleep(200 * time.Millisecond)
							}
						}
						handler.Client.DeleteUserReaction(context.Background(), msg.ChannelID, message.ID, u, emoji.Name)
					}
				} else if u != message.Author.ID {
					handler.Client.DeleteUserReaction(context.Background(), msg.ChannelID, message.ID, u, emoji.Name)
				}
			}
		}, 0)
	}
}
