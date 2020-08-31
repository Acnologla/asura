package commands

import (
	"asura/src/handler"
	"bytes"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	chessImage "github.com/cjsaylor/chessimage"
	"github.com/notnil/chess"
	"image/png"
	"strings"
	"sync"
	"time"
)

type chessGame struct {
	game     *chess.Game
	white    disgord.Snowflake
	black    disgord.Snowflake
	turn     int
	lastPlay time.Time
}

var (
	chessPlayers               = map[disgord.Snowflake]*chessGame{}
	chessMutex   *sync.RWMutex = &sync.RWMutex{}
)

func moveIncludes(moves []*chess.Move, str string) int {
	for i, move := range moves {
		if move.String() == str {
			return i
		}
	}
	return -1
}
func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"xadrez", "chess", "jogarxadrez"},
		Run:       runChess,
		Available: true,
		Cooldown:  2,
		Usage:     "j!chess <user>",
		Help:      "Jogue xadrez",
	})
}
func getGameImage(board *chessImage.Renderer) *bytes.Buffer {
	finalImage := bytes.NewBuffer([]byte{})
	image, _ := board.Render(chessImage.Options{
		AssetPath: "./resources/chess/",
	})
	png.Encode(finalImage, image)
	return finalImage
}
func runChess(session disgord.Session, msg *disgord.Message, args []string) {
	if len(args) > 0 {
		if strings.ToLower(args[0]) == "move" {
			if len(args) > 2 {
				args[1] = strings.ToLower(args[1])
				args[2] = strings.ToLower(args[2])
				currentGame := chessPlayers[msg.Author.ID]
				if currentGame == nil {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao esta jogando xadrez com ninguem, use j!chess <user> para começar a jogar")
					return
				}
				if (currentGame.black == msg.Author.ID && currentGame.turn == 0) || (currentGame.white == msg.Author.ID && currentGame.turn == 1) {
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Não é o seu turno")
					return
				}
				chessMutex.Lock()
				moves := currentGame.game.ValidMoves()
				validMove := moveIncludes(moves, args[1]+args[2])
				if validMove == -1 {
					chessMutex.Unlock()
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Movimento invalido")
					return
				}
				err := currentGame.game.Move(moves[validMove])
				if err != nil {
					chessMutex.Unlock()
					msg.Reply(context.Background(), session, msg.Author.Mention()+", Movimento invalido")
					return
				}
				var turnType string
				var nextPlay disgord.Snowflake
				if currentGame.turn == 0 {
					currentGame.turn = 1
					nextPlay = currentGame.black
					turnType = "Black"
				} else {
					currentGame.turn = 0
					nextPlay = currentGame.white
					turnType = "White"
				}
				currentGame.lastPlay = time.Now()
				chessMutex.Unlock()
				board, _ := chessImage.NewRendererFromFEN(currentGame.game.FEN())
				firstMove, _ := chessImage.TileFromAN(args[1])
				secondMove, _ := chessImage.TileFromAN(args[2])
				board.SetLastMove(chessImage.LastMove{
					From: firstMove,
					To:   secondMove,
				})
				user, err := session.GetUser(context.Background(), nextPlay)
				var username, avatar string
				if err == nil {
					username = user.Username
					avatar, _ = user.AvatarURL(512, false)
				}
				embed := &disgord.Embed{
					Title: "Xadrez",
					Color: 65535,
					Image: &disgord.EmbedImage{
						URL: "attachment://chess.png",
					},
				}
				currentMoves := currentGame.game.Moves()
				if len(currentMoves) > 0 {
					if currentMoves[len(currentMoves)-1].HasTag(chess.Check) {
						var squarePos string
						for key, value := range currentGame.game.Position().Board().SquareMap() {
							if value.Type().String() == "k" && value.Color().Name() == turnType {
								squarePos = strings.ToLower(key.String())
								break
							}
						}
						check, _ := chessImage.TileFromAN(squarePos)
						board.SetCheckTile(check)
					}
				}
				if currentGame.game.Method().String() == "Checkmate" {
					var currentPlay disgord.Snowflake
					if turnType == "white" {
						currentPlay = currentGame.black
					} else {
						currentPlay = currentGame.white
					}
					winner, err := session.GetUser(context.Background(), currentPlay)
					if err == nil {
						embed.Description = fmt.Sprintf("O %s venceu o xadrez", winner.Username)
					} else {
						embed.Description = fmt.Sprintf("O %s perdeu o xadrez", username)
					}
					delete(chessPlayers, currentPlay)
					delete(chessPlayers, nextPlay)
				} else {
					embed.Description = "Use j!chess <peça para mover> <objetivo>"
					embed.Footer = &disgord.EmbedFooter{
						IconURL: avatar,
						Text:    "Turno do: " + username,
					}
				}
				gameImage := getGameImage(board)
				msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
					Files: []disgord.CreateMessageFileParams{
						{gameImage, "chess.png", false},
					},
					Embed: embed,
				})
			} else {
				msg.Reply(context.Background(), session, msg.Author.Mention()+", Movimento invalido")
			}
			return
		}
	}
	if len(msg.Mentions) == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Mencione alguem para jogar xadrez")
		return
	}
	user := msg.Mentions[0]
	if chessPlayers[msg.Author.ID] != nil {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce ja esta em uma partida de xadrez")
		return
	} else if chessPlayers[user.ID] != nil {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Esse usuario ja esta em uma partida de xadrez")
		return
	}
	game := chess.NewGame()
	gameStruct := &chessGame{
		game:     game,
		turn:     0,
		white:    msg.Author.ID,
		black:    user.ID,
		lastPlay: time.Now(),
	}
	chessPlayers[msg.Author.ID] = gameStruct
	chessPlayers[user.ID] = gameStruct
	board, _ := chessImage.NewRendererFromFEN(game.FEN())
	gameImage := getGameImage(board)
	avatar, _ := msg.Author.AvatarURL(512, false)
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{
			{gameImage, "chess.png", false},
		},
		Embed: &disgord.Embed{
			Color:       65535,
			Title:       "Xadrez",
			Description: "Use j!chess move <peça para mover> <objetivo>",
			Image: &disgord.EmbedImage{
				URL: "attachment://chess.png",
			},
			Footer: &disgord.EmbedFooter{
				IconURL: avatar,
				Text:    "Turno do: " + msg.Author.Username,
			},
		},
	})
	for {
		time.Sleep(1 * time.Minute)
		if chessPlayers[msg.Author.ID] == nil {
			break
		}
		if time.Since(chessPlayers[msg.Author.ID].lastPlay).Seconds()/60 >= 3 {
			var text string
			turn := chessPlayers[msg.Author.ID].turn
			if turn == 0 {
				text = fmt.Sprintf("O %s ganhou o xadrez", user.Username)
			} else {
				text = fmt.Sprintf("O %s ganhou o xadrez", msg.Author.Username)
			}
			board, _ := chessImage.NewRendererFromFEN(chessPlayers[msg.Author.ID].game.FEN())
			gameImage := getGameImage(board)
			msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
				Files: []disgord.CreateMessageFileParams{
					{gameImage, "chess.png", false},
				},
				Embed: &disgord.Embed{
					Color:       65535,
					Description: text,
					Title:       "Xadrez",
					Image: &disgord.EmbedImage{
						URL: "attachment://chess.png",
					},
				},
			})
			delete(chessPlayers, msg.Author.ID)
			delete(chessPlayers, user.ID)
		}
	}
}
