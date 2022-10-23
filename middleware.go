package main

import (
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/network"
)

func CheckPacketSecurity(handler fuse.Handler) fuse.Handler {
	return func(packet *network.Packet, p fuse.Principal) {
		//todo check packet security
		handler(packet, p)

	}
}
