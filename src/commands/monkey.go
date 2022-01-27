package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"bytes"
	"context"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"macaco", "monkey"},
		Run:       runMonkey,
		Available: true,
		Cooldown:  5,
		Usage:     "j!monkey @user",
		Help:      "Transforme algum em um macaco",
	})
}

func runMonkey(session disgord.Session, msg *disgord.Message, args []string) {
	// Loads the template of nunca.png in memory
	file, _ := os.Open("resources/monkey.png")
	defer file.Close()

	// Decodes it into a image.Image
	img, err := png.Decode(file)
	if err != nil {
		return
	}

	// Download user image
	url := utils.GetImageURL(msg, args, 512, session)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar, err := utils.DownloadImage(replacer.Replace(url))

	if err != nil {
		msg.Reply(context.Background(), session, "Invalid image")
		return
	}
	// Here we draw the image to the "editor"
	avatar = resize.Resize(491, 406, avatar, resize.Lanczos3)
	dc := gg.NewContext(491, 406)
	dc.DrawImage(avatar, 0, 0)
	dc.DrawImage(img, 0, 0)
	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())

	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{{
			Reader:     bytes.NewReader(b.Bytes()),
			FileName:   "monkey.png",
			SpoilerTag: false},
		},
	})
}
