package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/linkernetworks/go-openvswitch/ovs"

	flags "github.com/jessevdk/go-flags"
	pb "github.com/linkernetworks/network-controller/messages"
	"google.golang.org/grpc"
)

type clientOptions struct {
	Server string `short:"s" long:"server " description:"target server address, [ip:port] for TCP or unix://[path] for UNIX" required:"true"`
}

var options clientOptions
var parser = flags.NewParser(&options, flags.Default)

// Testing Variable
const (
	bridgeName = "switch1"
)

//

type TestFunc func(*context.Context, *pb.NetworkControlClient) error

func Test(name string, ctx *context.Context, conn *pb.NetworkControlClient, fun TestFunc) {
	log.Println("[Start] Testing:", name)
	if err := fun(ctx, conn); err != nil {
		log.Fatal("[Fail] Testing:", name, err)
	}
	log.Println("[Pass] Testing:", name)
}

func main() {
	if _, err := parser.Parse(); err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	log.Println("Start to connect to", options.Server)
	conn, err := grpc.Dial(options.Server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Println("PASS")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("New the Network Controller Client")
	ncClient := pb.NewNetworkControlClient(conn)
	if ncClient == nil {
		log.Fatalf("Init NewNetworkControlClient Fail")
	}
	log.Println("PASS")

	Test("ping", &ctx, &ncClient, func(ctx *context.Context, nc *pb.NetworkControlClient) error {
		_, err := ncClient.Ping(
			*ctx,
			&pb.PingRequest{
				Ping: "Test Ping",
			},
		)
		return err
	})

	Test("Create Bridge", &ctx, &ncClient, func(ctx *context.Context, nc *pb.NetworkControlClient) error {
		_, err := ncClient.CreateBridge(
			*ctx,
			&pb.CreateBridgeRequest{
				BridgeName:   bridgeName,
				DatapathType: "system",
			},
		)
		return err
	})

	Test("Dump Ports", &ctx, &ncClient, func(ctx *context.Context, nc *pb.NetworkControlClient) error {
		ports, err := ncClient.DumpPorts(
			*ctx,
			&pb.DumpPortsRequest{
				BridgeName: bridgeName,
			},
		)
		if ports.Ports[0].ID != -1 {
			return fmt.Errorf("The port ID should be -1 for LOCAL port")
		}
		if ports.Ports[0].Name != bridgeName {
			return fmt.Errorf("The port name should same as bridge Name")
		}

		return err
	})
	Test("Delete Bridge", &ctx, &ncClient, func(ctx *context.Context, nc *pb.NetworkControlClient) error {
		_, err := ncClient.DeleteBridge(
			*ctx,
			&pb.DeleteBridgeRequest{
				BridgeName: bridgeName,
			},
		)
		return err
	})
}
