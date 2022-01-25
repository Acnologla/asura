package utils

import (
	"asura/src/interaction"
	"encoding/json"
)

func GenerateOptions(options ...*interaction.ApplicationCommandOption) []*interaction.ApplicationCommandOption {
	return options
}

func GetUser(options []*interaction.ApplicationCommandInteractionDataOption, i int) *interaction.User {
	val, _ := json.Marshal(options[i].Value)
	user := &interaction.User{}
	json.Unmarshal(val, user)
	return user
}
