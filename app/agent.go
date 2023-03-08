package app

import (
	"encoding/binary"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"streak/app/filter"
	"strings"
)

var currentIpNet *net.IPNet
var currentBroadcast net.IP

func RunAgent() {
	networkInterface := viper.GetString("network.interface")
	currentIpNet = getCurrentIpNet(networkInterface)
	currentBroadcast = getBroadcast(currentIpNet)

	go runAgent(networkInterface)
	runAgent(getLoopBackInterface())
}

func runAgent(networkInterface string) {
	if networkInterface == "" {
		panic("network interface unknown")
	}
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
			dstIp := net.ParseIP(p.NetworkLayer().NetworkFlow().Dst().String())

			if (currentIpNet.Contains(srcIp) || srcIp.IsLoopback()) && !(currentIpNet.Contains(dstIp) || dstIp.IsLoopback()) {
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

	// drop ipv4 broadcast
	if srcIp.Equal(net.IPv4bcast) || dstIp.Equal(net.IPv4bcast) {
		return true
	}

	// drop net broadcast
	if srcIp.Equal(currentBroadcast) || dstIp.Equal(currentBroadcast) {
		return true
	}

	return false
}

func getLoopBackInterface() string {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		flags := inter.Flags.String()
		if strings.Contains(flags, "up") && strings.Contains(flags, "loopback") {
			return inter.Name
		}
	}
	return ""
}

func getBroadcast(n *net.IPNet) net.IP {
	if n.IP.To4() == nil {
		return net.IP{}
	}
	ip := make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return ip
}
