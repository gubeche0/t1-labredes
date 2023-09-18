package main

import (
	"flag"
	"fmt"
	"net"
	"syscall"

	"github.com/gubeche0/raw-socket-t1-labredes/internal/layes"
)

var (
	MessagePort = flag.Int("message-port", 9000, "Port to recive message")
	CommandPort = flag.Int("command-port", 9001, "Port to recive command")

	OutputPort = flag.Int("out-port", 9090, "Port to send messages and commands")

	InterfaceName = flag.String("interface", "eth0", "Interface to use")

	MacSource = [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	MacDest   = [6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

func main() {
	flag.Parse()

	// fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_PACKET, syscall.IPPROTO_RAW)
	// fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	// fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	// fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_PACKET, syscall.IPPROTO_TCP)

	if err != nil {
		fmt.Println("Error creating socket:", err)
		return
	}

	fmt.Println("Socket created")
	fmt.Println("fd:", fd)

	defer syscall.Close(fd)

	if_info, err := net.InterfaceByName(*InterfaceName)
	if err != nil {
		fmt.Println("Error getting interface:", err)
		return
	}

	err = syscall.BindToDevice(fd, if_info.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	// sockAddr := syscall.SockaddrLinklayer{
	// 	// Protocol: htons(syscall.ETH_P_ALL),
	// 	// Protocol: syscall.ETH_P_ALL,
	// 	// Protocol: syscall.IPPROTO_TCP,
	// 	Protocol: syscall.IPPROTO_RAW,
	// 	// Protocol: syscall.,
	// 	// Ifindex:  2,
	// 	Ifindex: if_info.Index,
	// 	// Halen:   6,

	// 	// Hatype:  syscall.ARPHRD_ETHER,
	// 	// Pkttype: syscall.PACKET_HOST,
	// 	// Pkttype: 0,

	// 	// Addr:     MacDest[:],
	// }

	// sockAddr := syscall.SockaddrInet4{
	// 	Port: *MessagePort,
	// 	Addr: [4]byte{127, 0, 0, 1},
	// }

	sockAddr := syscall.SockaddrLinklayer{
		// Protocol: syscall.ETH_P_ALL,
		// Protocol: htons(syscall.ETH_P_ALL),
		// Protocol: syscall.IPPROTO_TCP,
		Protocol: htons(syscall.IPPROTO_TCP),
		// Protocol: 0x0800,
		Ifindex: if_info.Index,
		// Halen:   6,
		// Hatype:  syscall.ARPHRD_ETHER,
		// Pkttype: syscall.PACKET_HOST,
		// Pkttype: 0x8,
	}
	copy(sockAddr.Addr[:], MacDest[:])
	// copy(sockAddr.Addr[7:], []byte{0x23, 0x28})

	// copy(sockAddr.Addr[:], MacDest[:])
	// copy(sockAddr.Addr[:], []byte{6})

	data := []byte{0x01, 0x02, 0x03, 0x04}

	tcp := layes.NewTcpLayer()
	tcp.SourcePort = uint16(*OutputPort)
	tcp.DestinationPort = uint16(*MessagePort)
	tcp.Data = data
	tcp.Prepare()
	tcp.Checksum = 0x6530

	ipv4 := layes.NewIpv4Layer()
	ipv4.Protocol = layes.IPV4_PROTOCOL_TCP
	ipv4.Destino = [4]byte{255, 255, 255, 255}
	ipv4.Origem = [4]byte{0, 0, 0, 0}
	ipv4.Data = tcp.ToBytes()
	ipv4.Prepare()
	ipv4.Checksum = 0x7acd

	eth := layes.NewEthernetLayer(MacSource, MacDest, ipv4.ToBytes())

	// err = syscall.Sendto(fd, ipv4.ToBytes(), 0, &sockAddr)
	err = syscall.Sendto(fd, eth.ToBytes(), 0, &sockAddr)
	// err = syscall.Sendto(fd, make([]byte, 30), 0, &sockAddr)
	if err != nil {

		fmt.Println("Error sending data:", err)
		return
	}
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}
