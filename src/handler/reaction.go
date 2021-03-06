package handler

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

var ReactionHandlers = map[disgord.Snowflake]func(bool, disgord.Emoji, disgord.Snowflake){}
var ReactionLock = sync.RWMutex{}

func RegisterHandler(msg *disgord.Message, callback func(bool, disgord.Emoji, disgord.Snowflake), timeout int) {
	ReactionHandlers[msg.ID] = callback
	if timeout != 0 {
		time.Sleep(time.Duration(timeout) * time.Second)
		DeleteHandler(msg)
	}
}

func DeleteHandler(msg *disgord.Message) {
	ReactionLock.Lock()
	delete(ReactionHandlers, msg.ID)
	ReactionLock.Unlock()
}

func sendReaction(removed bool, id disgord.Snowflake, emoji disgord.Emoji, user disgord.Snowflake) {
	ReactionLock.RLock()
	if cb, found := ReactionHandlers[id]; found {
		ReactionLock.RUnlock()
		cb(removed, emoji, user)
		return
	}
	ReactionLock.RUnlock()
}

func OnReactionAdd(session disgord.Session, evt *disgord.MessageReactionAdd) {
	go sendReaction(false, evt.MessageID, *evt.PartialEmoji, evt.UserID)
}

func OnReactionRemove(session disgord.Session, evt *disgord.MessageReactionRemove) {
	go sendReaction(true, evt.MessageID, *evt.PartialEmoji, evt.UserID)
}
