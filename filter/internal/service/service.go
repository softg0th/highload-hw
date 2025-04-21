package service

import (
	"filter/internal/infra"
	"filter/internal/pkg"
)

type Service struct {
	infra          *infra.Infra
	processingPool *pkg.Pool[[]byte]
}

func NewService(infra *infra.Infra, pool *pkg.Pool[[]byte]) (*Service, error) {
	return &Service{
		infra:          infra,
		processingPool: pool,
	}, nil
}
