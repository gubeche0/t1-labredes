package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/gubeche0/raw-socket-t1-labredes/internal/layes"
)

func main() {
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
		// fmt.Printf("% X\n", buf[:numRead])
		// fmt.Printf("MAC Destino: %x:%x:%x:%x:%x:%x \n", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
		// fmt.Printf("MAC Origen: %x:%x:%x:%x:%x:%x \n\n", buf[6], buf[7], buf[8], buf[9], buf[10], buf[11])
		// fmt.Printf("Tipo: %x \n\n", buf[12:14])

		eth := layes.UnWrapEthernet(&buf)
		ipv4 := layes.UnWrapIpv4FromEthernet(eth)
		udp, err := layes.UnWrapUdpFromIpv4(ipv4)
		if err == nil {
			fmt.Println(eth)
			fmt.Println(ipv4)
			fmt.Println(udp)
		}
	}
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}
