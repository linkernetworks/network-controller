package main

import (
	"bytes"
	"encoding/binary"
	"log"

	pb "github.com/linkernetworks/network-controller/messages"

	"github.com/linkernetworks/go-openvswitch/ovs"
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

func (s *server) CreateBridge(ctx context.Context, req *pb.CreateBridgeRequest) (*pb.Response, error) {
	if err := s.OVS.CreateBridge(req.BridgeName, req.DatapathType); err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}
	return &pb.Response{
		Success: true,
		Reason:  "",
	}, nil
}

func (s *server) DeleteBridge(ctx context.Context, req *pb.DeleteBridgeRequest) (*pb.Response, error) {
	if err := s.OVS.DeleteBridge(req.BridgeName); err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}
	return &pb.Response{
		Success: true,
		Reason:  "",
	}, nil
}

func (s *server) AddDPDKPort(ctx context.Context, req *pb.AddPortRequest) (*pb.Response, error) {
	if err := s.OVS.AddDPDKPort(req.BridgeName, req.IfaceName, req.DpdkDevargs); err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}
	return &pb.Response{
		Success: true,
		Reason:  "",
	}, nil
}

func (s *server) AddPort(ctx context.Context, req *pb.AddPortRequest) (*pb.Response, error) {
	if err := s.OVS.AddPort(req.BridgeName, req.IfaceName); err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}
	return &pb.Response{
		Success: true,
		Reason:  "",
	}, nil
}

func (s *server) GetPort(ctx context.Context, req *pb.GetPortRequest) (*pb.GetPortResponse, error) {
	portOptions, err := s.OVS.GetPort(req.IfaceName)
	if err != nil {
		return &pb.GetPortResponse{
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}

	options := &pb.PortOptions{}
	if portOptions.Tag != nil {
		options.Tag = int32(*portOptions.Tag)
	}
	if portOptions.VLANMode != nil {
		options.VLANMode = *portOptions.VLANMode
	}
	// convert response trunk data(int) to grpc message trunk format(int32)
	for _, t := range portOptions.Trunk {
		options.Trunk = append(options.Trunk, int32(t))
	}
	return &pb.GetPortResponse{
		PortOptions: options,
		ServerResponse: &pb.Response{
			Success: true,
			Reason:  "",
		},
	}, nil
}

func (s *server) SetPort(ctx context.Context, req *pb.SetPortRequest) (*pb.Response, error) {
	portOptions := ovs.PortOptions{}
	if req.Options.VLANMode == "" {
		// set with vlan tag
		tag := int(req.Options.Tag)
		portOptions.Tag = &tag
		portOptions.VLANMode = nil
		portOptions.Trunk = nil
	} else {
		// set with vlan trunk
		portOptions.Tag = nil
		VLANMode := req.Options.VLANMode
		portOptions.VLANMode = &VLANMode
		// convert grpc message trunk format(int32) to request trunk data(int)
		for _, t := range req.Options.Trunk {
			portOptions.Trunk = append(portOptions.Trunk, int(t))
		}
	}

	if err := s.OVS.SetPort(req.IfaceName, portOptions); err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}
	return &pb.Response{
		Success: true,
		Reason:  "",
	}, nil
}

func (s *server) DeletePort(ctx context.Context, req *pb.DeletePortRequest) (*pb.Response, error) {
	if err := s.OVS.DeletePort(req.BridgeName, req.IfaceName); err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}
	return &pb.Response{
		Success: true,
		Reason:  "",
	}, nil
}

func (s *server) AddFlow(ctx context.Context, req *pb.AddFlowRequest) (*pb.Response, error) {
	if err := s.OVS.AddFlow(req.BridgeName, req.FlowString); err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}
	return &pb.Response{
		Success: true,
		Reason:  "",
	}, nil
}

func (s *server) DeleteFlow(ctx context.Context, req *pb.DeleteFlowRequest) (*pb.Response, error) {
	if err := s.OVS.DeleteFlow(req.BridgeName, req.FlowString); err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}
	return &pb.Response{
		Success: true,
		Reason:  "",
	}, nil
}

func (s *server) DumpFlows(ctx context.Context, req *pb.DumpFlowsRequest) (*pb.DumpFlowsResponse, error) {
	flows, err := s.OVS.DumpFlows(req.BridgeName)
	if err != nil {
		return &pb.DumpFlowsResponse{
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}

	flowsBytes := make([][]byte, len(flows))
	for _, flow := range flows {
		bytes, err := flow.MarshalText()
		if err != nil {
			return &pb.DumpFlowsResponse{
				ServerResponse: &pb.Response{
					Success: false,
					Reason:  err.Error(),
				},
			}, err
		}

		flowsBytes = append(flowsBytes, bytes)
	}

	return &pb.DumpFlowsResponse{
		Flows: flowsBytes,
		ServerResponse: &pb.Response{
			Success: true,
			Reason:  "",
		},
	}, nil
}

func (s *server) DumpPorts(ctx context.Context, req *pb.DumpPortsRequest) (*pb.DumpPortsResponse, error) {
	ports, err := s.OVS.DumpPorts(req.BridgeName)
	if err != nil {
		return &pb.DumpPortsResponse{
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}

	portsBytes := make([][]byte, len(ports))
	for _, port := range ports {

		buf := &bytes.Buffer{}
		if err := binary.Write(buf, binary.BigEndian, port); err != nil {
			return &pb.DumpPortsResponse{
				ServerResponse: &pb.Response{
					Success: false,
					Reason:  err.Error(),
				},
			}, err
		}
		portsBytes = append(portsBytes, buf.Bytes())
	}

	return &pb.DumpPortsResponse{
		Ports: portsBytes,
		ServerResponse: &pb.Response{
			Success: true,
			Reason:  "",
		},
	}, nil
}

func (s *server) DumpPort(ctx context.Context, req *pb.DumpPortRequest) (*pb.DumpPortResponse, error) {
	port, err := s.OVS.DumpPort(req.BridgeName, req.PortName)
	if err != nil {
		return &pb.DumpPortResponse{
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}

	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, port); err != nil {
		return &pb.DumpPortResponse{
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}

	return &pb.DumpPortResponse{
		Port: buf.Bytes(),
		ServerResponse: &pb.Response{
			Success: false,
			Reason:  "",
		},
	}, nil
}
