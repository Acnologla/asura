package handler

import (
	"github.com/andersfylling/disgord"
	"sync"
	"time"
)

type Message struct{
	disgord.Message
	sync.RWMutex
	Date time.Time
}

type Cache struct{
	*disgord.CacheLFUImmutable
	Messages map[disgord.Snowflake]*Message
}

func (cache *Cache) MessageCreate(data []byte) (*disgord.MessageCreate, error){
	msg, err := cache.CacheNop.MessageCreate(data)
	message := msg.Message
	if err != nil {
		return msg, err
	}
	if message.Author != nil && message.Type == disgord.MessageTypeDefault {
		if !message.Author.Bot{
			cache.Messages[message.ID] = &Message{
				Message: *message,
				Date: time.Now(),
			}
		}
	}
	return msg, nil
}

func (cache *Cache) MessageUpdate(data []byte) (*disgord.MessageUpdate, error){
	msg, err := cache.CacheNop.MessageUpdate(data)
	if err != nil {
		return msg, err
	}
	oldMsg, ok := cache.Messages[msg.Message.ID]
	if ok{
		oldMsg.Lock()
		oldMsgCopy := oldMsg.Message
		go OnMessageUpdate(&oldMsgCopy, msg.Message)
		oldMsg.Message = *msg.Message
		oldMsg.Date = time.Now()
		oldMsg.Unlock()
	}
	return msg, nil
}


func init(){
	//clear 30 minutes old messages
	go func(){
		time.Sleep(time.Minute * 30)
		for{
			cache := Client.Cache().(*Cache)
			for _, msg := range cache.Messages{
				msg.Lock()
				if time.Since(msg.Date).Minutes() >= 30{
					delete(cache.Messages, msg.ID)
				}
				msg.Unlock()
			}
			time.Sleep(time.Minute)
		}
	}()
}