package commands

import (
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils/rinha"
	"bytes"
	"context"
	"fmt"
	"image/png"
	"io"
	"math"
	"strconv"

	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"equipar", "equip", "equipgalo"},
		Run:       runEquip,
		Available: true,
		Cooldown:  2,
		Usage:     "j!equipar",
		Help:      "Equipe outro galo",
		Category:  1,
	})
}

func runEquip(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if len(args) == 0 {
		lenGalos := len(galo.Galos)
		dc := gg.NewContext(510, 95*int(math.Round((float64(lenGalos)/2))))
		for i, _galo := range galo.Galos {
			dc.LoadFontFace("./resources/Raleway-Bold.ttf", 30)
			galo := _galo.Type
			img := downloadedSprites[galo-1]
			dc.SetRGB255(255, 255, 255)
			dc.DrawRectangle(0, float64(i)*95, 510, 95)
			width := 0.0
			if i%2 != 0 {
				width = 250
			}
			dc.Fill()
			dc.SetRGB255(30, 30, 30)
			dc.DrawString(strconv.Itoa(i), 20+width, 55+float64(i/2)*95)
			if i != 0 && i != lenGalos {
				dc.DrawLine(0, 95*float64(i/2), 510, 95*float64(i/2))
				dc.Stroke()
			}

			dc.LoadFontFace("./resources/Raleway-Light.ttf", 15)
			name := rinha.Classes[galo].Name
			if _galo.Name != "" {
				name = _galo.Name
			}
			dc.DrawString(name, 70+57+width, 40+float64(i/2)*95)
			dc.DrawString(fmt.Sprintf("Level: %d", rinha.CalcLevel(_galo.Xp)), 70+57+width, 60+float64(i/2)*95)
			color := rinha.Classes[galo].Rarity.Color()
			dc.SetHexColor(fmt.Sprintf("%06x", color))
			dc.DrawRectangle(58+width, 18+float64(i/2)*95, 59, 59)
			dc.Fill()
			dc.DrawImage(img, 60+int(width), 20+(i/2)*95)

		}
		var b bytes.Buffer
		pw := io.Writer(&b)
		png.Encode(pw, dc.Image())
		avatar, _ := msg.Author.AvatarURL(512, true)
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Files: []disgord.CreateMessageFileParams{
				{bytes.NewReader(b.Bytes()), "galos.jpg", false},
			},
			Embed: &disgord.Embed{
				Footer: &disgord.EmbedFooter{
					IconURL: avatar,
					Text:    "Use j!equipar <numero do galo> para equipar um galo | use j!equipar <numero do galo> remove para vender um galo",
				},
				Color: 65535,
				Image: &disgord.EmbedImage{
					URL: "attachment://galos.jpg",
				},
			},
		})
	} else {
		battleMutex.RLock()
		if currentBattles[msg.Author.ID] != "" {
			battleMutex.RUnlock()
			msg.Reply(context.Background(), session, "Espere sua rinha terminar para equipar ou vender galos")
			return
		}
		battleMutex.RUnlock()
		value, err := strconv.Atoi(args[0])
		if err != nil {
			msg.Reply(context.Background(), session, "Use j!equipar <numero do galo> para equipar um galo | use j!equipar <numero do galo> remove para vender um galo")
			return
		}
		if value >= 0 && len(galo.Galos) > value {
			if len(args) >= 2 {
				if args[1] == "remove" || args[1] == "vender" {
					var gal rinha.SubGalo
					var price int
					rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
						gal = galo.Galos[value]
						for i := value; i < len(galo.Galos)-1; i++ {
							galo.Galos[i] = galo.Galos[i+1]
						}
						price = rinha.Sell(rinha.Classes[gal.Type].Rarity, gal.Xp)
						galo.Money += price
						galo.Galos = galo.Galos[0 : len(galo.Galos)-1]
						return galo, nil
					})
					newGalo := rinha.Classes[gal.Type]
					tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
					telemetry.Debug(fmt.Sprintf("%s Sell %s", tag, newGalo.Name), map[string]string{
						"galo":   newGalo.Name,
						"user":   strconv.FormatUint(uint64(msg.Author.ID), 10),
						"rarity": newGalo.Rarity.String(),
					})
					msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce vendeu o galo **%s** por **%d** de dinheiro com sucesso", msg.Author.Mention(), rinha.Classes[gal.Type].Name, price))
					return
				}
			}
			rinha.UpdateGaloDB(msg.Author.ID, func(galo rinha.Galo) (rinha.Galo, error) {
				newGalo := galo.Galos[value]
				old := rinha.SubGalo{
					Xp:        galo.Xp,
					Type:      galo.Type,
					Name:      galo.Name,
					GaloReset: galo.GaloReset,
				}
				galo.Type = newGalo.Type
				galo.Xp = newGalo.Xp
				galo.GaloReset = newGalo.GaloReset
				galo.Name = newGalo.Name
				galo.Galos[value] = old
				galo.Equipped = []int{}
				msg.Reply(context.Background(), session, fmt.Sprintf("%s, Voce trocou seu galo **%s** por **%s**", msg.Author.Mention(), rinha.Classes[old.Type].Name, rinha.Classes[galo.Type].Name))
				return galo, nil
			})

		} else {
			msg.Reply(context.Background(), session, "Numero invalido")

		}
	}
}
