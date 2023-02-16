package filter

import (
	"github.com/dgraph-io/ristretto"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"streak/app/common"
	"strings"
)

var ip2Name *ristretto.Cache
var name2Ip *ristretto.Cache

func init() {
	ip2Name, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1 << 24, // 16M
		MaxCost:     1 << 27, // 128MB
		BufferItems: 64,
	})
	name2Ip, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1 << 24, // 16M
		MaxCost:     1 << 27, // 128MB
		BufferItems: 64,
	})

}

func DnsFilter(packet gopacket.Packet) {
	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.DNS)
		for _, answer := range dns.Answers {
			ipSrc := packet.NetworkLayer().NetworkFlow().Src()
			ipDst := packet.NetworkLayer().NetworkFlow().Dst()

			domain := strings.ToLower(string(answer.Name))

			common.ReportDns(domain, common.GetShortDomain(domain), answer.Type.String(), ipDst.String(), ipSrc.String())
		}
	}
}
