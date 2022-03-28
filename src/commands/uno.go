package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"bytes"
	"context"
	"fmt"
	"image/png"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

var cards []string = []string{}

func shuffle(arr []string) []string {
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}
func findCard(cards []string, card string) string {
	for _, playerCard := range cards {
		if len(playerCard) > 5 {
			if playerCard[:5] == card {
				return playerCard
			}
		}
		if playerCard == card {
			return playerCard
		}
	}
	return ""
}

type player struct {
	id       disgord.Snowflake
	username string
	cards    []string
	uno      bool
}

type unoGame struct {
	cards    []string
	players  []*player
	board    []string
	refuse   int
	status   string
	lastPlay time.Time
}

func (game *unoGame) last() string {
	return game.board[len(game.board)-1]
}

func (game *unoGame) getBoard() *bytes.Buffer {
	file, _ := os.Open(fmt.Sprintf("resources/uno/%s.png", game.last()))
	defer file.Close()
	// Decodes it into a image.Image
	img, _ := png.Decode(file)
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, img)
	return &b
}

func (game *unoGame) buy(p *player, count int) {
	for i := 0; i < count; i++ {
		if len(game.cards) == 0 {
			gameCards := append([]string{}, game.board[:len(game.board)-2]...)
			for i, card := range gameCards {
				if strings.HasSuffix(card, "WILDP4") {
					gameCards[i] = "WILDP4"
				} else if strings.HasSuffix(card, "WILD") {
					gameCards[i] = "WILD"
				}
			}
			if 1 < len(gameCards) {
				game.cards = shuffle(gameCards)
				game.board = game.board[len(game.board)-1:]
			}
		}
		if len(game.cards) != 0 {
			newCard := game.cards[:1][0]
			game.cards = game.cards[1:]
			p.cards = append(p.cards, newCard)
		}
	}
}

func (game *unoGame) play(p *player, cards []string) string {
	card := findCard(p.cards, cards[0])
	if card == "" {
		return "Voce não tem essa carta copia a  da id de cima da carta"
	}
	last := game.last()
	if last[0] == card[0] || card[0] == 'W' || card[1:] == last[1:] {
		var action string
		if card[0] == 'W' {
			action = card[4:]
		} else {
			action = card[1:]
		}
		cardIndex := utils.IndexOf(p.cards, card)
		p.cards = utils.Splice(p.cards, cardIndex)
		if action != "REVERSE" {
			game.addToLast()
		}
		if action == "REVERSE" {
			for i, j := 0, len(game.players)-1; i < j; i, j = i+1, j-1 {
				game.players[i], game.players[j] = game.players[j], game.players[i]
			}
			game.board = append(game.board, card)
		} else if action == "SKIP" {
			skipped := game.players[0]
			game.players = game.players[1:]
			game.players = append(game.players, skipped)
			game.board = append(game.board, card)
		} else if action == "P2" {
			game.buy(game.players[0], 2)
			game.board = append(game.board, card)
		} else if action == "P4" {
			game.buy(game.players[0], 4)
			if cards[1] != "Y" && cards[1] != "B" && cards[1] != "G" && cards[1] != "R" {
				cards[1] = "Y"
			}
			game.board = append(game.board, cards[1]+card)
		} else if action == "" {
			if cards[1] != "Y" && cards[1] != "B" && cards[1] != "G" && cards[1] != "R" {
				cards[1] = "Y"
			}
			game.board = append(game.board, cards[1]+card)
		} else {
			game.board = append(game.board, card)
		}
		return ""
	} else {
		return "Jogada invalida"
	}
}
func (game *unoGame) addToLast() {
	p := game.players[0]
	game.players = game.players[1:]
	game.players = append(game.players, p)
}
func (game *unoGame) drawCards(p *player) *bytes.Buffer {
	dc := gg.NewContext(7*80+7*35, (int(len(p.cards)/7)-1)*110+330)
	dc.SetHexColor("#FFFFFF")
	dc.LoadFontFace("./resources/Raleway-Black.ttf", 20)
	i, h := 0, 30
	for _, card := range p.cards {
		file, _ := os.Open(fmt.Sprintf("resources/uno/%s.png", card))
		defer file.Close()
		img, _ := png.Decode(file)
		if i%7 == 0 {
			i = 0
			h += 110
		}
		img = resize.Resize(80, 80, img, resize.Lanczos3)
		dc.DrawImage(img, i*80+35, h)
		if len(card) > 5 {
			dc.DrawString(card[:5], float64(i*80+35+25), float64(h-5))
		} else {
			dc.DrawString(card, float64(i*80+35+25), float64(h-5))
		}
		i++
	}
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())
	return &b
}

func (game *unoGame) removePlayer(p *disgord.User) {
	var position int
	for i, p2 := range game.players {
		if p2.id == p.ID {
			position = i
			for _, card := range p2.cards {
				game.cards = append(game.cards, card)
			}
			break
		}
	}
	game.players = append(game.players[:position], game.players[position+1:]...)
}

func (game *unoGame) newPlayer(p *disgord.User) {
	newCards := append([]string{}, game.cards[:7]...)
	game.cards = game.cards[7:]
	newPlayer := player{
		id:       p.ID,
		username: p.Username,
		cards:    newCards,
	}
	game.players = append(game.players, &newPlayer)

}

func (game *unoGame) start() {
	var initialCard string
	for _, card := range game.cards {
		if len(card) == 2 {
			initialCard = card
			break
		}
	}
	game.board = append(game.board, initialCard)
	game.status = "n"
}

var (
	currentGames               = map[disgord.Snowflake]*unoGame{}
	gameMutex    *sync.RWMutex = &sync.RWMutex{}
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "uno",
		Description: translation.T("UnoHelp", "pt"),
		Run:         runUno,
		Cooldown:    3,
		Category:    handler.Games,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "create",
				Description: translation.T("UnoCreateHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "join",
				Description: translation.T("UnoJoinHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type: disgord.OptionTypeSubCommand,
				Name: "play",
				Options: utils.GenerateOptions(
					&disgord.ApplicationCommandOption{
						Type:        disgord.OptionTypeString,
						Name:        "id",
						Required:    true,
						Description: translation.T("UnoPlayIDHelp", "pt"),
					},
				),
				Description: translation.T("UnoPlayHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "buy",
				Description: translation.T("UnoBuyHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "uno",
				Description: translation.T("UnoUnoHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "leave",
				Description: translation.T("UnoLeaveHelp", "pt"),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "cards",
				Description: translation.T("UnoCardsHelp", "pt"),
			},
		),
	})
	colors := []string{"B", "R", "Y", "G"}
	for i := 0; i <= 9; i++ {
		cards = append(cards, fmt.Sprintf("B%d", i), fmt.Sprintf("B%d", i))
		cards = append(cards, fmt.Sprintf("R%d", i), fmt.Sprintf("R%d", i))
		cards = append(cards, fmt.Sprintf("Y%d", i), fmt.Sprintf("Y%d", i))
		cards = append(cards, fmt.Sprintf("G%d", i), fmt.Sprintf("G%d", i))
	}
	for i := 0; i < 4; i++ {
		cards = append(cards, "WILDP4")
		for o := 0; o < 2; o++ {
			cards = append(cards, colors[i]+"P2", colors[i]+"REVERSE", colors[i]+"SKIP", "WILD")
		}
	}
}

func getCards(ctx context.Context, itc *disgord.InteractionCreate, session disgord.Session) {
	gameMutex.RLock()
	game := currentGames[itc.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo nessa guilda, use /uno create para criar um",
			},
		})
		return
	}
	if game.status == "n" {
		var isInGame = isInGame(itc, game)
		if isInGame {
			var p *player
			for _, player2 := range game.players {
				if player2.id == itc.Member.UserID {
					p = player2
					break
				}
			}
			gameMutex.RUnlock()
			b := bytes.NewReader(game.drawCards(p).Bytes())
			chanID, err := handler.Client.User(itc.Member.UserID).CreateDM()
			if err == nil {
				go handler.Client.Channel(chanID.ID).CreateMessage(&disgord.CreateMessage{
					Content: "Para jogar uma carta usa /uno play <id da carta> ( o id é o coiso que ta encima da carta)",
					Files: []disgord.CreateMessageFile{
						{b, "uno.png", false},
					},
				})
			}
		} else {
			gameMutex.RUnlock()
			go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Você não está no jogo",
				},
			})
		}
	} else {
		gameMutex.RUnlock()
	}
}

func sendEmbed(ctx context.Context, itc *disgord.InteractionCreate, session disgord.Session, isItc bool) {
	gameMutex.RLock()
	defer gameMutex.RUnlock()
	game := currentGames[itc.GuildID]
	ch := session.Channel(itc.ChannelID)
	if !isItc {
		go ch.CreateMessage(&disgord.CreateMessage{
			Files: []disgord.CreateMessageFileParams{
				{bytes.NewReader(game.getBoard().Bytes()), "uno.png", false},
			},
			Embeds: []*disgord.Embed{
				{
					Title: "Uno",
					Footer: &disgord.EmbedFooter{
						Text: fmt.Sprintf("Turno do: %s", game.players[0].username),
					},
					Description: "Use /uno play <carta> para jogar!\n/uno leave (para sair)\nVeja suas cartas no privado (usando /cartas)\n/uno buy (para comprar uma carta e passar o turno\n caso não jogue ou não compre durante 3 minutos\n ira passar o turno e comprar automatico)\nCaso estiver no uno (uma carta) use /uno uno",
					Image: &disgord.EmbedImage{
						URL: "attachment://uno.png",
					},
				},
			},
		})
	} else {
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Files: []disgord.CreateMessageFileParams{
					{bytes.NewReader(game.getBoard().Bytes()), "uno.png", false},
				},
				Embeds: []*disgord.Embed{
					{
						Title: "Uno",
						Footer: &disgord.EmbedFooter{
							Text: fmt.Sprintf("Turno do: %s", game.players[0].username),
						},
						Description: "Use /uno play <carta> para jogar!\n/uno leave (para sair)\nVeja suas cartas no privado (usando /cartas)\n/uno buy (para comprar uma carta e passar o turno\n caso não jogue ou não compre durante 3 minutos\n ira passar o turno e comprar automatico)\nCaso estiver no uno (uma carta) use /uno uno",
						Image: &disgord.EmbedImage{
							URL: "attachment://uno.png",
						},
					},
				},
			},
		})
	}

}

func create(ctx context.Context, itc *disgord.InteractionCreate, session disgord.Session) {
	gameMutex.RLock()
	game := currentGames[itc.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		gameMutex.Lock()
		gameCards := append([]string{}, shuffle(cards)...)
		currentGames[itc.GuildID] = &unoGame{
			cards:    gameCards,
			board:    []string{},
			status:   "p",
			lastPlay: time.Now(),
		}
		currentGames[itc.GuildID].newPlayer(itc.Member.User)
		gameMutex.Unlock()
		itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title: "Uno",
						Thumbnail: &disgord.EmbedThumbnail{
							URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/UNO_Logo.svg/1024px-UNO_Logo.svg.png",
						},
						Description: fmt.Sprintf("Use **/uno join** para entrar ( O jogo começara em **1** minuto)\nParticipantes:\n%s", itc.Member.User.Username),
						Color:       16711680,
					},
				},
			},
		})
		channel := session.Channel(itc.ChannelID)
		go func() {
			time.Sleep(60 * time.Second)
			gameMutex.RLock()
			game = currentGames[itc.GuildID]
			if 1 >= len(game.players) {
				gameMutex.RUnlock()
				go channel.CreateMessage(&disgord.CreateMessageParams{
					Content: "Como só teve um jogador, o uno foi cancelado, uno /uno create para jogar novamente",
				})
				gameMutex.Lock()
				delete(currentGames, itc.GuildID)
				gameMutex.Unlock()
			} else {
				go func() {
					for {
						time.Sleep(60 * time.Second)
						gameMutex.RLock()
						game := currentGames[itc.GuildID]
						if game != nil {
							if time.Since(game.lastPlay).Seconds()/60 >= 2 {
								if game.refuse == 4 {
									gameMutex.RUnlock()
									gameMutex.Lock()
									delete(currentGames, itc.GuildID)
									gameMutex.Unlock()
									channel.CreateMessage(&disgord.CreateMessage{
										Content: "4 rodadas sem jogar consecutivas, jogo acabado",
									})
									return
								} else {
									go channel.CreateMessage(&disgord.CreateMessage{
										Content: fmt.Sprintf("O %s pulou seu turno e comprou uma carta agora é o turno de %s", game.players[0].username, game.players[1].username),
									})
									gameMutex.RUnlock()
									gameMutex.Lock()
									game = currentGames[itc.GuildID]
									game.buy(game.players[0], 1)
									game.addToLast()
									game.lastPlay = time.Now()
									game.refuse++
									gameMutex.Unlock()
									return
								}
							}
						} else {
							gameMutex.RUnlock()
							return
						}
						gameMutex.RUnlock()
					}
				}()
				for _, player2 := range game.players {
					go func(p *player) {
						chanID, err := handler.Client.User(p.id).CreateDM()
						if err == nil {
							handler.Client.Channel(chanID.ID).CreateMessage(&disgord.CreateMessage{
								Content: "Para jogar uma carta usa /uno play <id da carta> ( o id é o coiso que ta encima da carta)",
								Files: []disgord.CreateMessageFile{
									{
										Reader:   bytes.NewReader(game.drawCards(p).Bytes()),
										FileName: "uno.png",
									},
								},
							})
						}
					}(player2)
				}
				gameMutex.RUnlock()
				gameMutex.Lock()
				game.start()
				gameMutex.Unlock()
				sendEmbed(ctx, itc, session, false)
			}
		}()
	} else if game.status == "p" {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "O jogo em sua guilda esta em preparação, caso queira entrar use /uno join",
			},
		})
	} else {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "O jogo da sua guilda ja começou, espere ele terminar",
			},
		})
	}
}
func isInGame(itc *disgord.InteractionCreate, game *unoGame) bool {
	var isInGame bool
	for _, player := range game.players {
		if player.id == itc.Member.UserID {
			isInGame = true
			break
		}
	}
	return isInGame
}
func leave(ctx context.Context, itc *disgord.InteractionCreate, session disgord.Session) {
	gameMutex.RLock()
	game := currentGames[itc.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo na sua guilda",
			},
		})
		return
	}
	if game.status != "n" {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo na sua guilda",
			},
		})
		return
	}
	var isInGame = isInGame(itc, game)
	if !isInGame {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: " Voce não ta no jogo",
			},
		})
		return
	}
	gameMutex.RUnlock()
	gameMutex.Lock()
	game.removePlayer(itc.Member.User)
	game = currentGames[itc.GuildID]
	gameMutex.Unlock()
	if 1 >= len(game.players) {
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Como todos os jogadores sairam ou só sobrou 1 o jogo foi encerrado",
			},
		})
		gameMutex.Lock()
		delete(currentGames, itc.GuildID)
		gameMutex.Unlock()
		return
	}
	itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: "Voce saiu do jogo agora é o turno do " + game.players[0].username,
		},
	})
}

func buy(ctx context.Context, itc *disgord.InteractionCreate, session disgord.Session) {
	gameMutex.RLock()
	game := currentGames[itc.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo na sua guilda",
			},
		})
		return
	}
	if game.status != "n" {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo na sua guilda",
			},
		})
		return
	}
	var isInGame = isInGame(itc, game)
	if !isInGame {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Voce nao ta no jogo",
			},
		})
		return
	}
	gameMutex.RUnlock()
	gameMutex.Lock()
	var p *player
	for _, player2 := range game.players {
		if player2.id == itc.Member.UserID {
			p = player2
			break
		}
	}
	if p.id != game.players[0].id {
		gameMutex.Unlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não é seu turno",
			},
		})
		return
	}
	game.buy(p, 1)
	game.addToLast()
	game.lastPlay = time.Now()
	gameMutex.Unlock()
	getCards(ctx, itc, session)
	sendEmbed(ctx, itc, session, true)
}

func play(ctx context.Context, itc *disgord.InteractionCreate, session disgord.Session, args []string) {
	gameMutex.RLock()
	game := currentGames[itc.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo na sua guilda",
			},
		})
		return
	}
	if game.status != "n" {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo na sua guilda",
			},
		})
		return
	}
	if len(args) == 0 {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Diga uma carta para jogar",
			},
		})
		return
	}
	var isInGame = isInGame(itc, game)
	if !isInGame {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Voce nao ta no jogo",
			},
		})
		return
	}
	if strings.HasPrefix(args[0], "WILD") && 2 > len(args) {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Essa é uma carta preta usa  então voce tem que definir a cor apos a carta exemplo:\n/uno play wild <Y (Amarelo) | G (Verde) | B (Azul) | R (Vermelho) >",
			},
		})
		return
	}
	gameMutex.RUnlock()
	gameMutex.Lock()
	game = currentGames[itc.GuildID]
	var p *player
	for _, player2 := range game.players {
		if player2.id == itc.Member.UserID {
			p = player2
			break
		}
	}
	if p.id != game.players[0].id {
		gameMutex.Unlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não é seu turno",
			},
		})
		return
	}
	play := game.play(p, args)
	if play != "" {
		gameMutex.Unlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: play,
			},
		})
		return
	}
	game = currentGames[itc.GuildID]
	game.refuse = 0
	game.lastPlay = time.Now()
	gameMutex.Unlock()
	sendEmbed(ctx, itc, session, true)
	ch := session.Channel(itc.ChannelID)
	if len(p.cards) == 0 {
		gameMutex.Lock()
		if !p.uno {
			game.buy(p, 2)
			gameMutex.Unlock()
			go ch.CreateMessage(&disgord.CreateMessage{
				Content: "Voce comprou 2 cartas por não ter falado uno",
			})
			getCards(ctx, itc, session)
			return
		}
		go ch.CreateMessage(&disgord.CreateMessage{
			Embeds: []*disgord.Embed{{
				Title:       "Uno ganhou",
				Description: fmt.Sprintf("O %s venceu o uno", p.username),
				Color:       16711680,
				Thumbnail: &disgord.EmbedThumbnail{
					URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/UNO_Logo.svg/1024px-UNO_Logo.svg.png",
				},
			}}})
		delete(currentGames, itc.GuildID)
		gameMutex.Unlock()
	}

}
func uno(ctx context.Context, itc *disgord.InteractionCreate, session disgord.Session) {
	gameMutex.RLock()
	defer gameMutex.RUnlock()
	game := currentGames[itc.GuildID]
	if game == nil {
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo na sua guilda",
			},
		})
		return
	}
	if game.status != "n" {
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo na sua guilda",
			},
		})
		return
	}
	var isInGame = isInGame(itc, game)
	if !isInGame {
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Voce nao ta no jogo",
			},
		})
		return
	}
	var p *player
	for _, player2 := range game.players {
		if player2.id == itc.Member.UserID {
			p = player2
			break
		}
	}
	if len(p.cards) == 1 {
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: fmt.Sprintf("O %s ta no uno", itc.Member.User.Username),
			},
		})
		p.uno = true
	} else {
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Voce nao ta no uno nao",
			},
		})
	}
}
func join(ctx context.Context, itc *disgord.InteractionCreate, session disgord.Session) {
	gameMutex.RLock()
	game := currentGames[itc.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Não tem nenhum jogo em preparação nesse servidor, use /uno create para iniciar um",
			},
		})
		return
	}
	if game.status == "n" {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "O jogo ja começou",
			},
		})
		return
	}
	if len(game.players) >= 8 {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Este jogo ja encheu! maximo 8 pessoas",
			},
		})
		return
	}
	var isInGame = isInGame(itc, game)
	if isInGame {
		gameMutex.RUnlock()
		go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Voce ja ta no jogo",
			},
		})
		return
	}
	gameMutex.RUnlock()
	gameMutex.Lock()
	game = currentGames[itc.GuildID]
	game.newPlayer(itc.Member.User)
	gameMutex.Unlock()
	gameMutex.RLock()
	var text string
	for _, player := range game.players {
		text += player.username + "\n"
	}
	gameMutex.RUnlock()
	go itc.Reply(ctx, session, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Color: 16711680,
					Title: "Uno",
					Thumbnail: &disgord.EmbedThumbnail{
						URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/UNO_Logo.svg/1024px-UNO_Logo.svg.png",
					},
					Description: fmt.Sprintf("Use **/uno join** para entrar\nParticipantes:\n%s", text),
				},
			},
		},
	})
}

func runUno(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	command := itc.Data.Options[0].Name
	sess := handler.Client
	switch command {
	case "create":
		create(ctx, itc, sess)
	case "join":
		join(ctx, itc, sess)
	case "buy":
		buy(ctx, itc, sess)
	case "play":
		text := itc.Data.Options[0].Options[0].Value.(string)
		args := strings.Split(text, " ")
		play(ctx, itc, sess, args)
	case "uno":
		uno(ctx, itc, sess)
	case "cards":
		getCards(ctx, itc, sess)
	case "leave":
		leave(ctx, itc, sess)
	}
	return nil
}
