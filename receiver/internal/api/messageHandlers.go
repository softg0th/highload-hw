package api

import (
	"encoding/json"
	"fmt"
	"log"
	"receiver/internal/domain/entities"

	"github.com/valyala/fasthttp"
)

func (h *Handler) DebugGet(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.Write([]byte(`{"DEBUG":"ok"}`))
}

func (h *Handler) ReceiveMessage(ctx *fasthttp.RequestCtx) {
	var msg entities.Message

	if err := json.Unmarshal(ctx.PostBody(), &msg); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.Write([]byte(`{"error":"Not all fields are filled"}`))
		return
	}
	_, err := h.service.Sender.Send(ctx, msg)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetContentType("application/json")
		ctx.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	log.Println("New message!")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.Write([]byte(`{"DEBUG":"ok"}`))
}
