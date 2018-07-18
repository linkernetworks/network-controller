package main

import (
	"log"
	"runtime"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/network-controller/utils"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/plugins/pkg/ipam"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/linkernetworks/network-controller/docker"
	"github.com/linkernetworks/network-controller/nl"
	"golang.org/x/net/context"
)

func (s *server) FindNetworkNamespacePath(ctx context.Context, req *pb.FindNetworkNamespacePathRequest) (*pb.FindNetworkNamespacePathResponse, error) {
	log.Println("Start to Find Network")
	cli, err := docker.New()
	if err != nil {
		return &pb.FindNetworkNamespacePathResponse{
			Path: "",
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}

	containers, err := cli.ListContainer()
	if err != nil {
		return &pb.FindNetworkNamespacePathResponse{
			Path: "",
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}

	containerID, err := docker.FindK8SPauseContainerID(containers, req.PodName, req.Namespace, req.PodUUID)
	if err != nil {
		return &pb.FindNetworkNamespacePathResponse{
			Path: "",
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}
	if containerID == "" {
		return &pb.FindNetworkNamespacePathResponse{
			Path: "",
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}

	containerInfo, err := cli.InspectContainer(containerID)
	if err != nil {
		return &pb.FindNetworkNamespacePathResponse{
			Path: "",
			ServerResponse: &pb.Response{
				Success: false,
				Reason:  err.Error(),
			},
		}, err
	}

	return &pb.FindNetworkNamespacePathResponse{
		Path: docker.GetSandboxKey(containerInfo),
		ServerResponse: &pb.Response{
			Success: true,
			Reason:  "",
		},
	}, err
}

func (s *server) ConnectBridge(ctx context.Context, req *pb.ConnectBridgeRequest) (*pb.Response, error) {
	log.Println("Start to Connect Bridge")
	runtime.LockOSThread()
	netns, err := ns.GetNS(req.Path)
	if err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}

	log.Println("Get the netns object success")

	hostVethName := utils.GenerateVethName(req.PodUUID, req.ContainerVethName)
	log.Println("The host veth name", hostVethName)
	err = netns.Do(func(hostNS ns.NetNS) error {
		if _, _, err := nl.SetupVeth(req.ContainerVethName, hostVethName, 1500, hostNS); err != nil {
			return err
		}
		return nil
	})
	log.Println("Success setup veth")
	if err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}

	log.Println("Try to add port", hostVethName, " To ", req.BridgeName)
	if err := s.OVS.AddPort(req.BridgeName, hostVethName); err != nil {
		log.Println("Add port fail:", err, req.BridgeName, hostVethName)
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
		}, err
	}

	log.Println("Add Port Success")
	return &pb.Response{
		Success: true,
		Reason:  "",
	}, nil
}

func (s *server) ConfigureIface(ctx context.Context, req *pb.ConfigureIfaceRequest) (*pb.Response, error) {
	runtime.LockOSThread()
	log.Println("Start to configure interface")
	netns, err := ns.GetNS(req.Path)
	if err != nil {
		return &pb.Response{
			Success: false,
			Reason:  err.Error(),
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
			},
		}

		return ipam.ConfigureIface(req.ContainerVethName, result)
	})
	if err != nil {
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
