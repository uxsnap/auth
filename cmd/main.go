package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	desc "github.com/uxsnap/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthV1Server
}

func (c *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	getResp := &desc.GetResponse{
		Id:        req.GetId(),
		Name:      "Test",
		Email:     "test@mail.ru",
		Role:      0,
		CreatedAt: timestamppb.New(time.Now()),
	}

	return getResp, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))

	if err != nil {
		log.Fatal("Cannot create tcp connection!")
		return
	}

	grpcS := grpc.NewServer()
	reflection.Register(grpcS)

	if err != nil {
		log.Fatal("Cannot create grpc connection!")
		return
	}

	desc.RegisterAuthV1Server(grpcS, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = grpcS.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
