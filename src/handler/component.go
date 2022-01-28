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

func RegisterBHandler(itc *disgord.InteractionCreate, callback func(*disgord.InteractionCreate), timeout int) {
	ButtonsHandlers[itc.ID] = &ButtonHandler{
		callback:   callback,
		deleteChan: make(chan bool),
	}
	timeChannel := time.After(time.Duration(timeout) * time.Second)
	if timeout == 0 {
		<-ButtonsHandlers[itc.ID].deleteChan
	} else {
		select {
		case <-ButtonsHandlers[itc.ID].deleteChan:
		case <-timeChannel:
			ButtonLock.Lock()
			delete(ButtonsHandlers, itc.ID)
			ButtonLock.Unlock()
		}
	}

}

func DeleteBHandler(itc *disgord.InteractionCreate) {
	ButtonLock.Lock()
	button, ok := ButtonsHandlers[itc.ID]
	if ok {
		button.deleteChan <- true
	}
	delete(ButtonsHandlers, itc.ID)
	ButtonLock.Unlock()
}

func HandleButton(interaction *disgord.InteractionCreate) {
	ButtonLock.RLock()
	if btn, found := ButtonsHandlers[interaction.Message.Interaction.ID]; found {
		ButtonLock.RUnlock()
		btn.Lock()
		defer btn.Unlock()
		btn.callback(interaction)
		return
	}
	ButtonLock.RUnlock()
}

func ComponentInteraction(session disgord.Session, evt *disgord.InteractionCreate) {
	if evt.Member != nil {
		go HandleButton(evt)
	}
}
