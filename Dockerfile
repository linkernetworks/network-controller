FROM golang:1.10-alpine3.7

WORKDIR /go/src/github.com/linkernetworks/network-controller

RUN apk add --no-cache protobuf ca-certificates git

COPY ./ /go/src/github.com/linkernetworks/network-controller
RUN go install ./server/... ./client/...

FROM alpine:3.7
RUN apk add --no-cache ca-certificates rsync
WORKDIR /net-controller

COPY --from=0 /go/bin /go/bin
