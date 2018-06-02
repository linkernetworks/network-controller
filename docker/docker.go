package docker

import (
	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"golang.org/x/net/context"
)

// docker ps -a <Container ID>
func ListContainer() ([]types.Container, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return cli.ContainerList(context.Background(), types.ContainerListOptions{})
}

// docker inspect <Container ID>
func InspectContainer(containerID string) (types.ContainerJSON, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return cli.ContainerInspect(context.Background(), containerID)
}

// docker inspect <Container ID> | grep -E 'SandboxKey|Id'
func GetSandboxKey(containerInfo types.ContainerJSON) string {
	return containerInfo.NetworkSettings.SandboxKey
}
