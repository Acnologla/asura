package interaction

type ApplicationCommandType int

const (
	_ ApplicationCommandType = iota
	CHAT_INPUT
	USERCOMMANDTYPE
	MESSAGE
)

type ApplicationCommandOptionType int

const (
	_ ApplicationCommandOptionType = iota
	SUB_COMMAND
	SUB_COMMAND_GROUP
	STRING
	INTEGER
	BOOLEAN
	USER
	CHANNEL
	ROLE
	MENTIONABLE
	NUMBER
)

type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name"`
	Type    ApplicationCommandOptionType               `json:"type"`
	Value   interface{}                                `json:"value"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
	Focused bool                                       `json:"focused"`
}
