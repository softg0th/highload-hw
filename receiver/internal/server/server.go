package server

import (
	"context"
	"github.com/valyala/fasthttp"
	"log"
	"receiver/internal/api"
	"receiver/internal/domain/entities"
	"receiver/internal/infra"
	"receiver/internal/service"
)

type Server struct {
	app     *entities.App
	handler *api.Handler
	server  *fasthttp.Server
}

func setupRouter(handler *api.Handler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/api":
			if ctx.IsGet() {
				handler.DebugGet(ctx)
				return
			}
		case "/api/message":
			if ctx.IsPost() {
				handler.ReceiveMessage(ctx)
				return
			}
		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("Not Found")
		}
	}
}

func Setup(infra *infra.Infra, service *service.Service) *Server {
	handler := api.NewHandler(infra, service)
	router := setupRouter(handler)
	app := entities.NewApp(router)

	s := &fasthttp.Server{
		Handler: app.Handler,
		Name:    "fasthttp-server",
	}

	return &Server{
		app:     app,
		handler: handler,
		server:  s,
	}
}

func (s *Server) Start(port string) error {
	return s.server.ListenAndServe(":" + port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		if err := s.server.Shutdown(); err != nil {
			log.Fatalf("server shutdown error:%v", err)
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
