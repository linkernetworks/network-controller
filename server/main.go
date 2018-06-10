package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/network-controller/nl"
	ovs "github.com/linkernetworks/network-controller/openvswitch"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement messages.NetworkServer.
type server struct {
	OVS *ovs.OVSManager
}

func main() {

	var tcpAddr string
	var unixPath string
	var nlEventTracker bool
	flag.StringVar(&tcpAddr, "tcp", "", "Run as a TCP server and listen on target address")
	flag.StringVar(&unixPath, "unix", "", "Run as a UNIX server and listen on target path")
	flag.BoolVar(&nlEventTracker, "nl", false, "Run as a Netlink Event Tracker")

	flag.Parse()

	if tcpAddr == "" && unixPath == "" {
		log.Fatalf("You must use the one method(-tcp/-unix) to decide how the server listen to")
	}

	if tcpAddr != "" && unixPath != "" {
		log.Fatalf("You should only choose one method to listen to")
	}

	//Listen
	var lis net.Listener
	var err error
	if tcpAddr != "" {
		lis, err = net.Listen("tcp", tcpAddr)
	} else {
		lis, err = net.Listen("unix", unixPath)
	}

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPC
	s := grpc.NewServer()
	pb.RegisterNetworkControlServer(s, &server{OVS: ovs.New()})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	//for process
	stop := make(chan struct{})

	//The netlink event tracker
	var tracker *nl.NlEventHandler
	if nlEventTracker {
		tracker = nl.New()
		tracker.AddDeletedLinkHandler(nl.RemoveVethFromOVS)
	}
	// Stop all listener by catching interrupt signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal, lis net.Listener, s *grpc.Server) {
		sig := <-c
		log.Printf("caught signal: %s", sig.String())

		log.Printf("stopping grpc server...")
		s.GracefulStop()

		log.Printf("stopping tcp listener...")
		lis.Close()

		if tracker != nil {
			log.Printf("stopping netlink event tracker...")
			tracker.Stop()
		}

		if unixPath != "" {
			os.RemoveAll(unixPath)
		}

		log.Printf("all listener are stopped successfully")
		close(stop)
	}(sigc, lis, s)

	if tracker != nil {
		log.Printf("Starting the Netlink Event Tracker")
		go tracker.TrackNetlink()
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	<-stop
	os.Exit(0)
}
