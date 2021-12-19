package server

import (
	"asura/commandHandler"
	interactionPkg "asura/interaction"
	"encoding/json"

	"github.com/valyala/fasthttp"
)

var PublicKey string

// TODO verify request signature
func verifyRequest(ctx *fasthttp.RequestCtx) bool {
	return true
}

func Handler(ctx *fasthttp.RequestCtx) {
	if !verifyRequest(ctx) {
		return
	}
	if string(ctx.Method()) == fasthttp.MethodPost {
		ctx.SetContentType("application/json")
		var interaction interactionPkg.Interaction
		json.Unmarshal(ctx.PostBody(), &interaction)
		var response interactionPkg.InteractionResponse
		if interaction.Type == interactionPkg.PING {
			response.Type = interactionPkg.PONG
			val, _ := json.Marshal(response)
			ctx.SetBody(val)
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}
		if interaction.Type == interactionPkg.APPLICATION_COMMAND {
			commandHandler.Run(interaction)
		}
	}
}
