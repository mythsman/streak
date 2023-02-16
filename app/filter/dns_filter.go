package filter

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
)

func DnsFilter(packet gopacket.Packet) {
	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.DNS)
		for _, answer := range dns.Answers {
			if answer.Type == layers.DNSTypeA {
				logrus.Printf("dns A %s %s", answer.Name, answer.IP)
			} else if answer.Type == layers.DNSTypeCNAME {
				logrus.Printf("dns cname %s %s", answer.Name, answer.CNAME)
			}
		}
	}
}
