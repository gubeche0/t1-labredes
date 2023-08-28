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

	// syscall.Sendto()

	fmt.Println("Socket created")
	fmt.Println("fd:", fd)

	defer syscall.Close(fd)

	if_info, err := net.InterfaceByName("eth0")
	if err != nil {
		fmt.Println(err)
		return
	}

	// syscall.BindToDevice(fd, if_info.Name)
	err = syscall.BindToDevice(fd, "eth0")
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

	eth := layes.EthernetLayer{
		Origem:  MacSource,
		Destino: MacDest,
		Tipo:    [2]byte{0x08, 0x00},
		Data:    []byte{},
	}

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
