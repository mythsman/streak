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
		reportInnerOuter(ipSrc.String(), portSrc.String(), ipDst.String(), portDst.String(), len(packet.Data()))
	} else if !ipSrc.IsPrivate() && ipDst.IsPrivate() {
		reportInnerOuter(ipDst.String(), portDst.String(), ipSrc.String(), portSrc.String(), len(packet.Data()))
	} else {
		logrus.Infoln("may be lo ?", ipSrc.String(), ipDst.String())
	}
}

func reportInnerOuter(innerIp string, innerPort string, outerIp string, outerPort string, packetSize int) {
	domain := cache.QuerySni(innerIp, innerPort, outerIp, outerPort)
	if domain != "" {
		common.ReportTransport(domain, innerIp, outerIp, outerPort, packetSize)
		return
	}

	domain = cache.QueryDomain(outerIp)
	if domain != "" {
		common.ReportTransport(domain, innerIp, outerIp, outerPort, packetSize)
		return
	}
}
