package docker

import (
	"context"
	"io"
	"os"
	"testing"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"

	"github.com/stretchr/testify/assert"
)

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

func TestListContainer(t *testing.T) {
	cli, err := client.NewEnvClient()
	assert.NoError(t, err)

	image := "alpine"
	currentLength := 0
	c, err := ListContainer(cli)
	assert.NoError(t, err)
	currentLength = len(c)

	ctx := context.Background()
	r, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	assert.NoError(t, err)
	io.Copy(os.Stdout, r)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"tail", "-f", "/etc/hosts"},
		Tty:   false,
	}, nil, nil, "")
	assert.NoError(t, err)

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	assert.NoError(t, err)

	json, err := InspectContainer(cli, resp.ID)
	assert.NoError(t, err)
	assert.Equal(t, image, json.Config.Image)

	c, err = ListContainer(cli)
	assert.NoError(t, err)
	assert.Equal(t, currentLength+1, len(c))

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true})
	assert.NoError(t, err)
}
