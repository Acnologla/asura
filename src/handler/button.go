package handler

import (
	"fmt"
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

func RegisterBHandler(id disgord.Snowflake, callback func(*disgord.InteractionCreate), timeout int) {
	ButtonsHandlers[id] = &ButtonHandler{
		callback: callback,
	}
	if timeout != 0 {
		time.Sleep(time.Duration(timeout) * time.Second)
		DeleteBHandler(id)
	}
}

func DeleteBHandler(id disgord.Snowflake) {
	ButtonLock.Lock()
	delete(ButtonsHandlers, id)
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
	fmt.Println("receib")
	if evt.Type == disgord.InteractionMessageComponent && evt.Member != nil {
		go HandleButton(evt)
	} else if evt.Type == disgord.InteractionApplicationCommand {
		fmt.Println("x")
		handleSlash(evt)
	}
}
