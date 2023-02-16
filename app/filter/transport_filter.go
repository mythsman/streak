package filter

import (
	"github.com/google/gopacket"
	"github.com/sirupsen/logrus"
)

func TransportFilter(packet gopacket.Packet) {
	ipSrc := packet.NetworkLayer().NetworkFlow().Src()
	ipDst := packet.NetworkLayer().NetworkFlow().Dst()

	portSrc := packet.TransportLayer().TransportFlow().Src()
	portDst := packet.TransportLayer().TransportFlow().Dst()
	logrus.Printf("transport Layer %s:%s -> %s:%s", ipSrc, portSrc, ipDst, portDst)
}
