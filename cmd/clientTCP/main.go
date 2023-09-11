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

	// err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 0)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	// err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.IP_HDRINCL, 1)

	tcp := layes.NewTcpLayer()
	// udp.Data = make([]byte, 1472) // Max size of UDP packet ???
	tcp.SourcePort = uint16(*OutputPort)
	tcp.DestinationPort = uint16(*CommandPort)
	tcp.Data = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	// tcp.Data = make([]byte, 1472)
	tcp.Prepare()

	ipv4 := layes.NewIpv4Layer()
	ipv4.Protocol = layes.IPV4_PROTOCOL_TCP
	ipv4.Destino = [4]byte{192, 168, 0, 1}
	ipv4.Origem = [4]byte{0, 0, 0, 0}
	// ipv4.Checksum = 0x0000
	ipv4.Data = tcp.ToBytes()

	eth := layes.NewEthernetLayer(MacSource, MacDest, ipv4.ToBytes())

	sockAddr := syscall.SockaddrLinklayer{
		// Protocol: htons(syscall.ETH_P_ALL),
		Protocol: syscall.ETH_P_ALL,
		// Protocol: syscall.IPPROTO_TCP,
		// Protocol: syscall.,
		// Ifindex:  2,
		Ifindex: if_info.Index,
		// Halen:   6,

		// Hatype:  syscall.ARPHRD_ETHER,
		// Pkttype: syscall.PACKET_HOST,
		Pkttype: 0,

		// Addr:     MacDest[:],
	}

	// net.ListenIP()
	// syscall.RawSockaddrUnix{

	// }
	copy(sockAddr.Addr[:], MacDest[:])
	copy(sockAddr.Addr[:], []byte{6})

	// err = syscall.Bind(fd, &sockAddr)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

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
