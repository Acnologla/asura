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

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "pragra",
		Description: translation.T("PragaHelp", "pt"),
		Run:         runPraga,
		Cooldown:    10,
		Cache:       60,
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

func runPraga(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	url := utils.GetUrl(itc)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar, err := utils.DownloadImage(replacer.Replace(url))
	if err != nil {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "Invalid image",
			},
		}
	}
	file, _ := os.Open("resources/praga.jpg")
	defer file.Close()

	// Decodes it into a image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil
	}

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

	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Files: []disgord.CreateMessageFile{{
				Reader:     bytes.NewReader(b.Bytes()),
				FileName:   "praga.jpg",
				SpoilerTag: false},
			},
		},
	}

}
