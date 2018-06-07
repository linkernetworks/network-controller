package docker

import (
	"io"
	"os"
	"testing"

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
	if _, ok := os.LookupEnv("TEST_DOCKER"); !ok {
		t.Skip()
	}
	cli, err := New()
	assert.NoError(t, err)

	image := "alpine"

	//Current Number of container
	currentLength := 0
	c, err := cli.ListContainer()
	assert.NoError(t, err)
	currentLength = len(c)

	//Image Pull
	r, err := cli.Client.ImagePull(cli.Context, image, types.ImagePullOptions{})
	assert.NoError(t, err)
	io.Copy(os.Stdout, r)

	//Create Container
	resp, err := cli.Client.ContainerCreate(cli.Context, &container.Config{
		Image: image,
		Cmd:   []string{"tail", "-f", "/etc/hosts"},
		Tty:   false,
	}, nil, nil, "")
	defer cli.Client.ContainerRemove(cli.Context, resp.ID, types.ContainerRemoveOptions{Force: true})
	assert.NoError(t, err)

	//Start the container
	err = cli.Client.ContainerStart(cli.Context, resp.ID, types.ContainerStartOptions{})
	assert.NoError(t, err)

	//Inspect
	json, err := cli.InspectContainer(resp.ID)
	assert.NoError(t, err)
	assert.Equal(t, image, json.Config.Image)

	//The number of container should be original + 1
	c, err = cli.ListContainer()
	assert.NoError(t, err)
	assert.Equal(t, currentLength+1, len(c))
}
