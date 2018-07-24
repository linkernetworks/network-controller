package nl

import (
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"log"
)

// NlEventHandler is the structure for netlink event handler
type NlEventHandler struct {
	stop              chan struct{}
	LinkDeleteHandler []linkReceiver
}

type linkReceiver func(lu netlink.LinkUpdate) bool

// New will return an new netlink event handler
func New() *NlEventHandler {
	return &NlEventHandler{
		stop: make(chan struct{}),
	}
}

// AddDeletedLinkHandler will add DeletedLinkHandler
func (nl *NlEventHandler) AddDeletedLinkHandler(handler linkReceiver) {
	nl.LinkDeleteHandler = append(nl.LinkDeleteHandler, handler)
}

// TrackNetlink will track netlink
func (nl *NlEventHandler) TrackNetlink() error {

	stop := make(chan struct{})
	data := make(chan netlink.LinkUpdate)
	if err := netlink.LinkSubscribe(data, stop); err != nil {
		return err
	}

	for {
		select {
		case c := <-data:
			switch c.Header.Type {
			case unix.RTM_DELLINK:
				for _, v := range nl.LinkDeleteHandler {
					if v(c) == true {
						nl.Stop()
					}
				}

			}

		case <-nl.stop:
			log.Println("Receive Stop")
			var e struct{}
			stop <- e
			return nil
		}
	}
}

// Stop will stop netlink
func (nl *NlEventHandler) Stop() {
	if nl.stop == nil {
		return
	}
	close(nl.stop)
	nl.stop = nil
}
