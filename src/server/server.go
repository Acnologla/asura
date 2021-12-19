package server

import (
	"asura/src/handler"
	interactionPkg "asura/src/interaction"

	"encoding/json"

	"github.com/kevinburke/nacl/sign"
	"github.com/valyala/fasthttp"
)

var PublicKey sign.PublicKey

func Init(publicKey string) {
	PublicKey = []byte(publicKey)
}

// TODO verify request signature
func verifyRequest(ctx *fasthttp.RequestCtx) bool {
	signature := ctx.Request.Header.Peek("X-Signature-Ed25519")
	timestamp := ctx.Request.Header.Peek("X-Signature-Timestamp")
	body := ctx.Request.Body()
	return sign.Verify(PublicKey, append(append(signature, timestamp...), body...))
}

func Handler(ctx *fasthttp.RequestCtx) {
	print("a")
	if !verifyRequest(ctx) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}
	if string(ctx.Method()) == fasthttp.MethodPost {
		ctx.SetContentType("application/json")
		var interaction interactionPkg.Interaction
		json.Unmarshal(ctx.PostBody(), &interaction)
		var response *interactionPkg.InteractionResponse
		if interaction.Type == interactionPkg.PING {
			response.Type = interactionPkg.PONG
		}
		if interaction.Type == interactionPkg.APPLICATION_COMMAND {
			response = handler.Run(interaction)
		}
		val, _ := json.Marshal(response)
		ctx.SetBody(val)
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}
