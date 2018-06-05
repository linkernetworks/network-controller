package docker

import (
	"testing"

	"docker.io/go-docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestFindK8SPauseContainerID(t *testing.T) {
	containers := []types.Container{
		{
			ID:    "18720aa37e71",
			Names: []string{"k8s_POD_mongo-0_default_553be55e-532c-11e8-bea4-42010af0009f_0"},
		},
		{
			ID:    "12332aa37f99",
			Names: []string{"k8s_POD_redis_default_667be33e-178c-99e2-bea6-68090bd0001d_20"},
		},
	}
	cid, err := FindK8SPauseContainerID(containers, "mongo-0", "default", "553be55e-532c-11e8-bea4-42010af0009f")
	assert.NoError(t, err)
	assert.Equal(t, "18720aa37e71", cid, "they should be equal")
}

func TestFindK8SPauseContainerIDFail(t *testing.T) {
	containers := []types.Container{
		{
			ID:    "18720aa37e71",
			Names: []string{"k8s_POD_mongo-0_default_553be55e-532c-11e8-bea4-42010af0009f_0"},
		},
		{
			ID:    "12332aa37f99",
			Names: []string{"k8s_POD_redis_default_667be33e-178c-99e2-bea6-68090bd0001d_20"},
		},
	}
	//Invalid Regexp Pattern
	cid, err := FindK8SPauseContainerID(containers, "(+", "default", "553be55e-532c-11e8-bea4-42010af0009f")
	assert.Error(t, err)
	assert.Equal(t, "", cid)

	cid, err = FindK8SPauseContainerID(containers, "no-exist", "default", "553be55e-532c-11e8-bea4-42010af0009f")
	assert.NoError(t, err)
	assert.Equal(t, "", cid)
}
