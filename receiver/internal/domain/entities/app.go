package entities

import (
	"github.com/valyala/fasthttp"
)

type App struct {
	Handler fasthttp.RequestHandler
}

func NewApp(handler fasthttp.RequestHandler) *App {
	return &App{
		Handler: handler,
	}
}
