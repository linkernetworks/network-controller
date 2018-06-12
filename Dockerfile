FROM golang:1.10-alpine3.7

WORKDIR /go/src/github.com/linkernetworks/network-controller

RUN apk add --no-cache protobuf ca-certificates git make

COPY ./ /go/src/github.com/linkernetworks/network-controller
RUN go get -u github.com/golang/protobuf/proto &&\
    go get -u github.com/golang/protobuf/protoc-gen-go &&\
    go get -u github.com/kardianos/govendor
RUN govendor sync &&\
    make pb
RUN go install ./server/... ./client/...

FROM alpine:3.7
RUN apk add --no-cache ca-certificates openvswitch sudo
WORKDIR /go/bin

COPY --from=0 /go/bin/server /go/bin
COPY --from=0 /go/bin/client /go/bin
