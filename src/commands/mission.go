package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"bytes"
	"context"
	"fmt"
	"image/png"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "mission",
		Description: translation.T("MissionHelp", "pt"),
		Run:         runMission,
		Cooldown:    10,
		Category:    handler.Rinha,
	})
}

func runMission(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := itc.Member.User
	galo := database.User.GetUser(ctx, user.ID, "Missions")
	galo.Missions = getMissions(ctx, &galo)
	texts := rinha.MissionsToString(user.ID, &galo)
	galoImage, err := utils.DownloadImage(rinha.Sprites[0][0])
	if err != nil {
		return nil
	}
	re := regexp.MustCompile(`\(([^)]+)\)`)

	width, height := 580, 480
	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)

	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()
	radius := 50
	imageSize := radius * 2
	for i, mission := range galo.Missions {
		missionImage := galoImage
		if mission.Adv != 0 {
			missionImage, _ = utils.DownloadImage(rinha.Sprites[0][mission.Adv-1])
		}
		text := texts[i]
		splited := strings.Split(text, "\n")
		missionText := splited[0]
		money := splited[1]
		xp := splited[2]
		img := resize.Resize(uint(imageSize), uint(imageSize), missionImage, resize.Lanczos3)
		imageVerticalMargin := 20 + i*imageSize + 15*i
		dc.SetRGB(0, 0, 0)
		//radius := float64(imageSize) / (math.Sqrt(math.Pi))
		dc.DrawCircle(float64(10+radius), float64(imageVerticalMargin+radius), float64(radius))
		dc.Clip()
		dc.Fill()
		dc.DrawImage(img, 10, imageVerticalMargin)
		dc.ResetClip()
		dc.SetRGB(0, 0, 0)
		textMargin := float64(imageSize + 30)
		verticalMargin := float64(imageVerticalMargin) + 24
		dc.LoadFontFace("./resources/Raleway-Bold.ttf", 18)
		dc.DrawString(missionText, textMargin, verticalMargin)
		dc.LoadFontFace("./resources/Raleway-Light.ttf", 19)
		dc.DrawString(money, textMargin, verticalMargin+20)
		dc.DrawString(xp, textMargin, verticalMargin+40)
		//progress bar
		matches := re.FindAllStringSubmatch(text, -1)
		splitedMatches := strings.Split(matches[0][1], "/")
		progress, _ := strconv.Atoi(splitedMatches[0])
		total, _ := strconv.Atoi(splitedMatches[1])
		percentange := float64(progress) / float64(total)
		barWidth := float64(width) - 50 - textMargin
		barHeight := 15
		dc.DrawRoundedRectangle(textMargin, verticalMargin+50, barWidth, float64(barHeight), 4)
		dc.SetRGB(0.8, 0.8, 0.8)
		dc.Fill()
		if percentange > 0 {
			dc.SetRGB(0.2, 0.8, 0.2)
			dc.DrawRoundedRectangle(textMargin, verticalMargin+50, barWidth*percentange, float64(barHeight), 4)
			dc.Fill()
		}
	}

	var b bytes.Buffer
	pw := io.Writer(&b)
	png.Encode(pw, dc.Image())
	embed := &disgord.Embed{
		Color: 65535,
		Title: fmt.Sprintf("Missoes (%d/4)", len(galo.Missions)),
		Image: &disgord.EmbedImage{
			URL: "attachment://mission.png",
		},
	}
	if len(galo.Missions) != 4 {
		need := uint64(time.Now().Unix()) - galo.LastMission
		embed.Footer = &disgord.EmbedFooter{
			Text: translation.T("MissionTime", translation.GetLocale(itc), map[string]interface{}{
				"hours":   23 - (need / 60 / 60),
				"minutes": 59 - (need / 60 % 60),
			}),
		}
	}
	components := []*disgord.MessageComponent{}
	if (time.Now().Unix()-int64(galo.TradeMission))/60/60/24 >= 1 {
		for i := range galo.Missions {
			components = append(components, &disgord.MessageComponent{
				Type:     disgord.MessageComponentButton,
				Style:    disgord.Primary,
				Label:    translation.T("ChangeMission", translation.GetLocale(itc), i),
				CustomID: strconv.Itoa(i),
			})
		}
	}
	component := []*disgord.MessageComponent{{
		Type:       disgord.MessageComponentActionRow,
		Components: components,
	}}
	params := &disgord.CreateInteractionResponseData{
		Files: []disgord.CreateMessageFile{{
			Reader:     bytes.NewReader(b.Bytes()),
			FileName:   "mission.png",
			SpoilerTag: false},
		},
	}
	params.Embeds = []*disgord.Embed{embed}
	if len(components) > 0 {
		params.Components = component
	}
	handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: params,
	})
	if len(components) > 0 {
		handler.RegisterHandler(itc.ID, func(event *disgord.InteractionCreate) {
			if event.Member.User.ID == user.ID {
				i, _ := strconv.Atoi(event.Data.CustomID)
				database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
					if 0 > i || i >= len(u.Missions) || (time.Now().Unix()-int64(u.TradeMission))/60/60/24 < 3 {
						return u
					}
					newMission := rinha.CreateMission()
					newMission.ID = u.Missions[i].ID
					newMission.UserID = user.ID
					database.User.UpdateMissions(ctx, user.ID, &newMission, false)
					u.TradeMission = uint64(time.Now().Unix())
					return u
				}, "Missions")
				handler.DeleteHandler(itc.ID)
				handler.Client.SendInteractionResponse(ctx, event, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: translation.T("TradeMission", translation.GetLocale(event)),
					},
				})
			}
		}, 120)
	}
	return nil
}
