package nl

import (
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"log"
)

type NlEventHandler struct {
	stop              chan struct{}
	LinkDeleteHandler []LinkReceiver
}

type LinkReceiver func(lu netlink.LinkUpdate) bool

func New() *NlEventHandler {
	return &NlEventHandler{
		stop: make(chan struct{}),
	}
}

func (nl *NlEventHandler) AddDeletedLinkHandler(handler LinkReceiver) {
	nl.LinkDeleteHandler = append(nl.LinkDeleteHandler, handler)
}

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

func (nl *NlEventHandler) Stop() {
	if nl.stop == nil {
		return
	}
	log.Println("Stop the netlink event tracker")
	close(nl.stop)
	nl.stop = nil
}
