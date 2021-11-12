package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"bytes"
	"context"
	"image/jpeg"
	"io"
	"os"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"praga", "morrepraga", "morre"},
		Run:       runPraga,
		Available: true,
		Cooldown:  5,
		Usage:     "j!praga @user",
		Help:      "Morra praga",
	})
}

func runPraga(session disgord.Session, msg *disgord.Message, args []string) {
	// Loads the template of nunca.png in memory
	file, _ := os.Open("resources/praga.jpg")
	defer file.Close()

	// Decodes it into a image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return
	}

	// Download user image
	url := utils.GetImageURL(msg, args, 256, session)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar, err := utils.DownloadImage(replacer.Replace(url))

	if err != nil {
		msg.Reply(context.Background(), session, "Invalid image")
		return
	}
	// Resize the images
	avatar = resize.Resize(160, 160, avatar, resize.Lanczos3)

	// Here we draw the image to the "editor"
	dc := gg.NewContext(540, 325)
	dc.DrawImage(img, 0, 0)
	dc.DrawImage(avatar, 55, 90)

	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	jpeg.Encode(pw, dc.Image(), &jpeg.Options{
		Quality: 80,
	})

	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{{
			Reader:     bytes.NewReader(b.Bytes()),
			FileName:   "lixeira.jpg",
			SpoilerTag: false},
		},
	})
}
