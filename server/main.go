package main

import (
	"log"
	"net"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/network-controller/ovs"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement messages.NetworkServer.
type server struct{}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	return &pb.PingReply{Message: "Hello " + req.Name}, nil
}

func (s *server) AddPort(ctx context.Context, req *pb.AddPortRequest) (*pb.OVSReply, error) {
	if err := ovs.AddPort(req.BridgeName, req.IfaceName); err != nil {
		return &pb.OVSReply{Message: "failed"}, err
	}
	return &pb.OVSReply{Message: "Succeeded"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNetworkControlServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
