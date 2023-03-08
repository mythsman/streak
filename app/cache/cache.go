package cache

import (
	"github.com/dgraph-io/ristretto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"streak/app/common"
	"time"
)

var rDnsCache *ristretto.Cache
var capacity int64
var ttl int64
var metricInterval int64

func InitCache() {
	capacity = viper.GetInt64("cache.capacity")
	metricInterval = viper.GetInt64("cache.metric_interval")
	ttl = viper.GetInt64("cache.ttl")

	rDnsCache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: capacity * 10,
		MaxCost:     capacity,
		Metrics:     metricInterval > 0,
		BufferItems: 64,
	})
	if metricInterval > 0 {
		go scheduleMetrics()
	}
}

func scheduleMetrics() {
	ticker := time.NewTicker(time.Duration(metricInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		logrus.Infoln("cache metric", rDnsCache.Metrics.String())
	}
}

func SetDomain(ip string, domain string) {
	shortDomain := common.GetShortDomain(domain)
	rDnsCache.SetWithTTL(ip, shortDomain, 1, time.Duration(ttl)*time.Second)
	logrus.Debugln("cache set", ip, shortDomain)
}

func QueryDomain(ip net.IP) string {
	domain, found := rDnsCache.Get(ip.String())
	if found {
		logrus.Debugln("cache hit", ip, domain)
		return domain.(string)
	} else {
		logrus.Debugln("cache miss", ip)
		return ""
	}
}
