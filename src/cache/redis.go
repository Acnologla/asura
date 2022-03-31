package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

type _createInteractionResponseData struct {
	Content         string                      `json:"content"`
	Tts             bool                        `json:"tts,omitempty"`
	Embeds          []*disgord.Embed            `json:"embeds,omitempty"`
	Components      []*disgord.MessageComponent `json:"components"`
	Attachments     []*disgord.Attachment       `json:"attachments"`
	AllowedMentions *disgord.AllowedMentions    `json:"allowed_mentions,omitempty"`

	Files []_createMessageFile `json:"files"`

	SpoilerTagContent        bool `json:"-"`
	SpoilerTagAllAttachments bool `json:"-"`
}

type _createMessageFile struct {
	Bytes      []byte `json:"reader"`
	FileName   string `json:"filename"`
	SpoilerTag bool   `json:"-"`
}

type redisConfig struct {
	URL      string `json:"REDIS_URL"`
	Password string `json:"REDIS_PASSWORD"`
}

func loadRedisEnv() redisConfig {
	var config = os.Getenv("REDIS_CONFIG")
	var redisConfig redisConfig
	json.Unmarshal([]byte(config), &redisConfig)
	return redisConfig
}

func Init() {
	config := loadRedisEnv()
	Client = redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: config.Password,
	})
	cmd := Client.Ping(context.Background())
	err := cmd.Err()
	if err != nil {
		log.Fatal(err)
	}
}
func parseOpts(opts []*disgord.ApplicationCommandDataOption) string {
	str := ""
	for _, opt := range opts {
		if opt.Type == disgord.OptionTypeSubCommand || opt.Type == disgord.OptionTypeSubCommandGroup {
			return opt.Name + "-" + parseOpts(opt.Options)
		} else {
			str += fmt.Sprintf("%s:%v", opt.Name, opt.Value)
		}
	}
	return str
}
func ParseCmdOpts(itc *disgord.InteractionCreate) string {
	return itc.Data.Name + "-" + parseOpts(itc.Data.Options)
}

func GetCachedCommand(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponseData {
	parsed := ParseCmdOpts(itc)
	var res _createInteractionResponseData
	b, err := Client.Get(ctx, "/cmds/"+parsed).Bytes()
	if err != nil {
		return nil
	}
	json.Unmarshal(b, &res)
	f := []disgord.CreateMessageFile{}
	for _, file := range res.Files {
		f = append(f, disgord.CreateMessageFile{
			Reader:     bytes.NewReader(file.Bytes),
			FileName:   file.FileName,
			SpoilerTag: file.SpoilerTag,
		})
	}
	return &disgord.CreateInteractionResponseData{
		Content:         res.Content,
		Tts:             res.Tts,
		Embeds:          res.Embeds,
		Components:      res.Components,
		Attachments:     res.Attachments,
		AllowedMentions: res.AllowedMentions,
		Files:           f,
	}

}

type x = []_createMessageFile

func CacheCommand(ctx context.Context, itc *disgord.InteractionCreate, res *disgord.CreateInteractionResponse, minutes int) {
	parsed := ParseCmdOpts(itc)
	f := []_createMessageFile{}
	for i, file := range res.Data.Files {
		buf := new(bytes.Buffer)
		buf.ReadFrom(file.Reader)
		b := buf.Bytes()
		f = append(f, _createMessageFile{
			Bytes:      b,
			FileName:   file.FileName,
			SpoilerTag: file.SpoilerTag,
		})
		res.Data.Files[i].Reader = bytes.NewReader(b)
	}
	data := _createInteractionResponseData{
		Content:                  res.Data.Content,
		Tts:                      res.Data.Tts,
		Embeds:                   res.Data.Embeds,
		Components:               res.Data.Components,
		Attachments:              res.Data.Attachments,
		AllowedMentions:          res.Data.AllowedMentions,
		Files:                    f,
		SpoilerTagContent:        res.Data.SpoilerTagContent,
		SpoilerTagAllAttachments: res.Data.SpoilerTagAllAttachments,
	}
	r, _ := json.Marshal(data)
	Client.Set(ctx, "/cmds/"+parsed, r, time.Minute*time.Duration(minutes))
}
