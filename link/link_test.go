package link

import (
	"os"
	"testing"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/containernetworking/plugins/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/vishvananda/netlink"
)

func TestmakeVethPair(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_VETH"); !ok {
		t.SkipNow()
	}
	contVeth, err := makeVethPair("test-veth-ns1", "test-ovs-1", 1500)
	assert.NoError(t, err)

	assert.Equal(t, "test-veth-ns1", contVeth.Attrs().Name)
	assert.Equal(t, 1500, contVeth.Attrs().MTU)
	defer netlink.LinkDel(contVeth)
}

func TestInvalidmakeVethPair(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_VETH"); !ok {
		t.SkipNow()
	}
	contVeth, err := makeVethPair("invalid-test-veth-ns1", "invalid-test-ovs-1", -1800)
	// Err: numerical result out of range
	assert.Error(t, err)
	defer netlink.LinkDel(contVeth)
}

func TestSetupVeth(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_VETH"); !ok {
		t.SkipNow()
	}
	contIfName := "test-net0"
	hostVethName := "test-ovs-0"

	// Create a network namespace
	netns, err := testutils.NewNS()
	assert.NoError(t, err)

	err = netns.Do(func(hostNS ns.NetNS) error {
		// create the veth pair in the container and move host end into host netns
		hostVeth, containerVeth, err := SetupVeth(contIfName, hostVethName, 1500, hostNS)
		assert.NoError(t, err)

		assert.Equal(t, "test-net0", containerVeth.Name)
		assert.Equal(t, "test-ovs-0", hostVeth.Name)
		return nil
	})
	hostVeth, err := netlink.LinkByName(hostVethName)
	assert.NoError(t, err)
	defer netlink.LinkDel(hostVeth)
}
