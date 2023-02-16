package filter

import (
	"github.com/dreadl0ck/tlsx"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
)

func TlsFilter(packet gopacket.Packet) {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)

		if !tcp.SYN && !tcp.FIN && !tcp.RST && !(tcp.ACK && len(tcp.LayerPayload()) == 0) {
			clientHello := tlsx.GetClientHello(packet)
			if clientHello != nil {
				serverName := clientHello.SNI
				if serverName != "" {
					logrus.Printf("https %s", serverName)
					return
				}
			}
		}
	}
}
