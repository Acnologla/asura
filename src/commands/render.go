package commands

import (
	"asura/src/handler"
	"net/http"
	"time"
	"bytes"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/anaskhan96/soup"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"strings"
)

func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.EmulateViewport(1366, 768),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			*res, err = page.CaptureScreenshot().
				WithQuality(90).
				WithClip(&page.Viewport{
					X:     	0,
					Y:      0,
					Width:  1366,
					Height: 768,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"render", "viewsite", "rendersite", "view", "siteview"},
		Run:       runRender,
		Available: true,
		Cooldown:  10,
		Usage:     "j!render <site>",
		Help:      "Renderize um website",
	})
}

func runRender(session disgord.Session, msg *disgord.Message, args []string) {
	text := strings.Join(args, " ")
	if !strings.HasPrefix(text, "http") {
		text = "http://" + text
	}
	hresp, err := http.Get(text)
	if err != nil {
		msg.Reply(context.Background(), session, msg.Author.Mention()+ ", Site invalido")
		return
	}
	finalURL := hresp.Request.URL.String()
	resp, err := soup.Get(fmt.Sprintf("https://fortiguard.com/search?q=%s&engine=1",finalURL))
	if err != nil {
		msg.Reply(context.Background(), session, msg.Author.Mention()+ ", Site invalido")
		return
	}
	doc := soup.HTMLParse(resp)
	p := doc.Find("section", "class", "iprep").Find("h2").Find("a").Text()
	channel,err := handler.Client.GetChannel(context.Background(),msg.ChannelID)
	if err != nil{
		return
	}
	if p == "Pornography" && !channel.NSFW{
		msg.Reply(context.Background(), session, msg.Author.Mention()+ ", Voce não pode renderizar sites pornograficos")
		return
	}
	contex,_ := context.WithTimeout(context.Background(),15*time.Second)
	ctx, cancel := chromedp.NewContext(contex)
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(text, `html`, &buf)); err != nil {
		return
	}
	avatar, _ := msg.Author.AvatarURL(512, false)
	msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Files: []disgord.CreateMessageFileParams{
			{bytes.NewReader(buf), "render.jpg", false},
		},
		Embed: &disgord.Embed{
			Color: 65535,
			Description: fmt.Sprintf("[**%s**](%s)",strings.Join(args, " "),finalURL),
			Image: &disgord.EmbedImage{
				URL: "attachment://render.jpg",
			},
			Footer: &disgord.EmbedFooter{
				IconURL: avatar,
				Text:    msg.Author.Username,
			},
		},
	})
}