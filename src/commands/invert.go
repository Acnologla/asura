package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"bytes"
	"github.com/nfnt/resize"
	"context"
	"image/draw"
	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"image/png"
	"io"
	"image/color"
	"strings"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"invert", "inverter", "reverse"},
		Run:       runInvert,
		Available: true,
		Cooldown:  5,
		Usage:     "j!invert @user",
		Help:      "Inverta o avatar de alguem",
	})
}

func runInvert(session disgord.Session, msg *disgord.Message, args []string) {
	
	// Download user image
	url := utils.GetImageURL(msg, args,256)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar, err := utils.DownloadImage(replacer.Replace(url))

	if err != nil {
		msg.Reply(context.Background(), session, "Invalid image")
		return
	}
	// Resize the images
	// Here we draw the image to the "editor"
	avatar = resize.Resize(256, 256, avatar, resize.Lanczos3)

	cimg, ok := avatar.(draw.Image)
	if !ok{
		return
	}
	for i:=0; i <  256;i++{
		for j:=0; j < 256;j++{
			r,g,b,a := cimg.At(i,j).RGBA()
			cimg.Set(i,j,color.RGBA{
				uint8(255 - r),
				uint8(255 - g),
				uint8(255 - b),
				uint8(a),
			})
		}
	}
	dc := gg.NewContext(256, 256)
	dc.DrawImage(cimg, 0,0)

	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())

	session.SendMsg(context.Background(), msg.ChannelID, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{
			{bytes.NewReader(b.Bytes()), "inverted.png", false},
		},
	})
}
