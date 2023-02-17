package filter

import (
	"github.com/dreadl0ck/tlsx"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"streak/app/cache"
	"streak/app/common"
	"strings"
)

func TlsFilter(packet gopacket.Packet) {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)

		if !tcp.SYN && !tcp.FIN && !tcp.RST && !(tcp.ACK && len(tcp.LayerPayload()) == 0) {
			clientHello := tlsx.GetClientHello(packet)
			if clientHello != nil {
				serverName := clientHello.SNI
				if serverName != "" {
					ipSrc := packet.NetworkLayer().NetworkFlow().Src()
					ipDst := packet.NetworkLayer().NetworkFlow().Dst()

					portSrc := packet.TransportLayer().TransportFlow().Src()
					portDst := packet.TransportLayer().TransportFlow().Dst()

					domain := strings.ToLower(serverName)

					common.ReportTls(domain, common.GetShortDomain(domain), ipSrc.String(), ipDst.String())

					cache.SetSni(domain, ipSrc.String(), portSrc.String(), ipDst.String(), portDst.String())
				}
			}
		}
	}
}
