package handler

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

type ButtonHandler struct {
	callback func(*disgord.InteractionCreate)
	sync.Mutex
}

var ButtonsHandlers = map[disgord.Snowflake]*ButtonHandler{}
var ButtonLock = sync.RWMutex{}

func RegisterBHandler(msg *disgord.Message, callback func(*disgord.InteractionCreate), timeout int) {
	ButtonsHandlers[msg.ID] = &ButtonHandler{
		callback: callback,
	}
	if timeout != 0 {
		time.Sleep(time.Duration(timeout) * time.Second)
		DeleteBHandler(msg)
	}
}

func DeleteBHandler(msg *disgord.Message) {
	ButtonLock.Lock()
	delete(ButtonsHandlers, msg.ID)
	ButtonLock.Unlock()
}

func HandleButton(interaction *disgord.InteractionCreate) {
	ButtonLock.RLock()
	if btn, found := ButtonsHandlers[interaction.Message.ID]; found {
		ButtonLock.RUnlock()
		btn.Lock()
		defer btn.Unlock()
		btn.callback(interaction)
		return
	}
	ButtonLock.RUnlock()
}

func Interaction(session disgord.Session, evt *disgord.InteractionCreate) {
	if evt.Type == disgord.InteractionMessageComponent && evt.Member != nil {
		go HandleButton(evt)
	}
}
