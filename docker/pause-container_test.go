package docker

import (
	"testing"

	"docker.io/go-docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestFindK8SPauseContainerID(t *testing.T) {
	containers := []types.Container{
		types.Container{
			ID:    "18720aa37e71",
			Names: []string{"k8s_POD_mongo-0_default_553be55e-532c-11e8-bea4-42010af0009f_0"},
		},
		types.Container{
			ID:    "12332aa37f99",
			Names: []string{"k8s_POD_redis_default_667be33e-178c-99e2-bea6-68090bd0001d_20"},
		},
	}
	cid, err := FindK8SPauseContainerID(containers, "mongo-0", "default", "553be55e-532c-11e8-bea4-42010af0009f")
	assert.NoError(t, err)
	assert.Equal(t, "18720aa37e71", cid, "they should be equal")
}
