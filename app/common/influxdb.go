package common

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var writeApi api.WriteAPI

func InitInfluxdb() {
	influxdbUrl := viper.GetString("influxdb.url")
	if influxdbUrl == "" {
		return
	}

	client := influxdb2.NewClient(viper.GetString("influxdb.url"), viper.GetString("influxdb.token"))

	writeApi = client.WriteAPI(viper.GetString("influxdb.org"), viper.GetString("influxdb.bucket"))

	logrus.Infoln("Influxdb init success")
}

func ReportDns(domain string, shortDomain string, queryType string, client string, server string) {
	p := influxdb2.NewPointWithMeasurement("dns").
		AddTag("domain", shortDomain).
		AddTag("server", server).
		AddTag("client", client).
		AddTag("type", queryType).
		AddField("domain", domain).
		SetTime(time.Now())
	writeApi.WritePoint(p)
}

func ReportHttp(host string, path string, client string, server string) {
	p := influxdb2.NewPointWithMeasurement("http").
		AddTag("host", host).
		AddTag("server", server).
		AddTag("client", client).
		AddField("path", path).
		SetTime(time.Now())
	writeApi.WritePoint(p)
}

func ReportTls(domain string, client string, server string) {
	p := influxdb2.NewPointWithMeasurement("tls").
		AddTag("domain", domain).
		AddTag("server", server).
		AddTag("client", client).
		SetTime(time.Now())
	writeApi.WritePoint(p)
}

func ReportUdp(domain string, client string, server string, port string, tx uint64, rx uint64) {
	p := influxdb2.NewPointWithMeasurement("udp").
		AddTag("domain", domain).
		AddTag("server", server).
		AddTag("port", port).
		AddTag("client", client).
		AddField("rx", rx).
		AddField("tx", tx).
		SetTime(time.Now())
	writeApi.WritePoint(p)
}

func ReportTcp(domain string, client string, server string, port string, tx uint64, rx uint64) {
	p := influxdb2.NewPointWithMeasurement("tcp").
		AddTag("domain", domain).
		AddTag("server", server).
		AddTag("port", port).
		AddTag("client", client).
		AddField("rx", rx).
		AddField("tx", tx).
		SetTime(time.Now())
	writeApi.WritePoint(p)
}
