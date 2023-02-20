package app

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
)

var currentIpNet *net.IPNet

func RunAgent() {
	networkInterface := viper.GetString("network.interface")
	currentIpNet = getCurrentIpNet(networkInterface)

	handle, err := pcap.OpenLive(networkInterface, 1024, false, -1)
	if err != nil {
		logrus.Fatalln(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if ignorePacket(packet) {
			continue
		}
		//go func(val gopacket.Packet) {
		//	filter.DnsFilter(val)
		//	filter.TlsFilter(val)
		//	filter.HttpFilter(val)
		//	filter.TransportFilter(val)
		//}(packet)
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

func judgeType(srcIp net.IP, dstIp net.IP) string {
	if srcIp.Equal(currentIpNet.IP) && !dstIp.IsPrivate() {
		return "tx_public"
	}

	if !srcIp.IsPrivate() && dstIp.Equal(currentIpNet.IP) {
		return "rx_public"
	}

	if srcIp.Equal(currentIpNet.IP) && dstIp.IsPrivate() {
		return "tx_private"
	}

	if srcIp.IsPrivate() && dstIp.Equal(currentIpNet.IP) {
		return "rx_private"
	}

	if srcIp.IsPrivate() && !dstIp.IsPrivate() {
		return "tx_route"
	}

	if !srcIp.IsPrivate() && dstIp.IsPrivate() {
		return "rx_route"
	}

	if srcIp.IsPrivate() && dstIp.IsPrivate() {
		return "both private"
	}

	if !srcIp.IsPrivate() && !dstIp.IsPrivate() {
		return "both public?"
	}

	panic("packet type unknown")
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

	// drop unspecified
	if srcIp.IsUnspecified() || dstIp.IsUnspecified() {
		return true
	}

	packetType := judgeType(srcIp, dstIp)
	logrus.Infoln("packet type :", packetType)

	return false
}
