package common

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

var writeApi api.WriteAPI
var hostName string

func InitInfluxdb() {
	influxdbUrl := viper.GetString("influxdb.url")
	if influxdbUrl == "" {
		return
	}

	option := influxdb2.DefaultOptions()
	option.WriteOptions().SetBatchSize(viper.GetUint("influxdb.batch_size"))
	option.WriteOptions().SetFlushInterval(viper.GetUint("influxdb.flush_interval"))

	client := influxdb2.NewClientWithOptions(viper.GetString("influxdb.url"), viper.GetString("influxdb.token"), option)

	writeApi = client.WriteAPI(viper.GetString("influxdb.org"), viper.GetString("influxdb.bucket"))

	hostName, _ = os.Hostname()

	logrus.Infoln("Influxdb init success")
}

func ReportDns(domain string, client string, server string, detail string) {
	p := influxdb2.NewPointWithMeasurement("dns").
		AddTag("host", hostName).
		AddTag("domain", domain).
		AddTag("client", client).
		AddField("detail", detail).
		SetTime(time.Now())
	writeApi.WritePoint(p)
	logrus.Debugln("dns", domain, client, "->", server, detail)
}

func ReportHttp(domain string, client string, server string, detail string) {
	p := influxdb2.NewPointWithMeasurement("http").
		AddTag("host", hostName).
		AddTag("domain", domain).
		AddTag("client", client).
		AddField("detail", detail).
		SetTime(time.Now())
	writeApi.WritePoint(p)
	logrus.Debugln("http", domain, client, "->", server, detail)
}

func ReportTls(domain string, client string, server string, detail string) {
	p := influxdb2.NewPointWithMeasurement("tls").
		AddTag("host", hostName).
		AddTag("domain", domain).
		AddTag("client", client).
		AddField("detail", detail).
		SetTime(time.Now())
	writeApi.WritePoint(p)
	logrus.Debugln("tls", domain, client, "->", server, detail)
}

func ReportTransport(domain string, client string, server string, detail string) {
	p := influxdb2.NewPointWithMeasurement("transport").
		AddTag("host", hostName).
		AddTag("domain", domain).
		AddTag("client", client).
		AddField("detail", detail).
		SetTime(time.Now())
	writeApi.WritePoint(p)
	logrus.Debugln("transport", domain, client, "->", server, detail)
}
