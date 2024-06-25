package grpc

import (
	"antrein/dd-dashboard-config/application/common/resource"
	"antrein/dd-dashboard-config/application/common/usecase"
	"antrein/dd-dashboard-config/model/config"
	"context"

	pb "github.com/antrein/proto-repository/pb/bc"
	"google.golang.org/grpc"
)

type helloServer struct {
	pb.UnimplementedGreeterServer
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func ApplicationDelegate(cfg *config.Config, uc *usecase.CommonUsecase, rsc *resource.CommonResource) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()

	// Hello service
	helloServer := &helloServer{}
	pb.RegisterGreeterServer(grpcServer, helloServer)

	return grpcServer, nil
}
