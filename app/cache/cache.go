package cache

import (
	"github.com/dgraph-io/ristretto"
	"github.com/sirupsen/logrus"
	"streak/app/common"
	"time"
)

var rDnsCache *ristretto.Cache

var sniCache *ristretto.Cache

func init() {
	rDnsCache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1 << 24, // 16M
		MaxCost:     1 << 27, // 128MB
		BufferItems: 64,
	})
	sniCache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1 << 24, // 16M
		MaxCost:     1 << 27, // 128MB
		BufferItems: 64,
	})

}

func SetSni(domain string, ip1 string, port1 string, ip2 string, port2 string) {
	shortDomain := common.GetShortDomain(domain)
	sniKey := makeSniKey(ip1, port1, ip2, port2)
	sniCache.SetWithTTL(sniKey, shortDomain, 1, 1*time.Hour)
	logrus.Infoln("sni set", sniKey, shortDomain)
}

func QuerySni(ip1 string, port1 string, ip2 string, port2 string) string {
	sniKey := makeSniKey(ip1, port1, ip2, port2)
	domain, found := sniCache.Get(sniKey)
	if found {
		logrus.Infoln("sni hit", sniKey, domain)
		return domain.(string)
	} else {
		logrus.Infoln("sni miss", sniKey)
		return ""
	}
}

func makeSniKey(ip1 string, port1 string, ip2 string, port2 string) string {
	key1 := ip1 + ":" + port1
	key2 := ip2 + ":" + port2
	if key1 > key2 {
		return key1 + "->" + key2
	} else {
		return key2 + "->" + key1
	}
}

func SetDomain(ip string, domain string) {
	shortDomain := common.GetShortDomain(domain)
	rDnsCache.SetWithTTL(ip, shortDomain, 1, 1*time.Hour)
	logrus.Infoln("domain set", ip, shortDomain)
}

func QueryDomain(ip string) string {
	domain, found := rDnsCache.Get(ip)
	if found {
		logrus.Infoln("domain hit", ip, domain)
		return domain.(string)
	} else {
		logrus.Infoln("domain miss", ip)
		return ""
	}
}
