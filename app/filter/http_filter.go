package filter

import (
	"bufio"
	"bytes"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"net/http"
	"regexp"
	"streak/app/common"
	"strings"
)

func HttpFilter(packet gopacket.Packet) {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		if len(tcp.Payload) != 0 {
			reader := bufio.NewReader(bytes.NewReader(tcp.Payload))
			httpReq, err := http.ReadRequest(reader)
			if err == nil {
				ipSrc := packet.NetworkLayer().NetworkFlow().Src().String()
				ipDst := packet.NetworkLayer().NetworkFlow().Dst().String()

				host := strings.ToLower(httpReq.Host)

				if strings.Contains(host, ":") {
					host = strings.Split(host, ":")[0]
				}
				path := "http://" + host + parsePath(httpReq.RequestURI)
				if net.ParseIP(host) == nil {
					host = common.GetShortDomain(host)
				}
				common.ReportHttp(host, ipSrc, ipDst, path)
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
