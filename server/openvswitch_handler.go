package main

import (
	"log"

	pb "github.com/linkernetworks/network-controller/messages"

	"golang.org/x/net/context"
)

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
	if err := s.OVS.AddPort(req.BridgeName, req.IfaceName); err != nil {
		return &pb.OVSResponse{
			Success: false, Reason: err.Error(),
		}, err
	}
	return &pb.OVSResponse{Success: true}, nil
}

func (s *server) DeletePort(ctx context.Context, req *pb.DeletePortRequest) (*pb.OVSResponse, error) {
	if err := s.OVS.DeletePort(req.BridgeName, req.IfaceName); err != nil {
		return &pb.OVSResponse{
			Success: false, Reason: err.Error(),
		}, err
	}
	return &pb.OVSResponse{Success: true}, nil
}

func (s *server) AddFlow(ctx context.Context, req *pb.AddFlowRequest) (*pb.OVSResponse, error) {
	if err := s.OVS.AddFlow(req.BridgeName, req.FlowString); err != nil {
		return &pb.OVSResponse{
			Success: false, Reason: err.Error(),
		}, err
	}
	return &pb.OVSResponse{Success: true}, nil
}

func (s *server) DeleteFlow(ctx context.Context, req *pb.DeleteFlowRequest) (*pb.OVSResponse, error) {
	if err := s.OVS.DeleteFlow(req.BridgeName, req.FlowString); err != nil {
		return &pb.OVSResponse{
			Success: false, Reason: err.Error(),
		}, err
	}
	return &pb.OVSResponse{Success: true}, nil
}
