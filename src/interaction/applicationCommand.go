package interaction

import "github.com/andersfylling/disgord"

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

type ChannelType int

const (
	GUILD_TEXT ChannelType = iota
	DM
	GUILD_VOICE
	GROUP_DM
	GUILD_CATEGORY
	GUILD_NEWS
	GUILD_STORE
	_
	_
	_
	GUILD_NEWS_THREAD
	GUILD_PUBLIC_THREAD
	GUILD_PRIVATE_THREAD
	GUILD_STAGE_VOICE
)

type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name"`
	Type    ApplicationCommandOptionType               `json:"type"`
	Value   interface{}                                `json:"value"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
	Focused bool                                       `json:"focused"`
}

type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type ApplicationCommandOption struct {
	Type         ApplicationCommandOptionType      `json:"type"`
	Name         string                            `json:"name"`
	Description  string                            `json:"description"`
	Required     bool                              `json:"required"`
	Choices      []*ApplicationCommandOptionChoice `json:"choices"`
	Options      []*ApplicationCommandOption       `json:"options"`
	ChannelTypes []ChannelType                     `json:"channel_types"`
	MinValue     int                               `json:"min_value"`
	MaxValue     int                               `json:"max_value"`
	AutoComplete bool                              `json:"auto_complete"`
}

type ApplicationCommand struct {
	ID                disgord.Snowflake           `json:"id"`
	Type              ApplicationCommandType      `json:"type"`
	ApplicationID     disgord.Snowflake           `json:"application_id"`
	GuildId           disgord.Snowflake           `json:"guild_id"`
	Name              string                      `json:"name"`
	Description       string                      `json:"description"`
	Options           []*ApplicationCommandOption `json:"options"`
	DefaultPermission bool                        `json:"default_permission,omitempty"`
	Version           disgord.Snowflake           `json:"version"`
}
