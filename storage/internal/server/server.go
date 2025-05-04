package server

import (
	logstash_logger "github.com/KaranJagtiani/go-logstash"
	infr "storage/internal/infra"
	pb "storage/internal/proto"
)

type Server struct {
	pb.UnimplementedStorageServiceServer
	Infra  *infr.Infra
	Logger *logstash_logger.Logstash
}

func NewServer(s *infr.Infra, logger *logstash_logger.Logstash) *Server {
	return &Server{Infra: s, Logger: logger}
}
