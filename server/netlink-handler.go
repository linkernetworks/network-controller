package main

import (
	"crypto/sha256"
	"encoding/hex"

	pb "github.com/linkernetworks/network-controller/messages"
	ovs "github.com/linkernetworks/network-controller/openvswitch"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/linkernetworks/network-controller/docker"
	"github.com/linkernetworks/network-controller/link"
	"golang.org/x/net/context"
)

func (s *server) FindNetworkNamespacePath(ctx context.Context, req *pb.FindNetworkNamespacePathRequest) (*pb.FindNetworkNamespacePathResponse, error) {

	containers, err := docker.ListContainer()
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

	containerInfo, err := docker.InspectContainer(containerID)
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

	netns, err := ns.GetNS(req.Path)
	if err != nil {
		return &pb.ConnectBridgeResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	hash := sha256.New()
	hash.Write([]byte(req.PodUUID))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	hostVethName := "veth" + mdStr[0:8]

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

	if err := ovs.AddPort(req.BridgeName, hostVethName); err != nil {
		return &pb.ConnectBridgeResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	return &pb.ConnectBridgeResponse{Success: true}, nil
}
