package utils

import (
	"asura/src/handler"
	"context"
	"github.com/andersfylling/disgord"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"
	"strings"
)

// This function checks if the url has a valid image by checking the head of the url
// The head contains the Content-Length (that we limit the size to avoid downloading absurd res images)
// And the Content-type that will tell us if the url has an image
func checkImage(url string) bool {
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
	if mb > 10*1024*1024 {
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

// TODO: Remove this absurd quantity of ifs and make the function get 
// Persons without avatar
// This function is used to get a url for an iamge that will be used
// To a lot of functions. 
func GetImageURL(msg *disgord.Message, args []string) string {
	if len(msg.Mentions) > 0 {
		avatar, _ := msg.Mentions[0].AvatarURL(512, false)
		return avatar
	}
	if len(args) > 0 {
		converted, _ := strconv.Atoi(args[0])
		user, err := handler.Client.GetUser(context.Background(), disgord.NewSnowflake(uint64(converted)))
		if err == nil {
			avatar, _ := user.AvatarURL(512, false)
			return avatar
		}
	}
	if checkImage(strings.Join(args, "")) {
		return strings.Join(args, "")
	}
	avatar, _ := msg.Author.AvatarURL(512, false)
	return avatar
}