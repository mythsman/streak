package filter

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"streak/app/cache"
	"streak/app/common"
	"strings"
)

func DnsFilter(packet gopacket.Packet) {
	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.DNS)

		for _, answer := range dns.Answers {
			ipSrc := packet.NetworkLayer().NetworkFlow().Src()

			domain := strings.ToLower(string(answer.Name))

			if answer.Type == layers.DNSTypeA {
				cache.SetDomain(answer.IP.String(), domain)
			}

			common.ReportDns(common.GetShortDomain(domain), ipSrc.String(), domain)
		}
	}
}
