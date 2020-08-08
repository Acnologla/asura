package test

import (
	"asura/src/handler"
	"testing"
)

func TestParser(t *testing.T) {
	t.Run("Test Errored Parsing", func(t *testing.T) {
		args := []string {"atapo", "7", "12", "<@!1239012>"}
		if handler.ValidateNumber(args[0], 0, 10) == nil {
			t.Errorf("Expected error in %s",args[0])
		}
		if handler.ValidateString(args[1], 0, 10) != nil {
			t.Errorf("Expected sucess in %s",args[1])
		}
		if handler.ValidateNumber(args[2], 0, 10) == nil {
			t.Errorf("Expected error in %s",args[2])
		}
		if handler.ValidateMention(args[3]) == nil {
			t.Errorf("Expected success in %s",args[3])
		}
	})
}

