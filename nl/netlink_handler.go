package nl

import (
	ovs "github.com/linkernetworks/network-controller/openvswitch"
	"github.com/vishvananda/netlink"
	"log"
)

/*
	Remove the veth from the system.
	We get the veth name from the netlink.LinkUpdate
	and traverse all OpenvSwitches and try to remove the veth from its parent OVS

	return true will stop the netlink server handler
*/
func RemoveVethFromOVS(lu netlink.LinkUpdate) bool {
	L := lu.Attrs()
	if L.ParentIndex != 0 || L.MasterIndex != 0 {
		return false
	}

	o := ovs.New()

	bridges, err := o.ListBridges()
	if err != nil {
		return true
	}
	for _, v := range bridges {
		ports, _ := o.ListPorts(v)
		for _, port := range ports {
			if port == L.Name {
				log.Printf("Try to remove %s from %s\n", L.Name, v)
				o.DeletePort(v, L.Name)
				return false
			}
		}
	}

	return false
}
