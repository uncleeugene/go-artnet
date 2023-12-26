package main

import (
	"fmt"
	"net"
	"time"

	artnet "github.com/uncleeugene/go-artnet"
	"github.com/uncleeugene/go-artnet/packet"
)

func main() {

	dst := fmt.Sprintf("%s:%d", "255.255.255.255", packet.ArtNetPort)
	broadcastAddr, _ := net.ResolveUDPAddr("udp", dst)
	src := fmt.Sprintf("%s:%d", "2.12.12.12", packet.ArtNetPort)
	localAddr, _ := net.ResolveUDPAddr("udp", src)

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Printf("error opening udp: %s\n", err)
		return
	}

	p := &packet.ArtPollPacket{}
	b, err := p.MarshalBinary()
	if err != nil {
		fmt.Printf("error marshalling packet: %s\n", err)
		return
	}

	n, err := conn.WriteTo(b, broadcastAddr)
	if err != nil {
		fmt.Printf("error writing packet: %s\n", err)
		return
	}
	fmt.Printf("packet sent, wrote %d bytes\n", n)

	// wait 5 seconds for a reply
	timer := time.NewTimer(5 * time.Second)

	recvCh := make(chan []byte)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, addr, err := conn.ReadFromUDP(buf) // first packet you read will be your own
			if err != nil {
				fmt.Printf("error reading packet: %s\n", err)
				continue

			}
			fmt.Printf("packet received from %v, read %d bytes\n", addr.IP, n)
			if addr.IP.Equal(localAddr.IP) {
				// skip messages from myself
				continue
			}
			recvCh <- buf[:n]
		}
	}()

	for {
		select {
		case b := <-recvCh:
			p, err := packet.Unmarshal(b)
			if err != nil {
				fmt.Printf("error unmarshalling packet: %s\n", err)
				continue
			}
			cf := artnet.ConfigFromArtPollReply(*p.(*packet.ArtPollReplyPacket))
			fmt.Printf("got reply: %#v\n", cf)

		case <-timer.C:
			fmt.Printf("timeout reached\n")
			return
		}
	}
}
