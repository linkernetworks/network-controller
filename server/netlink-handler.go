package main

import (
	"net"
	"runtime"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/network-controller/utils"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/plugins/pkg/ipam"
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
	if containerID == "" {
		return &pb.FindNetworkNamespacePathResponse{
			Success: false, Reason: "ContainerID is empty.",
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

func (s *server) ConfigureIface(ctx context.Context, req *pb.ConfigureIfaceRequest) (*pb.ConfigureIfaceResponse, error) {
	runtime.LockOSThread()
	netns, err := ns.GetNS(req.Path)
	if err != nil {
		return &pb.ConfigureIfaceResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	err = netns.Do(func(_ ns.NetNS) error {
		result := &current.Result{}
		result.Interfaces = []*current.Interface{{Name: req.ContainerVethName}}

		ipv4, err := types.ParseCIDR(req.IP)
		if err != nil {
			return err
		}
		result.IPs = []*current.IPConfig{
			{
				Version:   "4",
				Interface: current.Int(0),
				Address:   *ipv4,
				Gateway:   net.ParseIP(req.Gateway),
			},
		}

		_, ipv4Net, err := net.ParseCIDR(req.IP)
		gatewayAddr, _, err := net.ParseCIDR(req.Gateway)
		if err != nil {
			return err
		}
		result.Routes = []*types.Route{{Dst: *ipv4Net, GW: gatewayAddr}}

		if err := ipam.ConfigureIface(req.ContainerVethName, result); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.ConfigureIfaceResponse{
			Success: false, Reason: err.Error(),
		}, err
	}

	return &pb.ConfigureIfaceResponse{Success: true}, nil
}
