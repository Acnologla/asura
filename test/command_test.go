package test

import (
	_ "asura/src/commands"
	"asura/src/handler"
	"testing"
)

func includes(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func TestCommandNoDupe(t *testing.T) {
	var aliases []string
	for _, command := range handler.Commands {
		for _, aliase := range command.Aliases {
			if includes(aliases, aliase) {
				t.Errorf("Dup aliases found, %s", aliase)
			}
			aliases = append(aliases, aliase)
		}
	}
}

func TestCommandLevenshtein(t *testing.T) {
	for _, command := range handler.Commands {
		for _, command2 := range handler.Commands {
			if command.Aliases[0] != command2.Aliases[0] {
				if 1 >= handler.CompareStrings(command.Aliases[0], command2.Aliases[0]) {
					t.Errorf("Dup Levenshtein found, %s ,%s", command.Aliases[0], command2.Aliases[0])
				}
			}
		}
	}
}

func TestCommandFields(t *testing.T) {
	for _, command := range handler.Commands {
		if len(command.Aliases) == 0 {
			t.Errorf("Commands must have at least one aliase")
		}
		if command.Cooldown <= 0 {
			t.Errorf(command.Aliases[0] + " - The cooldown should be greater than 0")
		}
		if command.Help == "" {
			t.Errorf(command.Aliases[0] + " - The command must have a help")
		}
		if command.Usage == "" {
			t.Errorf(command.Aliases[0] + " - The command must have a usage example")
		}
		if command.Run == nil {
			t.Errorf(command.Aliases[0] + " - The command must have a run function")

		}
	}
}
