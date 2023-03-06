package filter

import (
	"github.com/dreadl0ck/tlsx"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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

					domain := strings.ToLower(serverName)

					common.ReportTls(common.GetShortDomain(domain), ipSrc.String(), domain)
				}
			}
		}
	}
}
