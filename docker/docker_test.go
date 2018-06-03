package docker

import (
	"testing"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func TestListContainer(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	assert.NoError(t, err)

	imageName := "hello-world"

	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	assert.NoError(t, err)

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
		}, nil, nil, "")
	assert.NoError(t, err)

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	assert.NoError(t, err)

	var imgNameSlice []string
	containers, err := ListContainer()
	for _, container := range containers {
		imgNameSlice = append(imgNameSlice, container.Image)
	}
	assert.Contains(t, imgNameSlice, "hello-world")
}

func TestGetSandboxKey(t *testing.T) {
	containerInfo := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID: "43add2074a1bd05e6ba1e661de5b9ed109a058eabdb83032e4801389c78b035a",
		},
		NetworkSettings: &types.NetworkSettings{
			NetworkSettingsBase: types.NetworkSettingsBase{
				Bridge:                 "",
				SandboxID:              "66d0d6831275fd6102f22105cde6c7442bbde4974343c74adedd5b1650e1443d",
				HairpinMode:            false,
				LinkLocalIPv6Address:   "",
				LinkLocalIPv6PrefixLen: 0,
				SandboxKey:             "/var/run/docker/netns/66d0d6831275",
				SecondaryIPAddresses:   nil,
				SecondaryIPv6Addresses: nil,
			},
		},
	}
	netnsPath := GetSandboxKey(containerInfo)
	assert.Equal(t, "/var/run/docker/netns/66d0d6831275", netnsPath)
}
