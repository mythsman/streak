package filter

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"streak/app/cache"
	"streak/app/common"
)

func TransportFilter(packet gopacket.Packet) {
	ipSrc := net.ParseIP(packet.NetworkLayer().NetworkFlow().Src().String()).String()
	ipDst := net.ParseIP(packet.NetworkLayer().NetworkFlow().Dst().String()).String()

	portDst := packet.TransportLayer().TransportFlow().Dst()

	proto := "unknown"
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		proto = "tcp"
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		proto = "udp"
	}

	detail := proto + "://" + ipDst + ":" + portDst.String()

	domain := cache.QueryDomain(ipDst)
	if domain != "" {
		common.ReportTransport(domain, ipSrc, ipDst, detail)
	}
}
