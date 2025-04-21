package api

import (
	"receiver/internal/infra"
	"receiver/internal/service"
)

type Handler struct {
	infra   *infra.Infra
	service *service.Service
}

func NewHandler(infra *infra.Infra, handler *service.Service) *Handler {
	return &Handler{infra: infra, service: handler}
}
