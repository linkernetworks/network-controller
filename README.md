# network-controller [![Build Status](https://travis-ci.org/linkernetworks/network-controller.svg?branch=master)](https://travis-ci.org/linkernetworks/network-controller) [![codecov](https://codecov.io/gh/linkernetworks/network-controller/branch/master/graph/badge.svg)](https://codecov.io/gh/linkernetworks/network-controller) [![Go Report Card](https://goreportcard.com/badge/github.com/linkernetworks/network-controller)](https://goreportcard.com/report/github.com/linkernetworks/network-controller)  [![Docker Build Status](https://img.shields.io/docker/build/sdnvortex/network-controller.svg)](https://hub.docker.com/r/sdnvortex/network-controller/)


## Development

```shell
# generate protocol buffer
make pb

# make server binary
make server

# make client binary
make client

# make test (You should run this before push codes)
make test
```

### Develop with Vagrant
1. Download repository
```sh
$ git clone https://github.com/linkernetworks/network-controller.git
$ cd network-controller/
```
2. Run vagrant
```sh
$ vagrant up
```
3. Install vagrant-scp plugin
```sh
$ vagrant plugin install vagrant-scp
```
4. Get `~/.kube/config` from particular VM to Host Os
```sh
$ vagrant scp :~/.kube/config vagrant-config
```
5. Configure Access to vagrant Cluster
> Please install `kubectl` before you begin.
To be continue...

## Usage

### Run a Server
The network-controller server provide two ways to listen, TCP and Unix domain socket
If you want to run as a UNIX domain socket server, you should specify a path to store the sock file
and server will remove that file when server is been terminated
```shell
./server/network-controller-server -unix=/tmp/a.sock
```
For the TCP server, just use the `IP:PORT` format, for example, `0.0.0.0:50051`
```shell
./server/network-controller-server -tcp=0.0.0.0:50051
```

### Run a Client
The clinet is used for the kubernetes pod to create the veth and you can see the example yaml in `example/kubernetes/*.yaml` to see how to use this client.

For creating a veth for Pod, the client needs the following information
- Pod Name
- Pod Namespace
- Pod UUID
- Target Bridge Name
- The Interface Name in the container
- The server address

The first three variable can passed by the environemtn `POD_NAME`, `POD_NAMESPACE` and `POD_UUID`.

#### Bridge Name
`-b` or `--bridge`

#### Interface Name
`-n` or `--nic`

#### Server
The clinet support two way to connect to the server, TCP socket and UNIX domain socket.
In the TCP mode, use the `IP:PORT` format to connect to TCP server
```shell
./client/network-controller-client -server=0.0.0.0:50051
```
Fot the UNIX domain socket mode, you should use the `unix://PATH` format to connect to server.
Assume the path is `/tmp/a.sock` and you can use the following command to connect
```shell
./client/network-controller-client -server=unix:///tmp/a.sock

Following is an example of the client and you can see more use the `--help`.
```shell
./clinet/network-controller-client "-s", "unix:///tmp/grpc.sock", "-b", "br100", "-n", "eth100"
```
