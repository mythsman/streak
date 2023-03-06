package app

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"streak/app/filter"
)

var currentIpNet *net.IPNet

func RunAgent() {
	networkInterface := viper.GetString("network.interface")
	currentIpNet = getCurrentIpNet(networkInterface)

	handle, err := pcap.OpenLive(networkInterface, 1024, true, -1)
	if err != nil {
		logrus.Fatalln(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if ignorePacket(packet) {
			continue
		}
		go func(p gopacket.Packet) {
			filter.DnsFilter(p)
			filter.TlsFilter(p)
			filter.HttpFilter(p)

			// ignore return packet
			srcIp := net.ParseIP(p.NetworkLayer().NetworkFlow().Src().String())
			if currentIpNet.Contains(srcIp) {
				filter.TransportFilter(p)
			}

		}(packet)
	}
}

func getCurrentIpNet(networkInterface string) *net.IPNet {
	iface, _ := net.InterfaceByName(networkInterface)
	addrs, _ := iface.Addrs()
	for _, addr := range addrs {
		ipNet := addr.(*net.IPNet)
		ipv4 := ipNet.IP.To4()
		if ipv4 != nil {
			return ipNet
		}
	}
	panic("No ipv4 found for " + networkInterface)
}

func ignorePacket(packet gopacket.Packet) bool {

	// drop not transport
	if packet.TransportLayer() == nil {
		return true
	}

	srcIp := net.ParseIP(packet.NetworkLayer().NetworkFlow().Src().String())
	dstIp := net.ParseIP(packet.NetworkLayer().NetworkFlow().Dst().String())

	// drop ipv6
	if srcIp.To4() == nil || dstIp.To4() == nil {
		return true
	}

	// drop multicast
	if srcIp.IsMulticast() || dstIp.IsMulticast() {
		return true
	}

	// drop interface local multicast
	if srcIp.IsInterfaceLocalMulticast() || dstIp.IsInterfaceLocalMulticast() {
		return true
	}

	// drop link local multicast
	if srcIp.IsLinkLocalMulticast() || dstIp.IsLinkLocalMulticast() {
		return true
	}

	// drop unspecified
	if srcIp.IsUnspecified() || dstIp.IsUnspecified() {
		return true
	}

	return false
}
