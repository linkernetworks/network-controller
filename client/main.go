package main

import (
	"log"
	"time"

	pb "github.com/linkernetworks/network-controller/messages"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNetworkControlClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Ping(ctx, &pb.PingRequest{Ping: "PING"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Got: %s", r.Pong)
}
