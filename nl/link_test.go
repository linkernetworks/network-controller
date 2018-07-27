package nl

import (
	"net"
	"os"
	"testing"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/containernetworking/plugins/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/vishvananda/netlink"
)

func TestMakeVethPair(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_VETH"); !ok {
		t.SkipNow()
	}
	contVeth, err := makeVethPair("test-veth-ns1", "test-ovs-1", 1500)
	defer netlink.LinkDel(contVeth)
	assert.NoError(t, err)

	assert.Equal(t, "test-veth-ns1", contVeth.Attrs().Name)
	assert.Equal(t, 1500, contVeth.Attrs().MTU)
}

func TestInvalidmakeVethPair(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_VETH"); !ok {
		t.SkipNow()
	}
	contVeth, err := makeVethPair("invalid-test-veth-ns1", "invalid-test-ovs-1", -1800)
	// Err: numerical result out of range
	assert.Error(t, err)
	assert.Nil(t, contVeth)
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
	assert.NoError(t, err)

	hostVeth, err := netlink.LinkByName(hostVethName)
	assert.NoError(t, err)
	defer netlink.LinkDel(hostVeth)
}

func TestInvalidSetupVeth(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_VETH"); !ok {
		t.SkipNow()
	}
	contIfName := "invalid-test-net0"
	hostVethName := "invalid-test-ovs-0"

	// Create a network namespace
	netns, err := testutils.NewNS()
	assert.NoError(t, err)

	err = netns.Do(func(hostNS ns.NetNS) error {
		// create the veth pair in the container and move host end into host netns
		_, _, err := SetupVeth(contIfName, hostVethName, -1500, hostNS)
		assert.Error(t, err)
		return nil
	})
	assert.NoError(t, err)
}

func TestAddRoute(t *testing.T) {
	type args struct {
		ipn  *net.IPNet
		gwIP string
		dev  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test lo interface with null GW",
			args: args{
				ipn: &net.IPNet{
					IP:   net.IPv4(224, 0, 0, 0),
					Mask: net.CIDRMask(4, 32),
				},
				gwIP: "",
				dev:  "lo",
			},
			wantErr: true,
		}, {
			name: "Test lo interface with 0.0.0.0",
			args: args{
				ipn: &net.IPNet{
					IP:   net.IPv4(224, 0, 0, 0),
					Mask: net.CIDRMask(4, 32),
				},
				gwIP: "0.0.0.0",
				dev:  "lo",
			},
			wantErr: true,
		}, {
			name: "Test unknow interface",
			args: args{
				ipn: &net.IPNet{
					IP:   net.IPv4(192, 168, 0, 0),
					Mask: net.CIDRMask(24, 32),
				},
				gwIP: "0.0.0.0",
				dev:  "unknow",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a network namespace
			netns, err := testutils.NewNS()
			assert.NoError(t, err)
			err = netns.Do(func(hostNS ns.NetNS) error {
				if err := AddRoute(tt.args.ipn, tt.args.gwIP, tt.args.dev); (err != nil) != tt.wantErr {
					t.Errorf("AddRoute() error = %v, wantErr %v", err, tt.wantErr)
				}
				return nil
			})
			assert.NoError(t, err)
		})
	}
}
