package openvswitch

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ovsFile = "/usr/bin/ovs-vsctl"

var o *OVSManager

func TestMain(m *testing.M) {
	if _, ok := os.LookupEnv("TEST_OVS"); ok {
		o = New()
		retCode := m.Run()
		os.Exit(retCode)
	}
}

func changeVSCtl(t *testing.T) os.FileMode {
	info, err := os.Stat(ovsFile)
	assert.NoError(t, err)

	err = os.Chmod(ovsFile, 0666)
	assert.NoError(t, err)
	return info.Mode()
}

func resetVSCtl(mode os.FileMode) {
	os.Chmod(ovsFile, mode)
}

func TestBridgeOperations(t *testing.T) {
	bridges, err := o.ListBridges()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(bridges))

	bridgeName := "br0"
	const dpTypeSystem = "system"
	err = o.CreateBridge(bridgeName, dpTypeSystem)
	assert.NoError(t, err)

	bridges, err = o.ListBridges()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bridges))

	err = o.DeleteBridge(bridgeName)
	assert.NoError(t, err)

	bridges, err = o.ListBridges()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(bridges))
}

func TestBridgeOperationsFail(t *testing.T) {
	mode := changeVSCtl(t)
	defer resetVSCtl(mode)

	_, err := o.ListBridges()
	assert.Error(t, err)

	bridgeName := "br0"
	const dpTypeSystem = "system"
	err = o.CreateBridge(bridgeName, dpTypeSystem)
	assert.Error(t, err)

	_, err = o.ListBridges()
	assert.Error(t, err)

	err = o.DeleteBridge(bridgeName)
	assert.Error(t, err)
}

func TestAddDelPort(t *testing.T) {
	bridgeName := "br0"
	const dpTypeSystem = "system"
	err := o.CreateBridge(bridgeName, dpTypeSystem)
	defer o.DeleteBridge(bridgeName)

	hName := "test0"
	cName := "test0_peer"
	err = exec.Command("ip", "link", "add", hName, "type", "veth", "peer", "name", cName).Run()
	assert.NoError(t, err)
	defer exec.Command("ip", "link", "del", hName).Output()
	err = o.AddPort(bridgeName, hName)
	assert.NoError(t, err)

	ports, err := o.ListPorts(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ports))

	err = o.DeletePort(bridgeName, hName)
	assert.NoError(t, err)

	ports, err = o.ListPorts(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(ports))
}

func TestListPortsFail(t *testing.T) {
	bridgeName := "br0"
	ports, err := o.ListPorts(bridgeName)
	assert.Error(t, err)
	assert.Equal(t, 0, len(ports))
}

func TestAddDelPortFail(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}
	mode := changeVSCtl(t)
	defer resetVSCtl(mode)

	bridgeName := "br0"
	err := o.AddPort(bridgeName, "0")
	assert.Error(t, err)
	err = o.DeletePort(bridgeName, "0")
	assert.Error(t, err)
}

func TestFlowOperation(t *testing.T) {
	bridgeName := "br0"
	const dpTypeSystem = "system"
	err := o.CreateBridge(bridgeName, dpTypeSystem)
	defer o.DeleteBridge(bridgeName)
	assert.NoError(t, err)

	flowString := "cookie=1, actions=NORMAL"
	err = o.AddFlow(bridgeName, flowString)
	assert.NoError(t, err)

	flows, err := o.DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(flows))

	err = o.DeleteFlow(bridgeName, flowString)
	assert.NoError(t, err)

	flows, err = o.DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(flows))
}

func TestFlowOperationsFail(t *testing.T) {
	bridgeName := "br0"
	err := o.AddFlow(bridgeName, "")
	assert.Error(t, err)
	err = o.DeleteFlow(bridgeName, "")
	assert.Error(t, err)
	flows, err := o.DumpFlows(bridgeName)
	assert.Error(t, err)
	assert.Equal(t, 0, len(flows))
}

func TestDumpPorts(t *testing.T) {
	bridgeName := "br0"
	const dpTypeSystem = "system"
	err := o.CreateBridge(bridgeName, dpTypeSystem)
	defer o.DeleteBridge(bridgeName)

	hName := "test0"
	cName := "test0_peer"
	err = exec.Command("ip", "link", "add", hName, "type", "veth", "peer", "name", cName).Run()
	assert.NoError(t, err)
	defer exec.Command("ip", "link", "del", hName).Output()
	err = o.AddPort(bridgeName, hName)
	assert.NoError(t, err)

	ports, err := o.DumpPorts(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0x0), ports[0].Received.Packets)

	port, err := o.DumpPort(bridgeName, hName)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0x0), port.Received.Packets)

	err = o.DeletePort(bridgeName, hName)
	assert.NoError(t, err)
}

func TestDumpPortsFail(t *testing.T) {
	bridgeName := "br0"
	_, err := o.DumpPorts(bridgeName)
	assert.Error(t, err)

	_, err = o.DumpPort(bridgeName, "")
	assert.Error(t, err)
}
