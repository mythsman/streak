package app

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"streak/app/filter"
)

func RunAgent() {

	networkInterface := viper.GetString("network.interface")

	handle, err := pcap.OpenLive(networkInterface, 1024, false, -1)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Only process transport layer (tcp , udp)
		if packet.TransportLayer() != nil {
			filter.DnsFilter(packet)
			filter.TlsFilter(packet)
			filter.HttpFilter(packet)
			filter.TransportFilter(packet)
		}
	}
}
