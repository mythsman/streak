package filter

import (
	"bufio"
	"bytes"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net/http"
	"regexp"
	"streak/app/common"
)

func HttpFilter(packet gopacket.Packet) {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		if len(tcp.Payload) != 0 {
			reader := bufio.NewReader(bytes.NewReader(tcp.Payload))
			httpReq, err := http.ReadRequest(reader)
			if err == nil {
				ipSrc := packet.NetworkLayer().NetworkFlow().Src()
				ipDst := packet.NetworkLayer().NetworkFlow().Dst()
				common.ReportHttp(httpReq.Host, parsePath(httpReq.RequestURI), ipSrc.String(), ipDst.String())
			}
		}
	}
}

func parsePath(url string) string {
	pathPattern := regexp.MustCompile("(https?://[^/]*)?(/.*)")
	match := pathPattern.FindStringSubmatch(url)
	if len(match) >= 3 {
		return match[2]
	}
	return ""
}
