package app

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"streak/app/filter"
)

func Run() {
	InitLogger()

	InitConfig()

	InitInfluxdb()

	loop()
}

var dnsFilter = &filter.DnsFilter{}
var tlsFilter = &filter.TlsFilter{}
var httpFilter = &filter.HttpFilter{}
var transportFilter = &filter.TransportFilter{}

func loop() {

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
			dnsFilter.Filter(packet)
			tlsFilter.Filter(packet)
			httpFilter.Filter(packet)
			transportFilter.Filter(packet)
		}
	}
}
