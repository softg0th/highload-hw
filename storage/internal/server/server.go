package server

import (
	infr "storage/internal/infra"
	pb "storage/internal/proto"
)

type Server struct {
	pb.UnimplementedStorageServiceServer
	Infra *infr.Infra
}

func NewServer(s *infr.Infra) *Server {
	return &Server{Infra: s}
}
