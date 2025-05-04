package infra

import (
	"context"
	"filter/internal/entities"
	pb "filter/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RPCConn struct {
	stream pb.StorageService_StreamWithAckClient
}

func NewRPCConn(grpcAddress string) (*RPCConn, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewStorageServiceClient(conn)
	stream, err := client.StreamWithAck(context.Background())
	if err != nil {
		return nil, err
	}

	return &RPCConn{
		stream: stream,
	}, nil
}

func (r *RPCConn) StreamRequest(message entities.Message) error {
	err := r.stream.Send(&pb.GetMessageRequest{
		UserId:    message.UserId,
		Text:      message.Text,
		Timestamp: timestamppb.New(message.Timestamp),
	})
	if err != nil {
		return err
	}
	return nil
}
