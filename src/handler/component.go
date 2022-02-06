package handler

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

type ComponentHandler struct {
	callback   func(*disgord.InteractionCreate)
	deleteChan chan bool
	sync.Mutex
}

var ComponentHandlers = map[disgord.Snowflake]*ComponentHandler{}
var ComponentLock = sync.RWMutex{}

func RegisterHandler(itc *disgord.InteractionCreate, callback func(*disgord.InteractionCreate), timeout int) {
	ComponentHandlers[itc.ID] = &ComponentHandler{
		callback:   callback,
		deleteChan: make(chan bool),
	}
	timeChannel := time.After(time.Duration(timeout) * time.Second)
	if timeout == 0 {
		<-ComponentHandlers[itc.ID].deleteChan
	} else {
		select {
		case <-ComponentHandlers[itc.ID].deleteChan:
		case <-timeChannel:
			ComponentLock.Lock()
			delete(ComponentHandlers, itc.ID)
			ComponentLock.Unlock()
		}
	}

}

func DeleteHandler(itc *disgord.InteractionCreate) {
	ComponentLock.Lock()
	button, ok := ComponentHandlers[itc.ID]
	if ok {
		button.deleteChan <- true
	}
	delete(ComponentHandlers, itc.ID)
	ComponentLock.Unlock()
}

func HandleComponent(interaction *disgord.InteractionCreate) {
	ComponentLock.RLock()
	if btn, found := ComponentHandlers[interaction.Message.Interaction.ID]; found {
		ComponentLock.RUnlock()
		btn.Lock()
		defer btn.Unlock()
		btn.callback(interaction)
		return
	}
	ComponentLock.RUnlock()
}

func ComponentInteraction(session disgord.Session, evt *disgord.InteractionCreate) {
	if evt.Member != nil {
		go HandleComponent(evt)
	}
}
