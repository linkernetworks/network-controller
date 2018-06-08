package netlinkEvent

import (
	ovs "github.com/linkernetworks/network-controller/openvswitch"
	"github.com/vishvananda/netlink"
)

func RemoveVethFromOVS(lu netlink.LinkUpdate) error {
	L := lu.Attrs()
	if L.ParentIndex != 0 || L.MasterIndex != 0 {
		return nil
	}

	o = ovs.New()

	bridges, err := o.ListBridges()
	if err != nil {
		return err
	}
	for _, v := range bridges {
		ports, _ := o.ListPorts(v)
		for _, port := range ports {
			if port == L.Name {
				o.DeletePort(v, L.Name)
				return nil
			}
		}
	}

	return nil
}
