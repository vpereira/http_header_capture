package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func pumpOutChannel(messages chan string) {
	for cookie := range messages {
		fmt.Println(cookie)
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
					messages <- cookie
				}
			}
		}
	}
}
