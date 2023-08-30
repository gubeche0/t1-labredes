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
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	// fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Socket created")
	fmt.Println("fd:", fd)

	defer syscall.Close(fd)

	if_info, err := net.InterfaceByName(*InterfaceName)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = syscall.BindToDevice(fd, if_info.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	// err = syscall.SetLsfPromisc("eth0", true)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	udp := layes.NewUdpLayer()
	udp.SourcePort = uint16(*OutputPort)
	udp.DestinationPort = uint16(*CommandPort)
	udp.Data = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	// udp.Data = make([]byte, 1472) // Max size of UDP packet ???
	udp.Prepare()

	ipv4 := layes.NewIpv4Layer()
	ipv4.Protocol = layes.IPV4_PROTOCOL_UDP
	ipv4.Destino = [4]byte{192, 168, 0, 1}
	ipv4.Origem = [4]byte{0, 0, 0, 0}
	// ipv4.Checksum = 0x0000
	ipv4.Data = udp.ToBytes()
	ipv4.Prepare()

	eth := layes.NewEthernetLayer(MacSource, MacDest, ipv4.ToBytes())

	sockAddr := syscall.SockaddrLinklayer{
		Protocol: syscall.ETH_P_ALL,
		// Ifindex:  2,
		Ifindex: if_info.Index,
		Halen:   6,
		// Addr:     MacDest[:],
	}

	// syscall.RawSockaddrUnix{

	// }
	copy(sockAddr.Addr[:], MacDest[:])

	// syscall.Write(fd, eth.ToBytes())
	err = syscall.Sendto(fd, eth.ToBytes(), 0, &sockAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}
