package commands

import (
	"asura/src/handler"
	"asura/src/utils"
	"asura/src/utils/rinha"
	"asura/src/database"
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
		Aliases:   []string{"perfil", "profile"},
		Run:       runProfile,
		Available: true,
		Cooldown:  5,
		Usage:     "j!profile @user",
		Help:      "Veja o perfil",
		Category: 1,
	})
}

func GetGalosImages(galo *rinha.Galo) []string{
	images := []string{}
	for _, userGalo := range galo.Galos{
		images = append(images, rinha.Sprites[0][userGalo.Type -1])
	}
	return images
}


func GetClanPos(clanName string) int{
	q := database.Database.NewRef("clan").OrderByChild("xp")
	result, _ := q.GetOrdered(context.Background())
	for i, clanJSON := range result{
		if clanName ==  clanJSON.Key(){
			return i + 1
		}
	} 
	return 0
}

func runProfile(session disgord.Session, msg *disgord.Message, args []string) {
	// Loads the template of profile.png in memory
	file, err := os.Open("resources/galo/profile.png")
	defer file.Close()

	// Decodes it into a image.Image
	img, err := png.Decode(file)
	if err != nil {
		return
	}

	// Download user image
	user := utils.GetUser(msg,args,session)
	galo, _ := rinha.GetGaloDB(user.ID)
	url, _ := user.AvatarURL(128, false)
	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")
	avatar, err := utils.DownloadImage(replacer.Replace(url))

	if err != nil {
		msg.Reply(context.Background(), session, "Invalid image")
		return
	}
	// Resize the images
	avatar = resize.Resize(99, 99, avatar, resize.Lanczos3)
	name := "Galo do " + user.Username
	if galo.Name != "" {
		name = galo.Name
	}
	println(name)
	// Here we draw the image to the "editor"
	dc := gg.NewContext(419, 306)
	dc.LoadFontFace("resources/Raleway-Black.ttf", 1)
	dc.DrawImage(img, 0, 0)

	dc.DrawCircle(72, 142, 46)
	dc.Clip()
	// clipar avatar antes
	dc.DrawImage(avatar, 20, 89)
	dc.ResetClip()

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
