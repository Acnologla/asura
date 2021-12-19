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

//TODO implement all the struct
type Interaction struct {
	ID            Snowflake       `json:"id"`
	ApplicationID Snowflake       `json:"application_id"`
	Type          InteractionType `json:"type"`
}

//TODO implement all the struct

type InteractionResponse struct {
	Type InteractionCallbackType `json:"type"`
}
