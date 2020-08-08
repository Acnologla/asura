package utils

import (
	"image/png"
	"github.com/andersfylling/disgord"
	"image"
	"fmt"
	"net/http"
	"asura/src/handler"
	"strings"
	"context"
	"strconv"
)

func checkImage(url string) bool{
    req, _ := http.NewRequest("HEAD", url, nil)
    var client http.Client
    resp, _ := client.Do(req)
    mb, _ := strconv.Atoi(resp.Header["Content-Length"][0])
	if mb > 10 * 1024 * 1024 {
		return false
	}
	if !strings.HasPrefix(resp.Header["Content-Type"][0],"image") {
		return false
	}
	return true
}

func DownloadImage(url string) (image.Image, error) { 
	response, err := http.Get(url)
    if err != nil {
		fmt.Println("Went wrong!")
	}
	defer response.Body.Close()
	img, err := png.Decode(response.Body)
	if err !=nil{
		return nil, err
	}
	return img, nil
}


func GetImageURL(msg *disgord.Message, args []string) (avatar string){
	if len(msg.Mentions) >0 { // oargs tem q ter pra gente pegar o id ou a url da image
		 avatar, _ = msg.Mentions[0].AvatarURL(512,false)
		 return
	}else if len(args) > 0{
		converted, _ := strconv.Atoi(args[0])
		user,err := handler.Client.GetUser(context.Background(),disgord.NewSnowflake(uint64(converted)))
		if err == nil{
			avatar,_ = user.AvatarURL(512,false)
		}
	}else if avatar=strings.Join(args,"");checkImage(avatar){
	}else{
		avatar, _ = msg.Author.AvatarURL(512,false)
	}
	return
}