package commands

import (
	"asura/src/cache"
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"strconv"
	"strings"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

var downloadedSprites = []image.Image{}

func init() {
	if len(rinha.Sprites) > 0 {
		for _, sprite := range rinha.Sprites[0] {
			img, err := utils.DownloadImage(sprite)
			if err != nil {
				fmt.Println(sprite)
				fmt.Println(err)
				downloadedSprites = append(downloadedSprites, image.NewRGBA(image.Rect(0, 0, 55, 55)))
				continue
			}
			downloadedSprites = append(downloadedSprites, resize.Resize(55, 55, img, resize.Lanczos3))
		}
	}
	handler.RegisterCommand(handler.Command{
		Name:        "profile",
		Description: translation.T("ProfileHelp", "pt"),
		Run:         runProfile,
		Cooldown:    10,
		Category:    handler.Profile,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Type:        disgord.OptionTypeUser,
			Description: "user profile",
			Required:    true,
		}),
	})
}
func GetClanPos(ctx context.Context, clan *entities.Clan) int {
	if clan == nil || clan.Name == "" {
		return 0
	}
	return database.Clan.GetClanPos(ctx, clan) + 1
}

func runProfile(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := itc.Member.User
	if len(itc.Data.Options) > 0 {
		user = utils.GetUser(itc, 0)
	}
	if user.Bot {
		return nil
	}

	url, _ := user.AvatarURL(128, false)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar, err := utils.DownloadImage(replacer.Replace(url))

	if err != nil {
		avatar, err = utils.DownloadAvatar(user.ID, 128, false)
		if err != nil {
			return nil
		}
	}

	radius := 70.0
	// Resize the images

	galo := database.User.GetUser(ctx, user.ID, "Items", "Galos")
	handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: "Carregando...",
		},
	})
	avatar = resize.Resize(uint(radius*2), uint(radius*2), avatar, resize.Lanczos3)
	img, err := utils.DownloadImage(rinha.GetBackground(&galo))
	if err != nil {
		return nil
	}
	clan := database.Clan.GetUserClan(ctx, user.ID).Clan
	img = resize.Resize(600, 250, img, resize.Lanczos3)
	dc := gg.NewContext(600, 445)

	dc.DrawRoundedRectangle(0, 0, 600, 445, 10)
	dc.Fill()

	dc.DrawRoundedRectangle(0, 0, 600, 200, 10)
	dc.Clip()
	dc.DrawImage(img, 0, 0)
	dc.ResetClip()
	dc.DrawRoundedRectangle(0, 0, 600, 445, 10)
	dc.Clip()

	dc.SetRGB(0.8, 0.31, 0.31)
	dc.DrawRectangle(0, 170, 600, 50)
	dc.Fill()

	dc.SetRGB(0.9, 0.9, 0.9)
	dc.DrawRectangle(0, 220, 200, 285)
	dc.Fill()

	dc.DrawCircle(100, 190, radius+8)
	dc.SetRGB(1, 1, 1)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawCircle(100, 190, radius)
	dc.Clip()
	dc.Fill()

	dc.DrawImage(avatar, int(100-radius), int(190-radius))

	dc.ResetClip()
	dc.SetRGB(0.3, 0.3, 0.3)

	dc.LoadFontFace("./resources/Raleway-Bold.ttf", 13)
	dc.Stroke()

	dc.DrawString("INFO", 220, 375)
	dc.DrawLine(220, 380, 580, 380)
	dc.Stroke()
	dc.LoadFontFace("./resources/Raleway-Bold.ttf", 16)
	dc.DrawStringAnchored("WINS", 260, 400, 0.5, 0.5)
	dc.DrawStringAnchored("LOSES", 350, 400, 0.5, 0.5)
	dc.DrawStringAnchored("CLAN", 440, 400, 0.5, 0.5)
	dc.DrawStringAnchored("GALOS", 530, 400, 0.5, 0.5)

	dc.LoadFontFace("./resources/Raleway-Bold.ttf", 25)
	dc.DrawStringAnchored(strconv.Itoa(galo.Win), 260, 420, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(galo.Lose), 350, 420, 0.5, 0.5)
	dc.DrawStringAnchored("#"+strconv.Itoa(GetClanPos(ctx, clan)), 440, 420, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(len(galo.Galos)), 530, 420, 0.5, 0.5)

	dc.LoadFontFace("./resources/Raleway-Light.ttf", 14)

	dc.DrawString("Clan", 10, 285)
	dc.DrawString("Money", 10, 315)

	clanName := "Nenhum"
	if clan != nil {
		if clan.Name != "" {
			clanName = clan.Name

		}
	}
	dc.DrawStringAnchored(clanName, 190, 285, 1, 0)
	dc.DrawStringAnchored(strconv.Itoa(galo.Money), 190, 315, 1, 0)

	dc.SetRGB(1, 1, 1)
	dc.LoadFontFace("./resources/Raleway-Bold.ttf", 25)
	dc.DrawString(user.Username, 185, 204)

	dc.SetRGB(0.3, 0.3, 0.3)
	dc.DrawLine(10, 328, 190, 328)
	dc.Stroke()

	dc.SetRGB(0.7, 0.7, 0.7)
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			dc.DrawRectangle(float64(10+i*47), float64(345+j*47), 40, 40)
			dc.Fill()
		}
	}
	badgeN := 0
	badges := rinha.GetBadges(&galo)
	if galo.DungeonReset >= 1 {
		badges = append([]*rinha.Cosmetic{rinha.Cosmetics[23]}, badges...)
	}
	if rinha.IsVip(&galo) {
		badges = append([]*rinha.Cosmetic{rinha.Cosmetics[22]}, badges...)
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			if len(badges) <= badgeN {
				break
			}

			badge := badges[badgeN]
			badgeImg, err := utils.DownloadImage(badge.Value)
			badgeN++

			if err != nil {
				fmt.Println(err)
				continue
			}
			badgeImg = resize.Resize(40, 40, badgeImg, resize.Lanczos3)

			dc.DrawImage(badgeImg, 10+i*47, 345+j*47)
		}
	}
	dc.SetRGB(0.9, 0.9, 0.9)
	for i := 0; i < 5; i++ {
		dc.DrawRectangle(float64(220+i*75), 235, 55, 55)
		dc.Fill()
	}
	for i := 0; i < 5; i++ {
		dc.DrawRectangle(float64(220+i*75), 300, 55, 55)
		dc.Fill()
	}
	for i, g := range galo.Galos {
		if i == 10 {
			break
		}
		val := 0.0

		x := i * 75
		if i >= 5 {
			val = 65
			x -= 5 * 75
		}
		rarity := rinha.GetRarity(galo.Galos[i])
		if rarity == rinha.Mythic {
			grad := gg.NewLinearGradient(float64(218+x), 233+val, float64(218+x)+59, 233+val+59)
			grad.AddColorStop(0, color.RGBA{255, 0, 0, 255})
			grad.AddColorStop(0.2, color.RGBA{0, 255, 0, 255})
			grad.AddColorStop(0.4, color.RGBA{0, 0, 255, 255})
			grad.AddColorStop(0.6, color.RGBA{255, 255, 0, 255})
			grad.AddColorStop(0.8, color.RGBA{255, 0, 255, 255})
			grad.AddColorStop(1, color.RGBA{0, 255, 255, 255})
			dc.SetFillStyle(grad)
		} else if galo.Galos[i].Type != -1 {
			if galo.Galos[i].Evolved {
				dc.SetHexColor(fmt.Sprintf("%06x", 0))
			} else {
				colorE := rarity.Color()
				dc.SetHexColor(fmt.Sprintf("%06x", colorE))
			}
		}
		userGaloImg := rinha.GetGaloImage(g, galo.Items)
		var img *image.Image
		if strings.Contains(userGaloImg, "imgur") {
			img = &downloadedSprites[g.Type-1]
		} else {
			img = cache.GetImageFromCache(ctx, "profile", userGaloImg)
			if img == nil {
				imgD, err3 := utils.DownloadImage(userGaloImg)
				if err3 != nil {
					return nil
				}
				imgD = resize.Resize(55, 55, imgD, resize.Lanczos3)
				img = &imgD
				cache.CacheProfileImage(ctx, img, userGaloImg)
			}
		}

		dc.DrawRectangle(float64(218+x), 233+val, 59, 59)
		dc.Fill()
		dc.DrawImage(*img, 220+x, 235+int(val))
	}
	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())
	str := ""

	handler.Client.EditInteractionResponse(ctx, itc, &disgord.UpdateMessage{
		File: &disgord.CreateMessageFile{
			Reader:     bytes.NewReader(b.Bytes()),
			FileName:   "profile.jpg",
			SpoilerTag: false,
		},
		Content: &str,
	})

	return nil
}
