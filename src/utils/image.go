package utils

import (
	"image"
	"net/http"
	"strconv"
	"strings"
)

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
