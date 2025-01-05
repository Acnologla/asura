package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"bytes"
	"context"
	"fmt"
	"image/color"
	"image/png"
	"io"
	"math"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "equip",
		Description: translation.T("RunEquip", "pt"),
		Run:         runEquip,
		Cooldown:    8,
		Category:    handler.Profile,
	})
}

func genEquipOptions(user *entities.User) (opts []*disgord.SelectMenuOption) {
	for _, galo := range user.Galos {
		if !galo.Equip {
			class := rinha.Classes[galo.Type]
			opts = append(opts, &disgord.SelectMenuOption{
				Label:       class.Name,
				Value:       galo.ID.String(),
				Description: fmt.Sprintf("Equipar galo %s (level %d)", class.Name, rinha.CalcLevel(galo.Xp)),
			})
		}
	}

	if len(opts) == 0 {
		opts = append(opts, &disgord.SelectMenuOption{
			Label:       "Nenhum galo",
			Value:       "nil",
			Description: "Nenhum galo para equipar",
		})
	}
	return
}

func runEquip(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	galo := database.User.GetUser(ctx, itc.Member.UserID, "Galos", "Items")
	optsGalos := genEquipOptions(&galo)
	mythics := rinha.GetFromType(rinha.Mythic)
	hasAllMythics := len(galo.Galos) > 0
	for _, mythic := range mythics {
		if !rinha.HaveRooster(galo.Galos, mythic) {
			hasAllMythics = false
			break
		}
	}
	if hasAllMythics {
		completeAchievement(ctx, itc, 9)
	}

	width, height := 630, 119*math.Round(float64(len(galo.Galos))/2)+12
	dc := gg.NewContext(width, int(height))
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()
	radius := 50
	imageSize := radius * 2
	handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: "Carregando...",
		},
	})
	for i, g := range galo.Galos {
		image := rinha.GetGaloImage(g, galo.Items)
		imageD, err := utils.DownloadImage(image)
		if err != nil {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Tente novamente",
				},
			}
		}
		ga := rinha.Classes[g.Type]
		img := resize.Resize(uint(imageSize), uint(imageSize), imageD, resize.Lanczos3)
		y := i / 2
		x := i%2*330 + 20
		imageVerticalMargin := 20 + y*imageSize + 15*y
		rarity := rinha.GetRarity(galo.Galos[i])
		if rarity == rinha.Mythic {
			grad := gg.NewConicGradient(float64(x+radius), float64(imageVerticalMargin+radius), float64(radius+3))
			grad.AddColorStop(0, color.RGBA{255, 0, 0, 255})
			grad.AddColorStop(0.2, color.RGBA{0, 255, 0, 255})
			grad.AddColorStop(0.4, color.RGBA{0, 0, 255, 255})
			grad.AddColorStop(0.6, color.RGBA{255, 255, 0, 255})
			grad.AddColorStop(0.8, color.RGBA{255, 0, 255, 255})
			grad.AddColorStop(1, color.RGBA{0, 255, 255, 255})
			dc.SetFillStyle(grad)
		} else if g.Evolved {
			dc.SetHexColor(fmt.Sprintf("%06x", 0))
		} else {
			color := rarity.Color()
			dc.SetHexColor(fmt.Sprintf("%06x", color))
		}
		dc.DrawCircle(float64(x+radius), float64(imageVerticalMargin+radius), float64(radius+3))
		dc.Fill()
		dc.SetRGB(0, 0, 0)
		//radius := float64(imageSize) / (math.Sqrt(math.Pi))
		dc.DrawCircle(float64(x+radius), float64(imageVerticalMargin+radius), float64(radius))
		dc.Clip()
		dc.Fill()
		dc.DrawImage(img, x, imageVerticalMargin)
		dc.ResetClip()
		dc.SetRGB(0, 0, 0)
		textMargin := float64(imageSize+20) + float64(x)
		verticalMargin := float64(imageVerticalMargin) + 40
		dc.LoadFontFace("./resources/Raleway-Bold.ttf", 18)
		dc.DrawString(rinha.GetNameVanilla(ga.Name, *g), textMargin, verticalMargin)
		dc.LoadFontFace("./resources/Raleway-Light.ttf", 19)
		xp := fmt.Sprintf("Level: %d", rinha.CalcLevel(g.Xp))
		dc.DrawString(xp, textMargin, verticalMargin+20)

	}

	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())
	embed := &disgord.Embed{
		Color: 65535,
		Title: "Equip",
		Image: &disgord.EmbedImage{
			URL: "attachment://equip.png",
		},
	}
	r := entities.CreateMsg().
		Embed(embed).
		Component(entities.CreateComponent().
			Select(
				translation.T("EquipGaloPlaceholder", translation.GetLocale(itc)),
				"galoEquip",
				optsGalos))
	str := ""
	handler.Client.EditInteractionResponse(ctx, itc, &disgord.UpdateMessage{
		File: &disgord.CreateMessageFile{
			Reader:     bytes.NewReader(b.Bytes()),
			FileName:   "profile.jpg",
			SpoilerTag: false,
		},

		Components: &r.Data.Components,
		Content:    &str,
	})

	handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
		userIC := ic.Member.User
		name := ic.Data.CustomID
		if userIC.ID != itc.Member.UserID {
			return
		}
		if len(ic.Data.Values) == 0 {
			return
		}
		val := ic.Data.Values[0]
		if val == "nil" {
			return
		}
		itemID := uuid.MustParse(val)
		msg := ""
		database.User.UpdateUser(ctx, userIC.ID, func(u entities.User) entities.User {
			if isInRinha(ctx, userIC) != "" {
				msg = "IsInRinha"
				return u
			}
			if name == "galoEquip" {
				database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
					database.User.UpdateRooster(ctx, &u, itemID, func(r2 entities.Rooster) entities.Rooster {
						if r2.Type != 0 && r2.ID != r.ID {
							r.Equip = false
							r2.Equip = true
						}
						return r2
					})
					return r
				})
				msg = "EquipGalo"
			}
			return u
		}, "Galos")
		if msg != "" {
			handler.Client.SendInteractionResponse(ctx, ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: translation.T(msg, translation.GetLocale(ic)),
				},
			})
		}
	}, 120)
	return nil
}
