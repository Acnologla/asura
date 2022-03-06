package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"bytes"
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
			img, _ := utils.DownloadImage(sprite)
			downloadedSprites = append(downloadedSprites, resize.Resize(55, 55, img, resize.Lanczos3))
		}
	}
	handler.RegisterCommand(handler.Command{
		Name:        "profile",
		Description: translation.T("ProfileHelp", "pt"),
		Run:         runProfile,
		Cooldown:    15,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Type:        disgord.OptionTypeUser,
			Description: "user profile",
			Required:    true,
		}),
	})
}
func GetClanPos(clan *entities.Clan) int {
	if clan == nil || clan.Name == "" {
		return 0
	}
	return database.Clan.GetClanPos(clan) + 1
}

func runProfile(itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
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

	galo := database.User.GetUser(user.ID, "Items", "Galos")
	avatar = resize.Resize(uint(radius*2), uint(radius*2), avatar, resize.Lanczos3)
	img, err := utils.DownloadImage(rinha.GetBackground(&galo))
	if err != nil {
		return nil
	}
	clan := database.Clan.GetUserClan(user.ID).Clan
	img = resize.Resize(600, 250, img, resize.Lanczos3)
	dc := gg.NewContext(600, 430)

	dc.DrawRoundedRectangle(0, 0, 600, 430, 10)
	dc.Fill()

	dc.DrawRoundedRectangle(0, 0, 600, 200, 10)
	dc.Clip()
	dc.DrawImage(img, 0, 0)
	dc.ResetClip()
	dc.DrawRoundedRectangle(0, 0, 600, 430, 10)
	dc.Clip()

	dc.SetRGB(0.8, 0.31, 0.31)
	dc.DrawRectangle(0, 170, 600, 50)
	dc.Fill()

	dc.SetRGB(0.9, 0.9, 0.9)
	dc.DrawRectangle(0, 220, 200, 220)
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
	dc.DrawString("GALOS", 220, 245)
	dc.DrawLine(220, 250, 580, 250)
	dc.Stroke()

	dc.DrawString("INFO", 220, 345)
	dc.DrawLine(220, 350, 580, 350)
	dc.Stroke()
	dc.LoadFontFace("./resources/Raleway-Bold.ttf", 16)
	dc.DrawStringAnchored("WINS", 260, 380, 0.5, 0.5)
	dc.DrawStringAnchored("LOSES", 350, 380, 0.5, 0.5)
	dc.DrawStringAnchored("CLAN", 440, 380, 0.5, 0.5)
	dc.DrawStringAnchored("GALOS", 530, 380, 0.5, 0.5)

	dc.LoadFontFace("./resources/Raleway-Bold.ttf", 25)
	dc.DrawStringAnchored(strconv.Itoa(galo.Win), 260, 400, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(galo.Lose), 350, 400, 0.5, 0.5)
	dc.DrawStringAnchored("#"+strconv.Itoa(GetClanPos(clan)), 440, 400, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(len(galo.Galos)), 530, 400, 0.5, 0.5)

	dc.LoadFontFace("./resources/Raleway-Light.ttf", 13)

	dc.DrawString("Clan", 10, 280)
	dc.DrawString("Money", 10, 310)
	clanName := "Nenhum"
	if clan != nil {
		if clan.Name != "" {
			clanName = clan.Name

		}
	}
	dc.DrawStringAnchored(clanName, 190, 280, 1, 0)
	dc.DrawStringAnchored(strconv.Itoa(galo.Money), 190, 310, 1, 0)

	dc.SetRGB(1, 1, 1)
	dc.LoadFontFace("./resources/Raleway-Bold.ttf", 25)
	dc.DrawString(user.Username, 185, 204)

	dc.SetRGB(0.3, 0.3, 0.3)
	dc.DrawLine(10, 320, 190, 320)
	dc.Stroke()

	dc.SetRGB(0.7, 0.7, 0.7)
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			dc.DrawRectangle(float64(10+i*47), float64(332+j*47), 40, 40)
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
			badgeImg, _ := utils.DownloadImage(badge.Value)
			badgeImg = resize.Resize(40, 40, badgeImg, resize.Lanczos3)

			badgeN++
			dc.DrawImage(badgeImg, 10+i*47, 332+j*47)
		}
	}
	dc.SetRGB(0.9, 0.9, 0.9)
	for i := 0; i < 5; i++ {
		dc.DrawRectangle(float64(220+i*75), 265, 55, 55)
		dc.Fill()
	}

	for i := range galo.Galos {
		if i == 4 {
			break
		}
		if rinha.Classes[galo.Galos[i].Type].Rarity == rinha.Mythic {
			grad := gg.NewLinearGradient(float64(218+i*75), 263, float64(218+i*75)+59, 263+59)
			grad.AddColorStop(0, color.RGBA{255, 0, 0, 255})
			grad.AddColorStop(0.2, color.RGBA{0, 255, 0, 255})
			grad.AddColorStop(0.4, color.RGBA{0, 0, 255, 255})
			grad.AddColorStop(0.6, color.RGBA{255, 255, 0, 255})
			grad.AddColorStop(0.8, color.RGBA{255, 0, 255, 255})
			grad.AddColorStop(1, color.RGBA{0, 255, 255, 255})
			dc.SetFillStyle(grad)
		} else {
			colorE := rinha.Classes[galo.Galos[i].Type].Rarity.Color()
			dc.SetHexColor(fmt.Sprintf("%06x", colorE))
		}

		userGaloImg := downloadedSprites[galo.Galos[i].Type-1]

		dc.DrawRectangle(float64(218+i*75), 263, 59, 59)
		dc.Fill()
		dc.DrawImage(userGaloImg, 220+i*75, 265)
	}
	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Files: []disgord.CreateMessageFileParams{{
				Reader:     bytes.NewReader(b.Bytes()),
				FileName:   "profile.jpg",
				SpoilerTag: false},
			},
		},
	}

}
