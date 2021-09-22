package test

import (
	"asura/src/handler"
	"testing"
	"time"

	"github.com/andersfylling/disgord"
)

func TestFindCommand(t *testing.T) {
	pingCommand := handler.FindCommand("ping")
	if len(pingCommand.Aliases) == 0 {
		t.Errorf("Command ping must exists")
	}
}

func TestCompareStrings(t *testing.T) {
	distance := handler.CompareStrings("acno", "acnologia")
	if distance != 5 {
		t.Errorf("Distance must be 5")
	}
	distance = handler.CompareStrings("acno", "acno")
	if distance != 0 {
		t.Errorf("Distance must be 0 ")
	}
}

func TestGetFreeWorkers(t *testing.T) {
	if handler.GetFreeWorkers() != handler.BaseWorkers {
		t.Errorf("GetFreeWorkers Must return handler.BaseWorkers")
	}
}

func TestMessageCollector(t *testing.T) {
	go func() {
		time.Sleep(time.Second)
		handler.SendMsg(&disgord.Message{
			ChannelID: 0,
			Content:   "Test",
		})
	}()
	msg := handler.CreateMessageCollector(0, func(msg *disgord.Message) bool {
		return true
	})
	if msg.Content != "Test" {
		t.Errorf("Message content must be 'Test'")
	}
}

func TestButton(t *testing.T) {
	msg := &disgord.Message{
		ID: 0,
	}
	go func() {
		time.Sleep(time.Second)
		handler.HandleButton(&disgord.InteractionCreate{
			Message: msg,
			Type:    disgord.InteractionMessageComponent,
			Member:  &disgord.Member{},
			Data: &disgord.ApplicationCommandInteractionData{
				CustomID: "Test",
			},
		})
	}()
	done := false
	handler.RegisterBHandler(msg.ID, func(ic *disgord.InteractionCreate) {
		if ic.Data.CustomID != "Test" {
			t.Errorf("CustomId must be 'Test'")
		} else {
			done = true
		}
	}, 5)
	if !done {
		t.Errorf("Event failed")
	}
}
