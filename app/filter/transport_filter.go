package filter

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"streak/app/cache"
	"streak/app/common"
)

func TransportFilter(packet gopacket.Packet) {
	ipSrc := net.ParseIP(packet.NetworkLayer().NetworkFlow().Src().String())
	ipDst := net.ParseIP(packet.NetworkLayer().NetworkFlow().Dst().String())

	portDst := packet.TransportLayer().TransportFlow().Dst()

	proto := "unknown"
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		proto = "tcp"
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		proto = "udp"
	}

	detail := proto + "://" + ipDst.String() + ":" + portDst.String()

	domain := cache.QueryDomain(ipDst)
	if domain != "" {
		common.ReportTransport(domain, ipSrc.String(), ipDst.String(), detail)
	}else{
		common.ReportTransport("unknown", ipSrc.String(), ipDst.String(), detail)
	}
}
