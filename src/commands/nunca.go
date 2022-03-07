package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"bytes"
	"image/png"
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
		Name:        "nunca",
		Description: translation.T("NuncaHelp", "pt"),
		Run:         runNunca,
		Cooldown:    10,
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

func runNunca(itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
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
	file, _ := os.Open("resources/nunca.png")
	defer file.Close()

	// Decodes it into a image.Image
	img, err := png.Decode(file)
	if err != nil {
		return nil
	}
	img = resize.Resize(500, 500, img, resize.Lanczos3)
	avatar = resize.Resize(500, 350, avatar, resize.Lanczos3)
	// Here we draw the image to the "editor"
	dc := gg.NewContext(500, 500)
	dc.DrawImage(img, 0, 0)
	dc.DrawImage(avatar, 0, 150)

	// And here we encode it to send
	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())

	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Files: []disgord.CreateMessageFile{{
				Reader:     bytes.NewReader(b.Bytes()),
				FileName:   "lixeira.jpg",
				SpoilerTag: false},
			},
		},
	}

}
