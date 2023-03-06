package cache

import (
	"github.com/dgraph-io/ristretto"
	"github.com/sirupsen/logrus"
	"streak/app/common"
	"time"
)

var rDnsCache *ristretto.Cache

func init() {
	rDnsCache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1 << 24, // 16M
		MaxCost:     1 << 27, // 128MB
		BufferItems: 64,
	})
}

func SetDomain(ip string, domain string) {
	shortDomain := common.GetShortDomain(domain)
	rDnsCache.SetWithTTL(ip, shortDomain, 1, 1*time.Hour)
	logrus.Debugln("domain set", ip, shortDomain)
}

func QueryDomain(ip string) string {
	domain, found := rDnsCache.Get(ip)
	if found {
		logrus.Debugln("domain hit", ip, domain)
		return domain.(string)
	} else {
		logrus.Debugln("domain miss", ip)
		return ""
	}
}
