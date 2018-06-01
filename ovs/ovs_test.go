package ovs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAddBridge(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	defer DeleteBridge(bridgeName)
	assert.NoError(t, err)

	bridges, err := ListBridges()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bridges))
}

func TestDeleteBridge(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	assert.NoError(t, err)

	bridges, err := ListBridges()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bridges))

	err = DeleteBridge(bridgeName)
	assert.NoError(t, err)
}

func TestListBridges(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridges, err := ListBridges()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(bridges))
}

func TestAddFlow(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	defer DeleteBridge(bridgeName)
	assert.NoError(t, err)

	flowString := "cookie=1, actions=NORMAL"
	err = AddFlow(bridgeName, flowString)
	assert.NoError(t, err)

	flows, err := DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(flows))
}

func TestDeleteFlows(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	defer DeleteBridge(bridgeName)
	assert.NoError(t, err)

	flowString := "cookie=1, actions=NORMAL"
	err = AddFlow(bridgeName, flowString)
	assert.NoError(t, err)

	err = DeleteFlows(bridgeName, flowString)
	assert.NoError(t, err)

	flows, err := DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(flows))
}

func TestDumpFlows(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	defer DeleteBridge(bridgeName)
	assert.NoError(t, err)

	flows, err := DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(flows))
}
