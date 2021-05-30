package handler

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

type Message struct {
	disgord.Message
	sync.RWMutex
	Date time.Time
}

type Cache struct {
	sync.RWMutex
	*disgord.BasicCache
	Messages map[disgord.Snowflake]*Message
}

func (cache *Cache) MessageCreate(data []byte) (*disgord.MessageCreate, error) {
	msg, err := cache.CacheNop.MessageCreate(data)
	if err != nil {
		return nil, err
	}
	message := msg.Message
	if message.Author != nil && message.Type == disgord.MessageTypeDefault {
		if !message.Author.Bot {
			entry := &Message{
				Message: *(disgord.DeepCopy(message).(*disgord.Message)),
				Date:    time.Now(),
			}
			cache.Lock()
			cache.Messages[message.ID] = entry
			cache.Unlock()
		}
	}
	return msg, nil
}

func (cache *Cache) MessageUpdate(data []byte) (*disgord.MessageUpdate, error) {
	msg, err := cache.CacheNop.MessageUpdate(data)
	if err != nil {
		return msg, err
	}
	cache.RLock()
	oldMsg, ok := cache.Messages[msg.Message.ID]
	cache.RUnlock()
	if ok {
		oldMsg.Lock()
		oldMsgCopy := oldMsg.Message
		go OnMessageUpdate(&oldMsgCopy, msg.Message)
		oldMsg.Message = *(disgord.DeepCopy(msg.Message).(*disgord.Message))
		oldMsg.Date = time.Now()
		oldMsg.Unlock()
	}
	return msg, nil
}

func init() {
	//clear N minutes old messages
	const N = 30
	go func() {
		time.Sleep(time.Minute * 30)
		ticker := time.NewTicker(time.Duration(N) * time.Minute)
		cache := Client.Cache().(*Cache)
		for {
			now := time.Now()
			cache.Lock()
			for id, msg := range cache.Messages {
				if now.Sub(msg.Date).Minutes() >= N {
					delete(cache.Messages, id)
				}
			}
			cache.Unlock()
			<-ticker.C
		}
	}()
}
