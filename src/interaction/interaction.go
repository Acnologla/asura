package interaction

type InteractionType int
type InteractionCallbackType int

const (
	_ InteractionCallbackType = iota
	PONG
	_
	_
	CHANNEL_MESSAGE_WITH_SOURCE
	DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE
	DEFERRED_UPDATE_MESSAGE
	UPDATE_MESSAGE
	APPLICATION_COMMAND_AUTOCOMPLETE_RESULT
)

const (
	_ InteractionType = iota
	PING
	APPLICATION_COMMAND
	MESSAGE_COMPONENT
	APPLICATION_COMMAND_AUTOCOMPLETE
)

type Snowflake uint64

type InteractionData struct {
	ID      Snowflake                                  `json:"id"`
	Name    string                                     `json:"name"`
	Type    ApplicationCommandType                     `json:"type"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
}

//TODO implement all the struct
type Interaction struct {
	ID            Snowflake        `json:"id"`
	ApplicationID Snowflake        `json:"application_id"`
	Type          InteractionType  `json:"type"`
	Data          *InteractionData `json:"data"`
}

//TODO implement all the struct

type InteractionCallbackData struct {
	Tts     bool     `json:"tts"`
	Content string   `json:"content"`
	Embeds  []*Embed `json:"embeds"`
}

type InteractionResponse struct {
	Type InteractionCallbackType  `json:"type"`
	Data *InteractionCallbackData `json:"data"`
}
