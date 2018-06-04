package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/linkernetworks/network-controller/messages"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement messages.NetworkServer.
type server struct{}

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

	// Stop all listener by catching interrupt signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal, lis net.Listener, s *grpc.Server) {
		sig := <-c
		log.Printf("caught signal: %s", sig.String())

		log.Printf("stopping tcp listener...")
		lis.Close()

		log.Printf("stopping grpc server...")
		s.Stop()

		log.Printf("all listener are stopped successfully")
		os.Exit(0)
	}(sigc, lis, s)
}
