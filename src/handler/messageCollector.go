package handler

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

type MessageCollector struct {
	Filter      func(msg *disgord.Message) bool
	MessageChan chan *disgord.Message
	Timestamp   time.Time
}

var MessageHandlers = map[disgord.Snowflake][]MessageCollector{}
var MessageLock = sync.RWMutex{}

func CreateMessageCollector(id disgord.Snowflake, filter func(msg *disgord.Message) bool) *disgord.Message {
	MessageLock.Lock()
	_, ok := MessageHandlers[id]
	timestamp := time.Now()
	msgChan := make(chan *disgord.Message)
	messageCollector := MessageCollector{
		Filter:      filter,
		MessageChan: msgChan,
		Timestamp:   timestamp,
	}
	if ok {
		MessageHandlers[id] = append(MessageHandlers[id], messageCollector)
	} else {
		MessageHandlers[id] = []MessageCollector{messageCollector}
	}
	MessageLock.Unlock()
	var msgReturned *disgord.Message
	select {
	case <-time.After(time.Minute):
	case msgReturned = <-msgChan:
	}
	MessageLock.Lock()
	DeleteCollector(id, timestamp)
	MessageLock.Unlock()
	return msgReturned
}

func sendMsg(msg *disgord.Message) {
	MessageLock.RLock()
	collectors, ok := MessageHandlers[msg.ChannelID]
	MessageLock.RUnlock()
	if ok {
		for _, collector := range collectors {
			if collector.Filter(msg) {
				collector.MessageChan <- msg
			}
		}
	}
}

func DeleteCollector(id disgord.Snowflake, timestamp time.Time) {
	i := -1
	for j, value := range MessageHandlers[id] {
		if value.Timestamp.UnixNano() == timestamp.UnixNano() {
			i = j
			break
		}
	}
	if i == -1 {
		return
	}
	MessageHandlers[id][i] = MessageHandlers[id][len(MessageHandlers[id])-1]
	MessageHandlers[id] = MessageHandlers[id][:len(MessageHandlers[id])-1]
	if len(MessageHandlers[id]) == 0 {
		delete(MessageHandlers, id)
	}
}
