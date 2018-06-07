package main

import (
	"flag"
	"log"
	"time"

	pb "github.com/linkernetworks/network-controller/messages"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var serverAddr string
	flag.StringVar(&serverAddr, "server", "", "target server address, [ip:port] for TCP or unix://[path] for UNIX")
	flag.Parse()

	if serverAddr == "" {
		log.Fatalf("You should use the -server to specify the server address, 0.0.0.0:50051 for TCP and unix:///tmp/xxx.sock for UNIX")
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
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
