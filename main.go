package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/dreadl0ck/tlsx"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net/http"
	"regexp"
)

func listen(networkInterface string) {
	handle, err := pcap.OpenLive(networkInterface, 1024, false, -1)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Only process transport layer (tcp , udp , not icmp)
		if packet.TransportLayer() != nil {
			ipSrc := packet.NetworkLayer().NetworkFlow().Src()
			ipDst := packet.NetworkLayer().NetworkFlow().Dst()

			portSrc := packet.TransportLayer().TransportFlow().Src()
			portDst := packet.TransportLayer().TransportFlow().Dst()

			isTls := printTls(packet)

			isHttp := printHttp(packet)

			isDns := printDns(packet)

			if isTls || isHttp || isDns {
				log.Printf("known transportLayer %s:%s -> %s:%s", ipSrc, portSrc, ipDst, portDst)
			}
		}
	}
}

func printDns(packet gopacket.Packet) bool {
	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.DNS)
		for _, answer := range dns.Answers {
			if answer.Type == layers.DNSTypeA {
				log.Printf("dns A %s %s", answer.Name, answer.IP)
			} else if answer.Type == layers.DNSTypeAAAA {
				log.Printf("dns AAAA %s %s", answer.Name, answer.IP)
			} else if answer.Type == layers.DNSTypeCNAME {
				log.Printf("dns cname %s %s", answer.Name, answer.CNAME)
			} else if answer.Type == layers.DNSTypeTXT {
				log.Printf("dns txt %s %s", answer.Name, answer.TXT)
			}
		}
		return true
	}
	return false
}

func printTls(packet gopacket.Packet) bool {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)

		if !tcp.SYN && !tcp.FIN && !tcp.RST && !(tcp.ACK && len(tcp.LayerPayload()) == 0) {
			clientHello := tlsx.GetClientHello(packet)
			if clientHello != nil {
				serverName := clientHello.SNI
				log.Printf("https %s", serverName)
				return true
			}
		}
	}
	return false
}

func printHttp(packet gopacket.Packet) bool {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		if len(tcp.Payload) != 0 {
			reader := bufio.NewReader(bytes.NewReader(tcp.Payload))
			httpReq, err := http.ReadRequest(reader)
			if err == nil {
				log.Printf("http %s %s", httpReq.Host, parsePath(httpReq.RequestURI))
				return true
			}
		}
	}
	return false
}

func parsePath(url string) string {
	pathPattern := regexp.MustCompile("(https?://[^/]*)?(/.*)")
	match := pathPattern.FindStringSubmatch(url)
	if len(match) >= 3 {
		return match[2]
	}
	return ""
}

func main() {
	cliInterface := flag.String("i", "eth0", "Network interface to listen on")

	flag.Parse()

	listen(*cliInterface)
}
