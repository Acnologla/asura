package handler

import (
	"github.com/andersfylling/disgord"
	"sync"
	"time"
)


var ReactionHandlers =  map[disgord.Snowflake]func(bool, disgord.Emoji,disgord.Snowflake){}
var ReactionLock = &sync.RWMutex{}

func RegisterHandler(msg *disgord.Message, callback func(bool, disgord.Emoji,disgord.Snowflake), timeout int) {
	ReactionHandlers[msg.ID] = callback
	if timeout != 0 {
		go func() {
			time.Sleep(time.Duration(timeout) * time.Second)
			DeleteHandler(msg)
		}()
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
		cb(removed, emoji,user)
		return
	}
	ReactionLock.RUnlock()
}

func OnReactionAdd(session disgord.Session, evt *disgord.MessageReactionAdd) {
	sendReaction(false, evt.MessageID, *evt.PartialEmoji,evt.UserID)
}

func OnReactionRemove(session disgord.Session, evt *disgord.MessageReactionRemove) {
	sendReaction(true, evt.MessageID, *evt.PartialEmoji, evt.UserID)
}