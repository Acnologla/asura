package commands

import (
	"asura/src/handler"
	"context"
	"github.com/andersfylling/disgord"
	"asura/src/database"
	"asura/src/utils"
	"math"
	"bytes"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image/png"
	"io"
	"os"
	"strings"
	"fmt"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"galo"},
		Run:       runTest,
		Available: true,
		Cooldown:  10,
		Usage:     "j!galo",
		Help:      "Informação sobre seu galo",
	})
}

func runTest(session disgord.Session, msg *disgord.Message, args []string) {
	var person = msg.Author
	if len(msg.Mentions) != 0 {
		person = msg.Mentions[0]
	}
	galo, err := database.GetGaloDB(person.ID)
	if err == nil{
		file, err := os.Open("resources/galo.png")
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			return
		}
		url := utils.GetImageURL(msg, args,512)
		replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
		avatar, err := utils.DownloadImage(replacer.Replace(url))

		if err != nil {
			msg.Reply(context.Background(), session, "Invalid image")
			return
		}

		avatar = resize.Resize(127,127, avatar, resize.Lanczos3)

		dc := gg.NewContext(560, 365)
		dc.DrawImage(img, 0, 0)
		dc.DrawImage(avatar, 41, 108)
		if err := dc.LoadFontFace("./resources/Raleway-Light.ttf", 24); err != nil {
			panic(err)
		}
		dc.SetRGB(1,1,1)
		dc.DrawStringAnchored("Galo de " + person.Username, 195, 167, 0, 0.5)

		level := float64(math.Floor(math.Sqrt(float64(galo.Xp)/30)))
		baseExp := float64(math.Floor(math.Pow(float64(level) ,2)*30))
		nextExp := float64(math.Floor(math.Pow(float64(level)+1,2)*30))
		nextRetExp := (float64(galo.Xp) - baseExp) / (nextExp - baseExp)

		dc.SetRGB(float64(255)/255,float64(135)/255,float64(135)/255)
		dc.DrawRectangle(204, 204, 320, 26)
		dc.Fill()
		dc.SetRGB(float64(207)/255,float64(65)/255,float64(65)/255)
		dc.DrawRectangle(207, 207, nextRetExp*(314), 20)
		dc.Fill()

		if err := dc.LoadFontFace("./resources/Raleway-Light.ttf", 14); err != nil {
			panic(err)
		}

		dc.SetRGB(1,1,1)
		expText := fmt.Sprintf("%d/%d",(galo.Xp - int(baseExp)), int((nextExp - baseExp)))
		dc.DrawStringAnchored(expText, 204 + 334/2, 216, 0.5, 0.5)

		if err := dc.LoadFontFace("./resources/Raleway-Light.ttf", 20); err != nil {
			panic(err)
		}

		total := galo.Wins+galo.Loss
		var winrate float64
		if total == 0 {
			winrate = 0 
		} else {
			winrate = float64(galo.Wins)/float64(total)
		}

		dc.SetRGB(0.3, 0.3, 0.3)
		dc.DrawStringAnchored(fmt.Sprintf("WINS: %d",galo.Wins), 520, 260, 1, 0)
		dc.DrawStringAnchored(fmt.Sprintf("LOSS: %d",galo.Loss), 520, 283, 1, 0)
		dc.DrawStringAnchored(fmt.Sprintf("WINRATE: %.2f",winrate), 520, 306, 1, 0)
		dc.DrawStringAnchored(fmt.Sprintf("FIGHTS: %d",total), 520, 329, 1, 0)

		if err := dc.LoadFontFace("./resources/Raleway-Black.ttf", 34); err != nil {
			panic(err)
		}

		dc.DrawStringAnchored("LVL", 190+ 130/2, 285, 0.5, 0)
		dc.DrawStringAnchored(fmt.Sprintf("%d",int(level + 1)), 190+ 130/2, 315, 0.5, 0)

		var b bytes.Buffer
		pw := io.Writer(&b)
		png.Encode(pw, dc.Image())
		
		msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
			Files: []disgord.CreateMessageFileParams{
				{bytes.NewReader(b.Bytes()), "galo.jpg", false},
			},
		})
	}
	
}
