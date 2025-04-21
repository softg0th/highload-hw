package entities

import (
	pb "storage/internal/proto"
)

type Server struct {
	pb.UnimplementedStorageServiceServer
}
