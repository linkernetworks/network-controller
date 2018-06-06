package docker

import (
	"fmt"
	"regexp"

	"docker.io/go-docker/api/types"
)

// FindK8SPauseContainerID : Use pod-name, k8s-namespace, pod-uuid to match container-id in containers
func FindK8SPauseContainerID(containers []types.Container, PodName, Namespace, PodUUID string) (string, error) {
	// ex: k8s_POD_myinit_default_05ab36d8-65aa-11e8-b35e-42010af00248_0
	pattern := fmt.Sprintf("k8s_POD_%s_%s_%s_\\d+", PodName, Namespace, PodUUID)
	r, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	for _, container := range containers {
		// the first container name would be fine
		// don't know why container will have more than one name
		if r.MatchString(container.Names[0]) {
			return container.ID, nil
		}
	}
	return "", nil
}
