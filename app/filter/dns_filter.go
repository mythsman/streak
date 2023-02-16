package filter

import (
	"github.com/dgraph-io/ristretto"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
	"streak/app/common"
	"strings"
	"time"
)

var ip2Name *ristretto.Cache

func init() {
	ip2Name, _ = ristretto.NewCache(&ristretto.Config{
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

			if answer.Type == layers.DNSTypeA {
				logrus.Infoln("domain set", answer.IP.String(), common.GetShortDomain(domain))
				ip2Name.SetWithTTL(answer.IP.String(), common.GetShortDomain(domain), 1, 1*time.Hour)
			}

			common.ReportDns(domain, common.GetShortDomain(domain), answer.Type.String(), ipDst.String(), ipSrc.String())
		}
	}
}

func QueryDomain(ip string) string {
	domain, found := ip2Name.Get(ip)
	if found {
		logrus.Infoln("domain hit", ip, domain)
		return domain.(string)
	} else {
		logrus.Infoln("domain miss", ip)
		return ""
	}
}
