package main

import (
	"github.com/jessevdk/go-flags"
	pb "github.com/linkernetworks/network-controller/messages"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
	"time"
)

type InterfaceOptions struct {
	IP      string `short:"i" long:"ip" description:"The ip address of the interface, should be CIDR form"`
	Gateway string `short:"g" long:"gw" description:"The gateway of the inteface subnet"`
	VLAN    *int   `short:"v" long:"vlan" description:"The Vlan Tag of the interface"`
}

type ConnectOptions struct {
	Bridge    string `short:"b" long:"bridge" description:"Target bridge name" required:"true"`
	Interface string `short:"n" long:"nic" description:"The interface name in the container" required:"true"`
}

type ClientOptions struct {
	Server    string           `short:"s" long:"server " description:"target server address, [ip:port] for TCP or unix://[path] for UNIX" required:"true"`
	Connect   ConnectOptions   `group:"ConnectOptions"`
	Interface InterfaceOptions `group:"InterfaceOptions" `
}

var options ClientOptions

var parser = flags.NewParser(&options, flags.Default)

/*
!
-> -b: bridgeName
-> -n: interface name in the container
-> -v: vlanTag
-> -i: ip subnet
-> -g: gateway address of the routing rules

*/

func main() {
	//flag.Parse()
	if _, err := parser.Parse(); err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
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

	// get env value, include pod name, namespace, pod uuid, pod veth Name, bridge name
	// These are env values test case
	// os.Setenv("MY_POD_NAME", "pod1")
	// os.Setenv("MY_POD_NAMESPACE", "ns1")
	// os.Setenv("MY_POD_UUID", "1111")
	// os.Setenv("MY_POD_VETH_NAME", "veth1")
	// os.Setenv("MY_POD_BRIDGE_NAME", "br2")
	var pod_name, pod_namespace, pod_uuid, pod_vethname, pod_bridgename string
	pod_name = os.Getenv("MY_POD_NAME")
	pod_namespace = os.Getenv("MY_POD_NAMESPACE")
	pod_uuid = os.Getenv("MY_POD_UUID")
	pod_vethname = os.Getenv("MY_POD_VETH_NAME")
	pod_veth_list := strings.Split(pod_vethname, ",")
	pod_bridgename = os.Getenv("MY_POD_BRIDGE_NAME")

	if pod_name == "" || pod_namespace == "" || pod_uuid == "" || pod_vethname == "" || pod_bridgename == "" {
		log.Fatalf("The environment variables setup fault.")
	}
	log.Println(pod_name, pod_namespace, pod_uuid)
	// Find Network Namespace Path
	n, err := c.FindNetworkNamespacePath(ctx, &pb.FindNetworkNamespacePathRequest{PodName: pod_name, Namespace: pod_namespace, PodUUID: pod_uuid})
	if err != nil {
		log.Fatalf("There is something wrong with find network namespace pathpart.\n %v", err)
	}
	if n.Success {
		log.Printf("The path is %s.", n.Path)
		// Let's connect bridge
		for i := range pod_veth_list {
			log.Println(pod_veth_list[i])
			b, err := c.ConnectBridge(ctx, &pb.ConnectBridgeRequest{Path: n.Path, PodUUID: pod_uuid, ContainerVethName: string(pod_veth_list[i]), BridgeName: pod_bridgename})
			if err != nil {
				log.Fatalf("There is something wrong with connect bridge: %v", err)
			}
			if b.Success {
				log.Printf("Connecting bridge is sussessful. The reason is %s.", b.Reason)
			} else {
				log.Printf("Connecting bridge is not sussessful. The reason is %s.", b.Reason)
			}
		}
	} else {
		log.Printf("It's not success. The reason is %s.", n.Reason)
	}
}
