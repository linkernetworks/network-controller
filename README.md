# network-controller [![Build Status](https://travis-ci.org/linkernetworks/network-controller.svg?branch=master)](https://travis-ci.org/linkernetworks/network-controller) [![codecov](https://codecov.io/gh/linkernetworks/network-controller/branch/master/graph/badge.svg)](https://codecov.io/gh/linkernetworks/network-controller) [![Go Report Card](https://goreportcard.com/badge/github.com/linkernetworks/network-controller)](https://goreportcard.com/report/github.com/linkernetworks/network-controller)

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
The clinet support two way to connect to the server, TCP socket and UNIX domain socket.
In the TCP mode, use the `IP:PORT` format to connect to TCP server
```shell
./client/network-controller-client -server=0.0.0.0:50051
```
Fot the UNIX domain socket mode, you should use the `unix://PATH` format to connect to server.
Assume the path is `/tmp/a.sock` and you can use the following command to connect
```shell
./client/network-controller-client -server=unix:///tmp/a.sock
