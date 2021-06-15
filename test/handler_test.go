package test

import (
	"asura/src/handler"
	"testing"
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
