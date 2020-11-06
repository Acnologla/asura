package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"jogodavelha", "ttt", "tictactoe"},
		Run:       runTTT,
		Available: true,
		Cooldown:  8,
		Usage:     "j!jogodavelha @usuario",
		Help:      "Jogue jogo da velha",
	})
}

var emojis = []string{"1⃣", "2⃣", "3⃣", "4⃣", "5⃣", "6⃣", "7⃣", "8⃣", "9⃣"}

func text(tiles [9]int) (text string) {
	for i, tile := range tiles {
		if tile == 1 {
			text += ":x:"
		} else if tile == 2 {
			text += ":o:"
		} else {
			text += emojis[i]
		}
		if i == 2 || i == 5 {
			text += "\n"
		}
	}
	return
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

func runTTT(session disgord.Session, msg *disgord.Message, args []string) {
	if len(msg.Mentions) == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Mencione alguem para jogar jogo da velha")
		return
	}
	user := msg.Mentions[0]
	if user.ID == msg.Author.ID || user.Bot {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Usuario invalido")
		return
	}
	ctx := context.Background()
	tiles := [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	turn := 1
	message, err := msg.Reply(ctx, session, ":hash::x:\n\n"+text(tiles))
	if err != nil {
		return
	}
	for i := 0; i < len(emojis); i++ {
		utils.Try(func() error {
			return message.React(ctx, session, emojis[i])
		}, 5)
	}
	mes := session.Channel(message.ChannelID).Message(message.ID)
	handler.RegisterHandler(message, func(removed bool, emoji disgord.Emoji, u disgord.Snowflake) {
		turnUser := user
		letter := ":x:"
		if turn == 1 {
			turnUser = msg.Author
			letter = ":o:"
		}
		if !removed && turnUser.ID == u {
			if utils.Includes(emojis, emoji.Name) {
				tile := utils.IndexOf(emojis, emoji.Name)
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
						handler.DeleteHandler(message)
						go mes.DeleteAllReactions()
					}
					msgUpdater := mes.Update()
					msgUpdater.SetContent(fmt.Sprintf("%s%s\n\n%s", playType, letter, text(tiles)))
					msgUpdater.Execute()
					message.Unreact(ctx, session, emoji.Name)
					mes.Reaction(emoji.Name).DeleteUser(u)
				}
			}
		} else if u != message.Author.ID {
			mes.Reaction(emoji.Name).DeleteUser(u)
		}
	}, 60*5)
}
