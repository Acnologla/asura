package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"bytes"
	"context"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"strings"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:  "invert",
		Cache: 15,

		Description: translation.T("InvertHelp", "pt"),
		Run:         runInvert,
		Cooldown:    60,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Description: translation.T("InvertOpt1", "pt"),
			Type:        disgord.OptionTypeUser,
		}, &disgord.ApplicationCommandOption{
			Name:        "url",
			Description: translation.T("InvertOpt2", "pt"),
			Type:        disgord.OptionTypeString,
		}),
	})
}

func runInvert(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	url := utils.GetUrl(itc)

	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")

	avatar, err := utils.DownloadImage(replacer.Replace(url))
	if err != nil {
		return nil
	}

	avatar = resize.Resize(256, 256, avatar, resize.Lanczos3)
	cimg, ok := avatar.(draw.Image)
	if !ok {
		return nil
	}
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			r, g, b, a := cimg.At(i, j).RGBA()
			cimg.Set(i, j, color.RGBA{
				uint8(255 - r),
				uint8(255 - g),
				uint8(255 - b),
				uint8(a),
			})
		}
	}
	dc := gg.NewContext(256, 256)
	dc.DrawImage(cimg, 0, 0)

	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Files: []disgord.CreateMessageFile{{
				Reader:     bytes.NewReader(b.Bytes()),
				FileName:   "inverted.png",
				SpoilerTag: false},
			},
		},
	}

}
