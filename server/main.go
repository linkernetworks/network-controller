package main

import (
	"log"
	"net"
	"encoding/json"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/network-controller/ovs"
  "github.com/linkernetworks/network-controller/utils"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement messages.NetworkServer.
type server struct{}

func (s *server) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Word: "Echo Response: " + req.Word,
	}, nil
}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	log.Printf("Client Sent: %s", req.Ping)
	return &pb.PingResponse{Pong: "PONG"}, nil
}

func (s *server) AddPort(ctx context.Context, req *pb.AddPortRequest) (*pb.OVSResponse, error) {
	if err := ovs.AddPort(req.BridgeName, req.IfaceName); err != nil {
		return &pb.OVSResponse{
			Success: false, Reason: err.Error(),
		}, err
	}
	return &pb.OVSResponse{Success: true}, nil
}

func (s *server) DeletePort(ctx context.Context, req *pb.DeletePortRequest) (*pb.OVSResponse, error) {
	if err := ovs.DeletePort(req.BridgeName, req.IfaceName); err != nil {
		return &pb.OVSResponse{
			Success: false, Reason: err.Error(),
		}, err
	}
	return &pb.OVSResponse{Success: true}, nil
}

func (s *server) AddFlow(ctx context.Context, req *pb.AddFlowRequest) (*pb.OVSResponse, error) {
  var flow map[string]interface{}
  if err :=  json.Unmarshal([]byte(req.FlowString), &flow); err != nil {
    if err := ovs.AddFlow(req.BridgeName, utils.ConvertOVSFlow(flow)); err != nil {
	  	return &pb.OVSResponse{
	  		Success: false, Reason: err.Error(),
	  	}, err
	  }
  }
	return &pb.OVSResponse{Success: true}, nil
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
