package netlinkEvent

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
)

func TestNetlinkEventTrack(t *testing.T) {
	wait := make(chan struct{})
	tracker := New()
	var handler = func(lu netlink.LinkUpdate) bool {
		var e struct{}
		wait <- e
		return false
	}

	var stopHandler = func(lu netlink.LinkUpdate) bool {
		return true
	}

	tracker.AddDeletedLinkHandler(handler)
	tracker.AddDeletedLinkHandler(stopHandler)

	hName := "test0"
	cName := "test0_peer"
	err := exec.Command("ip", "link", "add", hName, "type", "veth", "peer", "name", cName).Run()
	require.NoError(t, err, "You should have the permission to create the veth to test the netlink")

	go tracker.TrackNetlink()

	exec.Command("ip", "link", "del", hName).Output()
	<-wait
	tracker.Stop()
}
