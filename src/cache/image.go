package cache

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"time"
)

func GetImageFromCache(sufix string, name string) *image.Image {
	bytesImg := Client.Get(context.Background(), fmt.Sprintf("/images/%s/%s", sufix, name))
	result, _ := bytesImg.Bytes()
	if len(result) == 0 {
		return nil
	}
	img, err := png.Decode(bytes.NewReader(result))
	if err != nil {
		return nil
	}
	return &img
}

func CacheProfileImage(img *image.Image, name string) {
	if img == nil {
		return
	}
	var buf bytes.Buffer
	png.Encode(&buf, *img)
	if buf.Len() > 0 {
		Client.Set(context.Background(), "/images/profile/"+name, buf.Bytes(), time.Minute*120)
	}
}
