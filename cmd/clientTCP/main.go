package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/gubeche0/raw-socket-t1-labredes/internal/layes"
)

var (
	MessagePort = flag.Int("message-port", 9000, "Port to send message")
	CommandPort = flag.Int("command-port", 9001, "Port to send command")

	OutputPort = flag.Int("out-port", 9090, "Port to send messages and commands")
)

func main() {
	flag.Parse()
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))

	// net.DialUnix()
	if err != nil {
		fmt.Println(err)
		return
	}

	// err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	// if err != nil {
	// 	panic(err)
	// }

	//

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
		// if eth.Origem != [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00} {
		// 	continue
		// }
		ipv4 := layes.UnWrapIpv4FromEthernet(eth)
		if ipv4.Protocol == layes.IPV4_PROTOCOL_TCP {
			// tcp, err := layes.UnWrapTcpFromIpv4(ipv4)
			// if err == nil {
			// 	fmt.Println(eth)
			// 	fmt.Println(ipv4)
			// 	fmt.Println(tcp)
			// }
		} else if ipv4.Protocol == layes.IPV4_PROTOCOL_UDP {
			udp, err := layes.UnWrapUdpFromIpv4(ipv4)
			if err == nil {
				fmt.Println(eth)
				fmt.Println(ipv4)
				fmt.Println(udp)
			}
		} else {
			fmt.Println(eth)
			fmt.Println(ipv4)
		}

	}
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}
