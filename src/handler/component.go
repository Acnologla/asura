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

func RegisterHandler(id disgord.Snowflake, callback func(*disgord.InteractionCreate), timeout int) {
	ComponentHandlers[id] = &ComponentHandler{
		callback:   callback,
		deleteChan: make(chan bool),
	}
	timeChannel := time.After(time.Duration(timeout) * time.Second)
	if timeout == 0 {
		<-ComponentHandlers[id].deleteChan
	} else {
		select {
		case <-ComponentHandlers[id].deleteChan:
		case <-timeChannel:
			ComponentLock.Lock()
			delete(ComponentHandlers, id)
			ComponentLock.Unlock()
		}
	}

}

func DeleteHandler(id disgord.Snowflake) {
	ComponentLock.Lock()
	button, ok := ComponentHandlers[id]
	if ok {
		button.deleteChan <- true
	}
	delete(ComponentHandlers, id)
	ComponentLock.Unlock()
}

func HandleComponent(interaction *disgord.InteractionCreate) {
	ComponentLock.RLock()
	btn, ok := ComponentHandlers[interaction.Message.ID]
	if !ok {
		if interaction.Message.Interaction != nil {
			btn, ok = ComponentHandlers[interaction.Message.Interaction.ID]
		}
	}
	if ok {
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
