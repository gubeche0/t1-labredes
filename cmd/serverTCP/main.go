package main

import (
	"fmt"
	"net"
	"syscall"

	"github.com/gubeche0/raw-socket-t1-labredes/internal/layes"
)

var (
	MacSource = [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	// MacSource = [6]byte{0x7e, 0xa1, 0x05, 0x5c, 0x37, 0x98}
	MacDest = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

func main() {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	// fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Socket created")
	fmt.Println("fd:", fd)

	defer syscall.Close(fd)

	if_info, err := net.InterfaceByName("eth0")
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

	ipv4 := layes.NewIpv4Layer()
	ipv4.Protocol = layes.IPV4_PROTOCOL_UDP
	ipv4.Destino = [4]byte{192, 168, 0, 1}
	ipv4.Origem = [4]byte{0, 0, 0, 0}
	// ipv4.Checksum = 0x0000

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
