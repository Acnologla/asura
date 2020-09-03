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

func rotateBoard(board []([]int), c bool) []([]int) {
	var rotatedBoard = makeBoard(len(board))
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

func isEqual(oldBoard []([]int), board []([]int)) bool {
	equal := true
	for i, row := range oldBoard {
		for j, tile := range row {
			if tile != board[i][j] {
				equal = false
				break
			}
		}
	}
	return equal
}

func drawBoard(board []([]int)) (text string) {
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

func makeBoard(size int) []([]int) {
	board := []([]int){}
	for i := 0; i < size; i++ {
		row := []int{}
		for j := 0; j < size; j++ {
			row = append(row, 0)
		}
		board = append(board, row)
	}
	return board
}
func deepClone(board []([]int)) []([]int){
	arr := []([]int){}
	for _, row := range board{
		newRow := []int{}
		for _,tile := range row{
			newRow = append(newRow,tile)
		}
		arr = append(arr,newRow)
	}
	return arr
}

func run2048(session disgord.Session, msg *disgord.Message, args []string) {
	size := 4
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err == nil {
			if 3 > n || n > 6 {
				msg.Reply(context.Background(), session, "O Tamanho do tabuleiro deve ser entre 3 e 6")
				return
			}
			size = n
		}
	}
	board := makeBoard(size)
	for i := 0; i < 2; i++ {
		board[rand.Intn(size)][rand.Intn(size)] = 2
	}
	lastPlay := time.Now()
	points := 0
	message, err := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Content: ":zero:\n\n" + drawBoard(board),
	})
	if err == nil {
		for i := 0; i < 4; i++ {
			utils.Try(func() error {
				return message.React(context.Background(), session, arrows[i])
			}, 5)
		}
		go func() {
			for {
				time.Sleep(time.Second)
				if time.Since(lastPlay).Seconds()/60 >= 2 {
					handler.DeleteHandler(message)
					handler.Client.DeleteAllReactions(context.Background(), msg.ChannelID, message.ID)
					msgUpdater := handler.Client.UpdateMessage(context.Background(), msg.ChannelID, message.ID)
					msgUpdater.SetContent(fmt.Sprintf(":skull:%s\n\n%s", drawPoints(points), drawBoard(board)))
					utils.Try(func() error {
						_, err := msgUpdater.Execute()
						return err
					}, 10)
					return
				}
			}
		}()
		handler.RegisterHandler(message, func(removed bool, emoji disgord.Emoji, u disgord.Snowflake) {
			if !removed {
				if u == msg.Author.ID {
					if utils.Includes(arrows, emoji.Name) {
						var oldBoard = deepClone(board)
						if emoji.Name == arrows[0] {
							board = rotateBoard(board, true)
							board = rotateBoard(board, true)
							slideLeft(board, &points)
							board = rotateBoard(board, false)
							board = rotateBoard(board, false)
						} else if emoji.Name == arrows[1] {
							board = rotateBoard(board, false)
							slideLeft(board, &points)
							board = rotateBoard(board, true)
						} else if emoji.Name == arrows[2] {
							board = rotateBoard(board, true)
							slideLeft(board, &points)
							board = rotateBoard(board, false)
						} else if emoji.Name == arrows[3] {
							slideLeft(board, &points)
						}					
						if !isEqual(oldBoard, board) {
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
								n := empty[rand.Intn(len(empty))]
								board[n/len(board)][n%len(board)] = 2
							}
							msgUpdater := handler.Client.UpdateMessage(context.Background(), msg.ChannelID, message.ID)
							msgUpdater.SetContent(fmt.Sprintf("%s\n\n%s", drawPoints(points), drawBoard(board)))
							utils.Try(func() error {
								_, err := msgUpdater.Execute()
								return err
							}, 10)
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