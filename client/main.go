package main

import (
	"flag"
	"log"
	"os"
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
	if n.Success == true {
		log.Printf("The path is %s.", n.Path)
		// Let's connect bridge
		b, err := c.ConnectBridge(ctx, &pb.ConnectBridgeRequest{Path: n.Path, PodUUID: pod_uuid, ContainerVethName: pod_vethname, BridgeName: pod_bridgename})
		if err != nil {
			log.Fatalf("There is something wrong with connect bridge: %v", err)
		}
		if b.Success == true {
			log.Printf("Connecting bridge is sussessful. The reason is %s.", b.Reason)
		} else {
			log.Printf("Connecting bridge is not sussessful. The reason is %s.", b.Reason)
		}
	} else {
		log.Printf("It's not success. The reason is %s.", n.Reason)
	}

	// ping
	r, err := c.Ping(ctx, &pb.PingRequest{Ping: "PING"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Got: %s", r.Pong)
}
