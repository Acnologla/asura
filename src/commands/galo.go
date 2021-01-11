package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"asura/src/utils/rinha"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image/png"
	"io"
	"fmt"
	"bytes"
	"os"
	"strings"
	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"galo", "galolevel", "meugalo"},
		Run:       runGalo,
		Available: true,
		Cooldown:  3,
		Usage:     "j!galo",
		Help:      "Informação sobre seu galo",
		Category:  1,
	})
}

func runGalo(session disgord.Session, msg *disgord.Message, args []string) {
	user := utils.GetUser(msg, args, session)
	galo, _ := rinha.GetGaloDB(user.ID)

	url, _ := user.AvatarURL(128, false)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar, err := utils.DownloadImage(replacer.Replace(url))

	if err != nil {
		msg.Reply(context.Background(), session, "Invalid image")
		return
	}

	name := "Galo do " + user.Username
	if galo.Name != "" {
		name = galo.Name
	}

	radius := 50.0
	// Resize the images
	avatar = resize.Resize(uint(radius*2), uint(radius*2), avatar, resize.Lanczos3)
	
	file, err := os.Open("resources/wall.png")
	
	if err != nil {
		return
	}

	defer file.Close()
	img, err := png.Decode(file)

	if err != nil {
		return
	}

	dc := gg.NewContext(320, 430)

	dc.DrawRoundedRectangle(0, 0, 320, 630, 10)
	dc.Fill()

	dc.DrawRoundedRectangle(0, 0, 320, 150, 10)
	dc.Clip()
	dc.DrawImage(img,0,0)
	dc.ResetClip()

	dc.DrawRoundedRectangle(0, 0, 320, 630, 10)
	dc.Clip()

	dc.SetRGB(0.8,0.31,0.31)
	dc.DrawRectangle(0, 140, 320, 40)
	dc.Fill()

	dc.DrawCircle(160, 70, radius+3)
	dc.SetRGB(1, 1, 1)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawCircle(160, 70, radius)
	dc.Clip()
	dc.Fill()

	dc.DrawImage(avatar, int(160-radius), int(70-radius))

	dc.ResetClip()
	dc.SetRGB(0.3, 0.3, 0.3)

	err = dc.LoadFontFace("./resources/Raleway-Bold.ttf", 13)
	dc.DrawString("INFO", 10, 205)
	dc.DrawLine(10,210,310,210)
	dc.Stroke()

	dc.DrawString("ATAQUES", 10, 295)
	dc.DrawLine(10,300,310,300)
	dc.Stroke()


	err = dc.LoadFontFace("./resources/Raleway-Light.ttf", 14)

	dc.SetRGB(0,0,0)
	dc.DrawString("Clan", 10, 230)
	dc.DrawStringAnchored("Imortais", 310, 230, 1, 0)
	

	dc.SetRGB(1, 1, 1)
	err = dc.LoadFontFace("./resources/Raleway-Light.ttf", 22)
	dc.DrawStringAnchored(name, 320/2, 160, 0.5, 0.5)

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
