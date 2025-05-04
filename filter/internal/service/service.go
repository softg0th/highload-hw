package service

import (
	"filter/internal/infra"
	"filter/internal/pkg"
	logstash_logger "github.com/KaranJagtiani/go-logstash"
)

type Service struct {
	infra          *infra.Infra
	processingPool *pkg.Pool[[]byte]
	Logger         *logstash_logger.Logstash
}

func NewService(infra *infra.Infra, pool *pkg.Pool[[]byte], logger *logstash_logger.Logstash) (*Service, error) {
	return &Service{
		infra:          infra,
		processingPool: pool,
		Logger:         logger,
	}, nil
}
