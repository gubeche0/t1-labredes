package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/gubeche0/raw-socket-t1-labredes/internal/layes"
)

var (
	MessagePort = flag.Int("message-port", 9000, "Port to send message")
	CommandPort = flag.Int("command-port", 9001, "Port to send command")

	OutputPort = flag.Int("out-port", 9090, "Port to send messages and commands")

	InterfaceName = flag.String("interface", "eth0", "Interface to listen")
)

func listingMessages() <-chan []byte {
	ch := make(chan []byte)

	go func() {

	}()

	return ch
}

func main() {
	flag.Parse()
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		fmt.Println(err)
		return
	}

	// err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	// if err != nil {
	// 	panic(err)
	// }
	// syscall.Bind(fd, &syscall.SockaddrInet4{})

	fmt.Println("Socket created")
	fmt.Println("fd:", fd)
	defer syscall.Close(fd)

	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	// net.Dial(network, address)

	for {
		buf := make([]byte, 1024)
		_, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
		}

		eth := layes.UnWrapEthernet(&buf)
		if eth.Tipo != layes.ETHERTYPE_IPV4 {
			continue
		}

		ipv4 := layes.UnWrapIpv4FromEthernet(eth)

		// @TODO: Check if is the correct destination
		// if ipv4.Destino != [4]byte{192, 168, 0, 1} {
		// 	continue
		// }

		if ipv4.Protocol != layes.IPV4_PROTOCOL_UDP {
			continue
		}

		udp, err := layes.UnWrapUdpFromIpv4(ipv4)
		if err != nil {
			log.Printf("Error on read UDP Packet: %v", err)
			continue
		}

		// if udp. {}

		if udp.DestinationPort != uint16(*MessagePort) && udp.DestinationPort != uint16(*CommandPort) {
			continue
		}

		fmt.Println(eth)
		fmt.Println(ipv4)
		fmt.Println(udp)

	}
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}
