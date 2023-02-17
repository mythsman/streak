package filter

import (
	"github.com/google/gopacket"
	"github.com/sirupsen/logrus"
	"net"
	"streak/app/cache"
	"streak/app/common"
)

func TransportFilter(packet gopacket.Packet) {
	ipSrc := net.ParseIP(packet.NetworkLayer().NetworkFlow().Src().String())
	ipDst := net.ParseIP(packet.NetworkLayer().NetworkFlow().Dst().String())

	portSrc := packet.TransportLayer().TransportFlow().Src()
	portDst := packet.TransportLayer().TransportFlow().Dst()

	if ipSrc.IsPrivate() && ipDst.IsPrivate() {
		logrus.Infoln("both private", ipSrc.String(), ipDst.String())
	} else if ipSrc.IsPrivate() && !ipDst.IsPrivate() {
		domain := cache.QueryDomain(ipDst.String())
		if domain != "" {
			common.ReportTransport(domain, ipSrc.String(), ipDst.String(), portDst.String(), len(packet.Data()))
		}
	} else if !ipSrc.IsPrivate() && ipDst.IsPrivate() {
		domain := cache.QueryDomain(ipSrc.String())
		if domain != "" {
			common.ReportTransport(domain, ipDst.String(), ipSrc.String(), portSrc.String(), len(packet.Data()))
		}
	} else {
		logrus.Infoln("may be lo ?", ipSrc.String(), ipDst.String())
	}

}
