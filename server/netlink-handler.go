package main

import (
	"runtime"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/network-controller/utils"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/linkernetworks/network-controller/docker"
	"github.com/linkernetworks/network-controller/link"
	"golang.org/x/net/context"
)

func (s *server) FindNetworkNamespacePath(ctx context.Context, req *pb.FindNetworkNamespacePathRequest) (*pb.FindNetworkNamespacePathResponse, error) {
	cli, err := docker.New()
	if err != nil {
		return nil, err
	}

	containers, err := cli.ListContainer()
	if err != nil {
		return &pb.FindNetworkNamespacePathResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	containerID, err := docker.FindK8SPauseContainerID(containers, req.PodName, req.Namespace, req.PodUUID)
	if err != nil {
		return &pb.FindNetworkNamespacePathResponse{
			Success: false, Reason: err.Error(),
		}, err
	}
	if containerID == "" {
		return &pb.FindNetworkNamespacePathResponse{
			Success: false, Reason: "ContainerID is empty.",
		}, err
	}

	containerInfo, err := cli.InspectContainer(containerID)
	if err != nil {
		return &pb.FindNetworkNamespacePathResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	return &pb.FindNetworkNamespacePathResponse{
		Success: true,
		Path:    docker.GetSandboxKey(containerInfo),
	}, err
}

func (s *server) ConnectBridge(ctx context.Context, req *pb.ConnectBridgeRequest) (*pb.ConnectBridgeResponse, error) {
	runtime.LockOSThread()
	netns, err := ns.GetNS(req.Path)
	if err != nil {
		return &pb.ConnectBridgeResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	hostVethName := utils.GenerateVethName(req.PodUUID)
	err = netns.Do(func(hostNS ns.NetNS) error {
		if _, _, err := link.SetupVeth(req.ContainerVethName, hostVethName, 1500, hostNS); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.ConnectBridgeResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	if err := s.OVS.AddPort(req.BridgeName, hostVethName); err != nil {
		return &pb.ConnectBridgeResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	return &pb.ConnectBridgeResponse{Success: true}, nil
}
