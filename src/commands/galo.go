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

	"github.com/fogleman/gg"

	"io"
	"strconv"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/nfnt/resize"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "galo",
		Description: translation.T("GaloHelp", "pt"),
		Run:         runGalo,
		Cooldown:    15,
		Category:    handler.Profile,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Type:        disgord.OptionTypeUser,
			Description: "user galo",
			Required:    true,
		}),
	})
}

var evolvedPrice = 10

func runGalo(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := itc.Member.User
	if len(itc.Data.Options) > 0 {
		user = utils.GetUser(itc, 0)
	}
	if user.Bot {
		return nil
	}
	u := database.User.GetUser(ctx, user.ID, "Items", "Galos")
	galo := rinha.GetEquippedGalo(&u)
	skills := rinha.GetEquipedSkills(galo)
	avatar, err := utils.DownloadImage(rinha.GetGaloImage(galo, u.Items))

	if err != nil {
		return nil
	}

	name := "Galo do " + user.Username
	if galo.Name != "" {
		name = galo.Name
	}

	radius := 50.0
	// Resize the images
	avatar = resize.Resize(uint(radius*2), uint(radius*2), avatar, resize.Lanczos3)

	img, err := utils.DownloadImage(rinha.GetBackground(&u))
	if err != nil {

		fmt.Println(err)
		return nil
	}
	img = resize.Resize(320, 200, img, resize.Lanczos3)
	dc := gg.NewContext(320, 450)

	dc.DrawRoundedRectangle(0, 0, 320, 630, 10)
	dc.Fill()

	dc.DrawRoundedRectangle(0, 0, 320, 150, 10)
	dc.Clip()
	dc.DrawImage(img, 0, 0)
	dc.ResetClip()

	dc.DrawRoundedRectangle(0, 0, 320, 630, 10)
	dc.Clip()

	dc.SetRGB(0.8, 0.31, 0.31)
	dc.DrawRectangle(0, 140, 320, 40)
	dc.Fill()

	dc.DrawCircle(160, 70, radius+3)
	//	dc.SetRGB(1, 1, 1)
	if rinha.Classes[galo.Type].Rarity == rinha.Mythic {
		grad := gg.NewConicGradient(160, 70, radius+3)
		grad.AddColorStop(0, color.RGBA{255, 0, 0, 255})
		grad.AddColorStop(0.2, color.RGBA{0, 255, 0, 255})
		grad.AddColorStop(0.4, color.RGBA{0, 0, 255, 255})
		grad.AddColorStop(0.6, color.RGBA{255, 255, 0, 255})
		grad.AddColorStop(0.8, color.RGBA{255, 0, 255, 255})
		grad.AddColorStop(1, color.RGBA{0, 255, 255, 255})
		dc.SetFillStyle(grad)
	} else if galo.Evolved {
		dc.SetHexColor(fmt.Sprintf("%06x", 0))
	} else {
		color := rinha.Classes[galo.Type].Rarity.Color()
		dc.SetHexColor(fmt.Sprintf("%06x", color))
	}

	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawCircle(160, 70, radius)
	dc.Clip()
	dc.Fill()

	dc.DrawImage(avatar, int(160-radius), int(70-radius))

	dc.ResetClip()
	dc.SetRGB(0.3, 0.3, 0.3)

	dc.LoadFontFace("./resources/Raleway-Bold.ttf", 13)
	dc.SetRGB255(196, 196, 196)
	dc.DrawRoundedRectangle(10, 195, 300, 20, 10)
	dc.Fill()
	dc.SetRGB255(208, 80, 80)
	level := rinha.CalcLevel(galo.Xp)
	curLevelXP := float64(rinha.CalcXP(level))
	nextLevelXp := float64(rinha.CalcXP(level + 1))
	percentage := (float64(galo.Xp-int(curLevelXP)) * 100) / (nextLevelXp - curLevelXP)
	if percentage >= 3 {
		dc.DrawRoundedRectangle(10, 195, 300*(percentage/100), 20, 10)
		dc.Fill()
	}
	dc.SetRGB(0.3, 0.3, 0.3)
	dc.DrawStringAnchored("HABILIDADES EQUIPADAS", 320/2, 295, 0.5, 0.5)
	dc.DrawLine(10, 310, 310, 310)
	dc.Stroke()

	dc.LoadFontFace("./resources/Raleway-Light.ttf", 14)
	dc.SetRGB255(255, 255, 255)
	dc.DrawStringAnchored(fmt.Sprintf("%d/%d", galo.Xp-int(curLevelXP), int(nextLevelXp-curLevelXP)), 320/2, 203.5, 0.5, 0.5)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Tipo", 10, 240)
	dc.DrawString("Level", 10, 255)
	dc.DrawString("Item", 10, 270)
	levelText := strconv.Itoa(rinha.CalcLevel(galo.Xp))
	if galo.Resets > 0 {
		levelText += "["
		levelText += strconv.Itoa(galo.Resets)
		levelText += "]"
	}
	dc.DrawStringAnchored(rinha.Classes[galo.Type].Name, 310, 240, 1, 0)
	dc.DrawStringAnchored(levelText, 310, 255, 1, 0)
	item := rinha.GetEquippedItem(&u)
	if item != -1 {
		i := rinha.Items[item]
		dc.DrawStringAnchored(i.Name, 310, 270, 1, 0)
	} else {
		dc.DrawStringAnchored("Nenhum", 310, 270, 1, 0)

	}

	for i, skill := range skills {
		rinhaSkill := rinha.Skills[galo.Type-1][skill.Skill]
		margin := float64(335 + (25 * i))
		text, _ := rinha.SkillToString(rinhaSkill)
		dc.DrawString(text, 10, margin)
		min, max := rinha.CalcDamage(rinhaSkill, galo)
		dc.DrawStringAnchored(fmt.Sprintf("Dano: %d - %d", min, max-1), 310, margin, 1, 0)
	}

	dc.SetRGB(1, 1, 1)
	dc.LoadFontFace("./resources/Raleway-Light.ttf", 22)
	dc.DrawStringAnchored(name, 320/2, 160, 0.5, 0.5)

	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())
	data := &disgord.CreateInteractionResponseData{
		Files: []disgord.CreateMessageFile{{
			Reader:     bytes.NewReader(b.Bytes()),
			FileName:   "galo.jpg",
			SpoilerTag: false},
		},
	}
	if level >= 35 {
		data.Components = []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Label:    "Upgrade",
						CustomID: "upgrade",
						Style:    disgord.Primary,
					},
				},
			},
		}
	}
	if !galo.Evolved && galo.Resets > 4 && u.AsuraCoin > evolvedPrice {
		component := &disgord.MessageComponent{
			Type:     disgord.MessageComponentButton,
			Label:    "Evoluir",
			CustomID: "Evolve",
			Style:    disgord.Primary,
		}

		if len(data.Components) == 1 {
			data.Components[0].Components = append(data.Components[0].Components, component)
		} else {
			data.Components = []*disgord.MessageComponent{
				{
					Type: disgord.MessageComponentActionRow,
					Components: []*disgord.MessageComponent{
						component,
					},
				},
			}
		}

	}
	if len(data.Components) > 0 {
		handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: data,
		})
		handler.RegisterHandler(itc.ID, func(ic *disgord.InteractionCreate) {
			done := false
			isEvolved := ic.Data.CustomID == "Evolve"
			if ic.Member.UserID == user.ID {
				database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
					database.User.UpdateEquippedRooster(ctx, u, func(r entities.Rooster) entities.Rooster {
						if isEvolved {
							if u.AsuraCoin > evolvedPrice && r.Resets > 4 && !r.Evolved {
								r.Evolved = true
								u.AsuraCoin -= evolvedPrice
								done = true
							}
						} else {
							level := rinha.CalcLevel(r.Xp)
							if level >= 35 {
								done = true
								r.Equipped = []int{}
								r.Xp = 0
								r.Resets++
							}
						}
						return r
					})
					return u
				}, "Galos")
				if done {
					if isEvolved {
						ic.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
							Type: disgord.InteractionCallbackChannelMessageWithSource,
							Data: &disgord.CreateInteractionResponseData{
								Content: "Galo evoluido com sucesso",
							},
						})
					} else {
						ic.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
							Type: disgord.InteractionCallbackChannelMessageWithSource,
							Data: &disgord.CreateInteractionResponseData{
								Content: translation.T("GaloUpgrade", translation.GetLocale(itc)),
							},
						})
					}
				}
			}
		}, 100)
	} else {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: data,
		}
	}
	return nil
}
