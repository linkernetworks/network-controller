#!/bin/bash
protoc ./messages/messages.proto --go_out=plugins=grpc:.
