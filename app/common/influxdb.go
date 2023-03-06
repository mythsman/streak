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

func ReportDns(domain string, client string, detail string) {
	p := influxdb2.NewPointWithMeasurement("dns").
		AddTag("domain", domain).
		AddTag("client", client).
		AddField("detail", detail).
		SetTime(time.Now())
	writeApi.WritePoint(p)
	logrus.Infoln("dns", domain, client, detail)
}

func ReportHttp(domain string, client string, detail string) {
	p := influxdb2.NewPointWithMeasurement("http").
		AddTag("domain", domain).
		AddTag("client", client).
		AddField("detail", detail).
		SetTime(time.Now())
	writeApi.WritePoint(p)
	logrus.Infoln("http", domain, client, detail)
}

func ReportTls(domain string, client string, detail string) {
	p := influxdb2.NewPointWithMeasurement("tls").
		AddTag("domain", domain).
		AddTag("client", client).
		AddField("detail", detail).
		SetTime(time.Now())
	writeApi.WritePoint(p)
	logrus.Infoln("tls", domain, client, detail)
}

func ReportTransport(domain string, client string, detail string) {
	p := influxdb2.NewPointWithMeasurement("transport").
		AddTag("domain", domain).
		AddTag("client", client).
		AddField("detail", detail).
		SetTime(time.Now())
	writeApi.WritePoint(p)
	logrus.Infoln("transport", domain, client, detail)
}
