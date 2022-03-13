package utils

import (
	"asura/src/handler"
	"image"
	"net/http"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

func DownloadAvatar(id disgord.Snowflake, size int, gif bool) (image.Image, error) {
	user, _ := handler.Client.User(id).WithFlags(disgord.IgnoreCache).Get()
	avatar, _ := user.AvatarURL(size, gif)
	avatar = strings.Replace(avatar, ".webp", ".png", 1)
	return DownloadImage(avatar)
}

func CheckImage(url string) bool {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false
	}
	var client http.Client
	resp, err := client.Do(req)
	if err != nil || len(resp.Header["Content-Length"]) == 0 || len(resp.Header["Content-Type"]) == 0 {
		return false
	}
	defer resp.Body.Close()
	mb, _ := strconv.Atoi(resp.Header["Content-Length"][0])
	if mb > 3*1024*1024 {
		return false
	}
	if !strings.HasPrefix(resp.Header["Content-Type"][0], "image") {
		return false
	}
	return true
}

// This function has simply download an image and return it.
func DownloadImage(url string) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func GetUrl(itc *disgord.InteractionCreate) string {
	url, _ := itc.Member.User.AvatarURL(1024, false)

	replacer := strings.NewReplacer(".gif", ".png", ".webp", ".png")

	if len(itc.Data.Options) > 0 {
		if itc.Data.Options[0].Name == "user" {
			user := GetUser(itc, 0)
			url, _ = user.AvatarURL(1024, false)
		} else {
			_url := itc.Data.Options[0].Value.(string)

			if CheckImage(replacer.Replace(url)) {
				url = _url
			}
		}
	}
	return url
}
