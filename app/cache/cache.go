package cache

import (
	"fmt"
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

func SetSni(domain string, ipSrc string, portSrc string, ipDst string, portDst string) {
	shortDomain := common.GetShortDomain(domain)
	sniCache.SetWithTTL(makeSniKey(ipSrc, portSrc, ipDst, portDst), shortDomain, 1, 1*time.Hour)
	logrus.Infoln("sni set", ipDst, shortDomain)
}

func QuerySni(ipSrc string, portSrc string, ipDst string, portDst string) string {
	domain, found := sniCache.Get(makeSniKey(ipSrc, portSrc, ipDst, portDst))
	if found {
		logrus.Infoln("sni hit", ipDst, domain)
		return domain.(string)
	} else {
		logrus.Infoln("sni miss", ipDst)
		return ""
	}
}

func makeSniKey(ipSrc string, portSrc string, ipDst string, portDst string) string {
	return fmt.Sprintf("%s:%s->%s:%s", ipSrc, portSrc, ipDst, portDst)
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
