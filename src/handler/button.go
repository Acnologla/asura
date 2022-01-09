package handler

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

type ButtonHandler struct {
	callback   func(*disgord.InteractionCreate)
	deleteChan chan bool
	sync.Mutex
}

var ButtonsHandlers = map[disgord.Snowflake]*ButtonHandler{}
var ButtonLock = sync.RWMutex{}

func RegisterBHandler(msg *disgord.Message, callback func(*disgord.InteractionCreate), timeout int) {
	ButtonsHandlers[msg.ID] = &ButtonHandler{
		callback:   callback,
		deleteChan: make(chan bool),
	}
	timeChannel := time.After(time.Duration(timeout) * time.Second)
	if timeout == 0 {
		<-ButtonsHandlers[msg.ID].deleteChan
	} else {
		select {
		case <-ButtonsHandlers[msg.ID].deleteChan:
		case <-timeChannel:
			ButtonLock.Lock()
			delete(ButtonsHandlers, msg.ID)
			ButtonLock.Unlock()
		}
	}
}

func DeleteBHandler(msg *disgord.Message) {
	ButtonLock.Lock()
	button, ok := ButtonsHandlers[msg.ID]
	if ok {
		button.deleteChan <- true
	}
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
