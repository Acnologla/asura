package server

import (
	"asura/src/handler"

	"encoding/json"

	"github.com/andersfylling/disgord"
	"github.com/kevinburke/nacl/sign"
	"github.com/valyala/fasthttp"
)

var PublicKey sign.PublicKey

func Init(publicKey string) {
	PublicKey = []byte(publicKey)
}

func verifyRequest(ctx *fasthttp.RequestCtx) bool {
	signature := ctx.Request.Header.Peek("X-Signature-Ed25519")
	timestamp := ctx.Request.Header.Peek("X-Signature-Timestamp")
	body := ctx.Request.Body()
	return sign.Verify(PublicKey, append(append(signature, timestamp...), body...))
}

func ExecuteInteraction(interaction *disgord.InteractionCreate) *disgord.InteractionResponse {
	if interaction.Type == disgord.InteractionApplicationCommand {
		return handler.Run(interaction)
	}
	return nil
}

func Handler(ctx *fasthttp.RequestCtx) {
	if !verifyRequest(ctx) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}
	if string(ctx.Method()) == fasthttp.MethodPost {
		ctx.SetContentType("application/json")
		var interaction *disgord.InteractionCreate
		json.Unmarshal(ctx.PostBody(), &interaction)
		var response *disgord.InteractionResponse
		if interaction.Type == disgord.InteractionPing {
			response.Type = disgord.InteractionCallbackPong
		} else {
			response = ExecuteInteraction(interaction)
		}
		val, _ := json.Marshal(response)
		ctx.SetBody(val)
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}
