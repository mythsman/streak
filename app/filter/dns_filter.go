package filter

import (
	"github.com/dgraph-io/ristretto"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"streak/app/common"
	"strings"
)

// https://en.wikipedia.org/wiki/List_of_Internet_top-level_domains
var specialDomains map[string]bool

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

	specialDomains = make(map[string]bool)

	for _, name := range [...]string{"com", "org", "net", "int", "edu", "gov", "mil"} {
		specialDomains[name] = true
	}

}

func DnsFilter(packet gopacket.Packet) {
	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.DNS)
		for _, answer := range dns.Answers {
			ipSrc := packet.NetworkLayer().NetworkFlow().Src()
			ipDst := packet.NetworkLayer().NetworkFlow().Dst()

			domain := strings.ToLower(string(answer.Name))

			common.ReportDns(domain, getShortDomain(domain), answer.Type.String(), ipDst.String(), ipSrc.String())
		}
	}
}

func getShortDomain(domain string) string {
	split := strings.Split(domain, ".")
	if len(split) <= 2 {
		return domain
	}
	secondLevel := split[len(split)-2]
	if specialDomains[secondLevel] {
		return strings.Join(split[len(split)-3:], ".")
	} else {
		return strings.Join(split[len(split)-2:], ".")
	}
}
