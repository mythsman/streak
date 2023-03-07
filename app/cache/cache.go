package cache

import (
	"github.com/dgraph-io/ristretto"
	"github.com/sirupsen/logrus"
	"net"
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
	logrus.Debugln("cache set", ip, shortDomain)
}

func QueryDomain(ip net.IP) string {
	if ip.IsLoopback() || ip.IsPrivate() {
		return ""
	}
	domain, found := rDnsCache.Get(ip.String())
	if found {
		logrus.Debugln("cache hit", ip, domain)
		return domain.(string)
	} else {
		logrus.Debugln("cache miss", ip)
		return ""
	}
}
