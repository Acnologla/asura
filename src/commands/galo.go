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
		Cooldown:  1,
		Usage:     "j!galo",
		Help:      "Informação sobre seu galo",
	})
}

func runTest(session disgord.Session, msg *disgord.Message, args []string) {
	galo, err := database.GetGaloDB(msg.Author.ID)

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

		avatar = resize.Resize(160,160, avatar, resize.Lanczos3)

		dc := gg.NewContext(555, 204)
		dc.DrawImage(img, 0, 0)
		dc.DrawImage(avatar, 20, 21)
		if err := dc.LoadFontFace("./resources/Raleway-Light.ttf", 24); err != nil {
			panic(err)
		}
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored("Galo de " + msg.Author.Username, 204 + 334/2, 40, 0.5, 0.5)

		level := float64(math.Floor(math.Sqrt(float64(galo.Xp)/30)))
		baseExp := float64(math.Floor(math.Pow(float64(level) ,2)*30))
		nextExp := float64(math.Floor(math.Pow(float64(level)+1,2)*30))
		nextRetExp := (float64(galo.Xp) - baseExp) / (nextExp - baseExp)

		dc.SetRGB(1,1,1)
		dc.DrawRectangle(204 + 17, 66, 300, 20)
		dc.Fill()
		dc.SetRGB(0.5, 0.5, 0.5)
		dc.DrawRectangle(204 + 20, 69, nextRetExp*(300-17), 14)
		dc.Fill()

		if err := dc.LoadFontFace("./resources/Raleway-Light.ttf", 14); err != nil {
			panic(err)
		}

		dc.SetRGB(0,0,0)
		expText := fmt.Sprintf("%d/%d",(galo.Xp - int(baseExp)), int((nextExp - baseExp)))
		dc.DrawStringAnchored(expText, 204 + 334/2, 76, 0.5, 0.5)

		if err := dc.LoadFontFace("./resources/Raleway-Black.ttf", 34); err != nil {
			panic(err)
		}
		dc.SetRGB(0.3, 0.3, 0.3)
		fmt.Printf("%d \n\r",204 +(5/6 * 334))
		dc.DrawStringAnchored("LVL", 204 + 334/6, 130, 0.5, 0)
		dc.DrawStringAnchored("WINS", 204 + 334/2, 130, 0.5, 0)
		dc.DrawStringAnchored("TYPE", 204 +(5*334)/6, 130, 0.5, 0)
		if err := dc.LoadFontFace("./resources/Raleway-Black.ttf", 25); err != nil {
			panic(err)
		}
		dc.DrawStringAnchored(fmt.Sprintf("%d",int(level+1)), 204 + 334/6, 154, 0.5, 0.5)
		dc.DrawStringAnchored("343", 204 + 334/2, 150, 0.5, 0.5)
		dc.DrawStringAnchored("FIRE", 204 +(5*334)/6, 154, 0.5, 0.5)

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
