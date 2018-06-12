package main

import (
	"log"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/network-controller/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type podOptions struct {
	Name string `long:"pName" description:"The Pod Name, can set by environement variable" env:"POD_NAME" required:"true"`
	NS   string `long:"pNS" description:"The namespace of the Pod, can set by environement variable" env:"POD_NAMESPACE" required:"true"`
	UUID string `long:"pUUID" description:"The UUID of the Pod, can set by environement variable" env:"POD_UUID" required:"true"`
}

type interfaceOptions struct {
	IP      string `short:"i" long:"ip" description:"The ip address of the interface, should be CIDR form"`
	Gateway string `short:"g" long:"gw" description:"The gateway of the inteface subnet"`
	VLAN    *int   `short:"v" long:"vlan" description:"The Vlan Tag of the interface"`
}

type connectOptions struct {
	Bridge    string `short:"b" long:"bridge" description:"Target bridge name" required:"true"`
	Interface string `short:"n" long:"nic" description:"The interface name in the container" required:"true"`
}

type clientOptions struct {
	Server    string           `short:"s" long:"server " description:"target server address, [ip:port] for TCP or unix://[path] for UNIX" required:"true"`
	Connect   connectOptions   `group:"ConnectOptions"`
	Interface interfaceOptions `group:"InterfaceOptions" `
	Pod       podOptions       `group:"InterfaceOptions" `
}

var options clientOptions
var parser = flags.NewParser(&options, flags.Default)

func main() {
	//flag.Parse()
	if _, err := parser.Parse(); err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	// Verify IP address
	if utils.IsValidCIDR(options.Interface.IP) {
		log.Fatalf("IP address is not correct: %s", options.Interface.IP)
	}

	// Verify gateway address
	if utils.IsValidIP(options.Interface.Gateway) {
		log.Fatalf("Gateway address is not correct: %s", options.Interface.Gateway)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(options.Server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewNetworkControlClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println(options.Pod.Name, options.Pod.NS, options.Pod.UUID)
	// Find Network Namespace Path
	n, err := c.FindNetworkNamespacePath(ctx, &pb.FindNetworkNamespacePathRequest{
		PodName:   options.Pod.Name,
		Namespace: options.Pod.NS,
		PodUUID:   options.Pod.UUID})
	if err != nil {
		log.Fatalf("There is something wrong with find network namespace pathpart.\n %v", err)
	}

	if !n.Success {
		log.Fatalf("Find network namespace path is fail. The reason is %s.", n.Reason)
	}

	log.Printf("The path is %s.", n.Path)
	// Let's connect bridge
	b, err := c.ConnectBridge(ctx, &pb.ConnectBridgeRequest{
		Path:              n.Path,
		PodUUID:           options.Pod.UUID,
		ContainerVethName: options.Connect.Interface,
		BridgeName:        options.Connect.Bridge})

	if err != nil {
		log.Fatalf("There is something wrong with connect bridge: %v", err)
	}
	if b.Success {
		log.Printf("Connecting bridge is sussessful. The reason is %s.", b.Reason)
	} else {
		log.Printf("Connecting bridge is not sussessful. The reason is %s.", b.Reason)
	}
}
