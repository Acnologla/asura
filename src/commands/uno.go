package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"bytes"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image/png"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

var cards []string = []string{}

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
	if err := dc.LoadFontFace("./resources/Raleway-Black.ttf", 20); err != nil {
	}
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
	handler.Register(handler.Command{
		Aliases:   []string{"uno", "unoplay", "u", "playuno"},
		Run:       runUno,
		Available: true,
		Cooldown:  2,
		Usage:     "j!uno <ação>",
		Help:      "Jogue uno",
	})
	colors := []string{"B", "R", "Y", "G"}
	rand.Seed(time.Now().UnixNano())
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
func shuffle(arr []string) []string {
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func getCards(msg *disgord.Message, session disgord.Session) {
	gameMutex.RLock()
	defer gameMutex.RUnlock()
	game := currentGames[msg.GuildID]
	if game == nil {
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Não tem nenhum jogo nessa guilda, use j!uno create para criar um")
		return
	}
	if game.status == "n" {
		var isInGame = isInGame(msg, game)
		if isInGame {
			var p *player
			for _, player2 := range game.players {
				if player2.id == msg.Author.ID {
					p = player2
					break
				}
			} 
			b := bytes.NewReader(game.drawCards(p).Bytes())
			chanID, err := handler.Client.CreateDM(context.Background(), msg.Author.ID)
			if err == nil {
				go handler.Client.SendMsg(context.Background(), chanID.ID, &disgord.CreateMessageParams{
					Content: "Para jogar uma carta usa j!uno play <id da carta> ( o id é o coiso que ta encima da carta)",
					Files: []disgord.CreateMessageFileParams{
						{b, "uno.png", false},
					},
				})
			}
		} else {
			go msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce não ta no jogo")
		}
	}
}

func sendEmbed(msg *disgord.Message, session disgord.Session) {
	gameMutex.RLock()
	defer gameMutex.RUnlock()
	game := currentGames[msg.GuildID]
	go msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{
			{bytes.NewReader(game.getBoard().Bytes()), "uno.png", false},
		},
		Embed: &disgord.Embed{
			Title: "Uno",
			Footer: &disgord.EmbedFooter{
				Text: fmt.Sprintf("Turno do: %s", game.players[0].username),
			},
			Description: "Use j!uno play <carta> para jogar!\nj!uno leave (para sair)\nVeja suas cartas no privado (usando j!cartas)\nj!uno buy (para comprar uma carta e passar o turno\n caso não jogue ou não compre durante 3 minutos\n ira passar o turno e comprar automatico)\nCaso estiver no uno (uma carta) use j!uno uno",
			Image: &disgord.EmbedImage{
				URL: "attachment://uno.png",
			},
		},
	})
}

func create(msg *disgord.Message, session disgord.Session) {
	gameMutex.RLock()
	game := currentGames[msg.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		gameMutex.Lock()
		gameCards := append([]string{}, shuffle(cards)...)
		currentGames[msg.GuildID] = &unoGame{
			cards:  gameCards,
			board:  []string{},
			status: "p",
			lastPlay: time.Now(),
		}
		currentGames[msg.GuildID].newPlayer(msg.Author)
		gameMutex.Unlock()
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title: "Uno",
				Thumbnail: &disgord.EmbedThumbnail{
					URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/UNO_Logo.svg/1024px-UNO_Logo.svg.png",
				},
				Description: fmt.Sprintf("Use **j!uno join** para entrar ( O jogo começara em **1** minuto)\nParticipantes:\n%s", msg.Author.Username),
				Color:       16711680,
			},
		})
		go func() {
			time.Sleep(60 * time.Second)
			gameMutex.RLock()
			game = currentGames[msg.GuildID]
			if 1 >= len(game.players) {
				gameMutex.RUnlock()
				go msg.Reply(context.Background(), session, "Como só teve um jogador, o uno foi cancelado, uno j!uno create para jogar novamente")
				gameMutex.Lock()
				delete(currentGames, msg.GuildID)
				gameMutex.Unlock()
			} else {
				go func() {
					for {
						time.Sleep(60 * time.Second)
						gameMutex.RLock()
						game := currentGames[msg.GuildID]
						if game != nil {
							if time.Since(game.lastPlay).Seconds()/60 >= 2 {
								if game.refuse == 4 {
									gameMutex.RUnlock()
									gameMutex.Lock()
									delete(currentGames, msg.GuildID)
									gameMutex.Unlock()
									msg.Reply(context.Background(), session, "4 rodadas sem jogar consecutivas, jogo acabado")
									return
								} else {
									go msg.Reply(context.Background(), session, fmt.Sprintf("O %s pulou seu turno e comprou uma carta agora é o turno de %s", game.players[0].username, game.players[1].username))
									gameMutex.RUnlock()
									gameMutex.Lock()
									game = currentGames[msg.GuildID]
									game.buy(game.players[0], 1)
									game.addToLast()
									game.lastPlay = time.Now()
									game.refuse++
									gameMutex.Unlock()
									continue
								}
							}
						} else {
							gameMutex.RUnlock()
							return
						}
					}
				}()
				for _, player2 := range game.players {
					go func(p *player) {
						chanID, err := handler.Client.CreateDM(context.Background(), p.id)
						if err == nil {
							handler.Client.SendMsg(context.Background(), chanID.ID, &disgord.CreateMessageParams{
								Content: "Para jogar uma carta usa j!uno play <id da carta> ( o id é o coiso que ta encima da carta)",
								Files: []disgord.CreateMessageFileParams{
									{bytes.NewReader(game.drawCards(p).Bytes()), "uno.png", false},
								},
							})
						}
					}(player2)
				}
				gameMutex.RUnlock()
				gameMutex.Lock()
				game.start()
				gameMutex.Unlock()
				sendEmbed(msg, session)
			}
		}()
	} else if game.status == "p" {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, "O jogo em sua guilda esta em preparação, caso queira entrer use j!uno join")
	} else {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, "O jogo da sua guilda ja começou, espere ele terminar")
	}
}
func isInGame(msg *disgord.Message, game *unoGame) bool {
	var isInGame bool
	for _, player := range game.players {
		if player.id == msg.Author.ID {
			isInGame = true
			break
		}
	}
	return isInGame
}
func leave(msg *disgord.Message, session disgord.Session) {
	gameMutex.RLock()
	game := currentGames[msg.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, "Não tem nenhum jogo na sua guilda")
		return
	}
	if game.status != "n" {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, "Não tem nenhum jogo na sua guilda")
		return
	}
	var isInGame = isInGame(msg, game)
	if !isInGame {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce não ta no jogo")
		return
	}
	gameMutex.RUnlock()
	gameMutex.Lock()
	game.removePlayer(msg.Author)
	game = currentGames[msg.GuildID]
	gameMutex.Unlock()
	if 1 >= len(game.players) {
		go msg.Reply(context.Background(), session, "Como todos os jogadores sairam ou só sobrou 1 o jogo foi encerrado")
		gameMutex.Lock()
		delete(currentGames, msg.GuildID)
		gameMutex.Unlock()
		return
	}
	msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce saiu do jogo agora é o turno do "+game.players[0].username)
}

func buy(msg *disgord.Message, session disgord.Session) {
	gameMutex.RLock()
	game := currentGames[msg.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, "Não tem nenhum jogo na sua guilda")
		return
	}
	if game.status != "n" {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, "Não tem nenhum jogo na sua guilda")
		return
	}
	var isInGame = isInGame(msg, game)
	if !isInGame {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao ta no jogo")
		return
	}
	gameMutex.RUnlock()
	gameMutex.Lock()
	var p *player
	for _, player2 := range game.players {
		if player2.id == msg.Author.ID {
			p = player2
			break
		}
	}
	if p.id != game.players[0].id {
		gameMutex.Unlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Não é seu turno")
		return
	}
	game.buy(p, 1)
	game.addToLast()
	game.lastPlay = time.Now()
	gameMutex.Unlock()
	getCards(msg, session)
	sendEmbed(msg, session)
}

func play(msg *disgord.Message, session disgord.Session, args []string) {
	gameMutex.RLock()
	game := currentGames[msg.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, "Não tem nenhum jogo na sua guilda")
		return
	}
	if game.status != "n" {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, "Não tem nenhum jogo na sua guilda")
		return
	}
	if len(args) == 0 {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Diga uma carta para jogar")
		return
	}
	var isInGame = isInGame(msg, game)
	if !isInGame {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao ta no jogo")
		return
	}
	if strings.HasPrefix(args[0], "WILD") && 2 > len(args) {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Essa é uma carta preta usa  então voce tem que definir a cor apos a carta exemplo:\nj!uno play wild <Y (Amarelo) | G (Verde) | B (Azul) | R (Vermelho) >")
		return
	}
	gameMutex.RUnlock()
	gameMutex.Lock()
	game = currentGames[msg.GuildID]
	var p *player
	for _, player2 := range game.players {
		if player2.id == msg.Author.ID {
			p = player2
			break
		}
	}
	if p.id != game.players[0].id {
		gameMutex.Unlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Não é seu turno")
		return
	}
	play := game.play(p, args)
	if play != "" {
		gameMutex.Unlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", "+play)
		return
	}
	game = currentGames[msg.GuildID]
	game.refuse = 0
	game.lastPlay = time.Now()
	gameMutex.Unlock()
	sendEmbed(msg, session)
	if len(p.cards) == 0 {
		gameMutex.Lock()
		if !p.uno {
			game.buy(p, 2)
			gameMutex.Unlock()
			go msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce comprou 2 cartas por não ter falado uno")
			getCards(msg, session)
			return
		}
		go msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       "Uno ganhou",
				Description: fmt.Sprintf("O %s venceu o uno", p.username),
				Color:       16711680,
				Thumbnail: &disgord.EmbedThumbnail{
					URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/UNO_Logo.svg/1024px-UNO_Logo.svg.png",
				},
			},
		})
		delete(currentGames, msg.GuildID)
		gameMutex.Unlock()
	}

}
func uno(msg *disgord.Message, session disgord.Session) {
	gameMutex.RLock()
	defer gameMutex.RUnlock()
	game := currentGames[msg.GuildID]
	if game == nil {
		go msg.Reply(context.Background(), session, "Não tem nenhum jogo na sua guilda")
		return
	}
	if game.status != "n" {
		go msg.Reply(context.Background(), session, "Não tem nenhum jogo na sua guilda")
		return
	}
	var isInGame = isInGame(msg, game)
	if !isInGame {
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao ta no jogo")
		return
	}
	var p *player
	for _, player2 := range game.players {
		if player2.id == msg.Author.ID {
			p = player2
			break
		}
	}
	if len(p.cards) == 1 {
		go msg.Reply(context.Background(), session, fmt.Sprintf("O %s ta no uno", msg.Author.Username))
		p.uno = true
	} else {
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao ta no uno nao")
	}
}
func join(msg *disgord.Message, session disgord.Session) {
	gameMutex.RLock()
	game := currentGames[msg.GuildID]
	if game == nil {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Não tem nenhum jogo em preparação nesse servidor, use j!uno create para iniciar um")
		return
	}
	if game.status == "n" {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", O jogo ja começou")
		return
	}
	if len(game.players) >= 8 {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Este jogo ja encheu! maximo 8 pessoas")
		return
	}
	var isInGame = isInGame(msg, game)
	if isInGame {
		gameMutex.RUnlock()
		go msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce ja ta no jogo")
		return
	}
	gameMutex.RUnlock()
	gameMutex.Lock()
	game = currentGames[msg.GuildID]
	game.newPlayer(msg.Author)
	gameMutex.Unlock()
	var text string
	for _, player := range game.players {
		text += player.username + "\n"
	}
	go msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Color: 16711680,
			Title: "Uno",
			Thumbnail: &disgord.EmbedThumbnail{
				URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/UNO_Logo.svg/1024px-UNO_Logo.svg.png",
			},
			Description: fmt.Sprintf("Use **j!uno join** para entrar\nParticipantes:\n%s", text),
		},
	})
}

func runUno(session disgord.Session, msg *disgord.Message, args []string) {
	if len(args) == 0 {
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Content: msg.Author.Mention(),
			Embed: &disgord.Embed{
				Color:       16711680,
				Title:       "Uno",
				Description: "Uso certo:\nj!uno play carta (Joga uma carta na mesa caso for seu turno)\nj!uno buy ( compre uma carta e pule seu turno caso não tiver como jogar, pulara automatico apos 3 minutos sem jogar no seu turno)\nj!uno cartas (veja suas cartas)\nj!uno uno (Fale que esta no uno)\nj!uno create (cria um jogo no seu servidor)\nj!uno leave (caso queira sair)\nj!uno join (entra no jogo do seu servidor caso ainda não tiver começado)",
			},
		})
	} else {
		function := strings.ToLower(args[0])
		args = args[1:]
		if len(args) >= 1 {
			args[0] = strings.ToUpper(args[0])
		}
		if len(args) >= 2 {
			args[1] = strings.ToUpper(args[1])
		}
		if function == "create" {
			create(msg, session)
		} else if function == "join" {
			join(msg, session)
		} else if function == "leave" || function == "sair" {
			leave(msg, session)
		} else if function == "cartas" || function == "cards" {
			getCards(msg, session)
		} else if function == "uno" {
			uno(msg, session)
		} else if function == "play" || function == "jogar" {
			if len(args) == 0 {
				msg.Reply(context.Background(), session, msg.Author.Mention()+", Diga uma carta para jogar, j!uno play <id da carta>")
				return
			}
			play(msg, session, args)
		} else if function == "buy" || function == "comprar" {
			buy(msg, session)
		} else {
			msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
				Content: msg.Author.Mention(),
				Embed: &disgord.Embed{
					Color:       16711680,
					Title:       "Uno",
					Description: "Uso certo:\nj!uno play carta (Joga uma carta na mesa caso for seu turno)\nj!uno buy ( compre uma carta e pule seu turno caso não tiver como jogar, pulara automatico apos 3 minutos sem jogar no seu turno)\nj!uno cartas (veja suas cartas)\nj!uno uno (Fale que esta no uno)\nj!uno create (cria um jogo no seu servidor)\nj!uno leave (caso queira sair)\nj!uno join (entra no jogo do seu servidor caso ainda não tiver começado)",
				},
			})
			return
		}
	}
}
