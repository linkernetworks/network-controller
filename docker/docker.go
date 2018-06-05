package docker

import (
	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"golang.org/x/net/context"
)

// docker ps -a <Container ID>
func ListContainer(cli *client.Client) ([]types.Container, error) {
	return cli.ContainerList(context.Background(), types.ContainerListOptions{})
}

// docker inspect <Container ID>
func InspectContainer(cli *client.Client, containerID string) (types.ContainerJSON, error) {
	return cli.ContainerInspect(context.Background(), containerID)
}

// docker inspect <Container ID> | grep -E 'SandboxKey|Id'
func GetSandboxKey(containerInfo types.ContainerJSON) string {
	return containerInfo.NetworkSettings.SandboxKey
}
