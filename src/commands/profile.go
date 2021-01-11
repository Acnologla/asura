package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"asura/src/database"
	"bytes"
	"context"
	"os"
	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image/png"
	"image"
	"image/jpeg"
	"io"
	"strings"
	"strconv"
)

var downloadedSprites = []image.Image{}

func init() {
	for _, sprite := range rinha.Sprites[0]{
		img, _ := utils.DownloadImage(sprite)
		downloadedSprites = append(downloadedSprites, resize.Resize(55,55, img, resize.Lanczos3))
	}
	handler.Register(handler.Command{
		Aliases:   []string{"perfil", "profile"},
		Run:       runProfile,
		Available: true,
		Cooldown:  10,
		Usage:     "j!profile @user",
		Help:      "Veja o perfil",
		Category: 1,
	})
}



func GetClanPos(clanName string) int{
	q := database.Database.NewRef("clan").OrderByChild("xp")
	result, _ := q.GetOrdered(context.Background())
	for  i := len(result) - 1; 0 <= i; i-- {
		if clanName ==  result[i].Key(){
			return len(result) - i
		}
	} 
	return 0
}

func runProfile(session disgord.Session, msg *disgord.Message, args []string) {
	// Download user image
	user := utils.GetUser(msg,args,session)
	url, _ := user.AvatarURL(128, false)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar, err := utils.DownloadImage(replacer.Replace(url))

	if err != nil {
		msg.Reply(context.Background(), session, "Invalid image")
		return
	}

	radius := 70.0
	// Resize the images
	avatar = resize.Resize(uint(radius*2), uint(radius*2), avatar, resize.Lanczos3)
	file, err := os.Open("resources/wall.jpg")

	if err != nil {
		return
	}
	defer file.Close()
	img, err := jpeg.Decode(file)

	if err != nil {
		return
	}
	dc := gg.NewContext(600, 430)

	dc.DrawRoundedRectangle(0, 0, 600, 430, 10)
	dc.Fill()

	dc.DrawRoundedRectangle(0, 0, 600, 200, 10)
	dc.Clip()
	dc.DrawImage(img,0,0)
	dc.ResetClip()
	dc.DrawRoundedRectangle(0, 0, 600, 430, 10)
	dc.Clip()

	dc.SetRGB(0.8,0.31,0.31)
	dc.DrawRectangle(0, 170, 600, 50)
	dc.Fill()

	dc.SetRGB(0.9,0.9,0.9)
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

	err = dc.LoadFontFace("./resources/Raleway-Bold.ttf", 13)
	dc.DrawString("GALOS", 220, 245)
	dc.DrawLine(220,250,580,250)
	dc.Stroke()

	dc.DrawString("INFO", 220, 345)
	dc.DrawLine(220,350,580,350)
	dc.Stroke()

	galo, _ := rinha.GetGaloDB(user.ID)

	err = dc.LoadFontFace("./resources/Raleway-Bold.ttf", 16)
	dc.DrawStringAnchored("WINS", 260, 380, 0.5, 0.5)
	dc.DrawStringAnchored("LOSES", 350, 380, 0.5, 0.5)
	dc.DrawStringAnchored("CLAN", 440, 380, 0.5, 0.5)
	dc.DrawStringAnchored("GALOS", 530, 380, 0.5, 0.5)

	err = dc.LoadFontFace("./resources/Raleway-Bold.ttf", 25)
	dc.DrawStringAnchored(strconv.Itoa(galo.Win), 260, 400, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(galo.Lose), 350, 400, 0.5, 0.5)
	dc.DrawStringAnchored("#" + strconv.Itoa(GetClanPos(galo.Clan)), 440, 400, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(len(galo.Galos)+1), 530, 400, 0.5, 0.5)

	err = dc.LoadFontFace("./resources/Raleway-Light.ttf", 13)

	dc.DrawString("Clan", 10, 280)
	dc.DrawString("Money", 10, 310)
	clan := "Nenhum"
	if galo.Clan != ""{
		clan = galo.Clan
	}
	dc.DrawStringAnchored(clan, 190, 280, 1, 0)
	dc.DrawStringAnchored(strconv.Itoa(galo.Money), 190, 310, 1, 0)

	dc.SetRGB(1, 1, 1)
	err = dc.LoadFontFace("./resources/Raleway-Bold.ttf", 25)
	dc.DrawString(user.Username, 185, 204)

	dc.SetRGB(0.3,0.3,0.3)
	dc.DrawLine(10,320,190,320)
	dc.Stroke()

	dc.SetRGB(0.7, 0.7, 0.7)
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			dc.DrawRectangle(float64(10 + i*47), float64(332 + j*47), 40, 40)
			dc.Fill()
		}
	}

	dc.SetRGB(0.9, 0.9, 0.9)
	for i := 0; i < 5; i++ {
		dc.DrawRectangle(float64(220 + i*75), 265, 55, 55)
		dc.Fill()
	}
	dc.DrawImage(downloadedSprites[galo.Type-1], 220 , 265)
	for i := range galo.Galos{
		if i == 4 {
			break
		}
		userGaloImg := downloadedSprites[galo.Galos[i].Type - 1]
		dc.DrawImage(userGaloImg, 220 + (i+1)*75, 265)
	}
	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())

	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{
			{bytes.NewReader(b.Bytes()), "profile.jpg", false},
		},
	})
}
