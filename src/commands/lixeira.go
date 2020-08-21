package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"bytes"
	"context"
	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image/png"
	"io"
	"os"
	"strings"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"lixeira", "trashcan", "lixo"},
		Run:       runLixeira,
		Available: true,
		Cooldown:  5,
		Usage:     "j!lixeira @user",
		Help:      "Um lixo total",
	})
}

func runLixeira(session disgord.Session, msg *disgord.Message, args []string) {
	// Loads the template of nunca.png in memory
	file, err := os.Open("resources/lixeira.png")
	defer file.Close()

	// Decodes it into a image.Image
	img, err := png.Decode(file)
	if err != nil {
		return
	}

	// Download user image
	url := utils.GetImageURL(msg, args,128,session)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar, err := utils.DownloadImage(replacer.Replace(url))

	if err != nil {
		msg.Reply(context.Background(), session, "Invalid image")
		return
	}
	// Resize the images
	avatar = resize.Resize(38, 38, avatar, resize.Lanczos3)

	// Here we draw the image to the "editor"
	dc := gg.NewContext(324, 132)
	dc.DrawImage(img, 0, 0)
	dc.DrawImage(avatar, 258, 52)

	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())

	session.SendMsg(context.Background(), msg.ChannelID, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{
			{bytes.NewReader(b.Bytes()), "lixeira.jpg", false},
		},
	})
}
