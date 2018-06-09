package main

import (
	"bytes"
	"encoding/binary"
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

func (s *server) DumpFlows(ctx context.Context, req *pb.DumpFlowsRequest) (*pb.DumpFlowsResponse, error) {
	flows, err := s.OVS.DumpFlows(req.BridgeName)
	if err != nil {
		return &pb.DumpFlowsResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	flowsBytes := make([][]byte, len(flows))
	for _, flow := range flows {
		bytes, err := flow.MarshalText()
		if err != nil {
			return &pb.DumpFlowsResponse{
				Success: false, Reason: err.Error(),
			}, err
		}

		flowsBytes = append(flowsBytes, bytes)
	}

	return &pb.DumpFlowsResponse{
		Success: true,
		Flows:   flowsBytes,
	}, nil
}

func (s *server) DumpPorts(ctx context.Context, req *pb.DumpPortsRequest) (*pb.DumpPortsResponse, error) {
	ports, err := s.OVS.DumpPorts(req.BridgeName)
	if err != nil {
		return &pb.DumpPortsResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	portsBytes := make([][]byte, len(ports))
	for _, port := range ports {

		buf := &bytes.Buffer{}
		if err := binary.Write(buf, binary.BigEndian, port); err != nil {
			return &pb.DumpPortsResponse{
				Success: false, Reason: err.Error(),
			}, err
		}

		portsBytes = append(portsBytes, buf.Bytes())
	}

	return &pb.DumpPortsResponse{
		Success: true,
		Ports:   portsBytes,
	}, nil
}

func (s *server) DumpPort(ctx context.Context, req *pb.DumpPortRequest) (*pb.DumpPortResponse, error) {
	port, err := s.OVS.DumpPort(req.BridgeName, req.PortName)
	if err != nil {
		return &pb.DumpPortResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, port); err != nil {
		return &pb.DumpPortResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	return &pb.DumpPortResponse{
		Success: true,
		Port:    buf.Bytes(),
	}, nil
}
