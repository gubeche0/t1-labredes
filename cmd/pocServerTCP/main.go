package main

import (
	"bufio"
	"flag"
	"log"
	"net"
)

var (
	MessagePort = flag.Int("message-port", 9000, "Port to recive message")
	CommandPort = flag.Int("command-port", 9001, "Port to recive command")
	User        = flag.String("user", "", "User to connect")

	Listen = flag.String("listen", "0:0:0:0", "Listen to connect")
)

func main() {
	flag.Parse()

	serve, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Error to create server: ", err)
	}

	// for {
	conn, err := serve.Accept()
	if err != nil {
		log.Fatal("Error to accept connection: ", err)
	}

	// go handleConnection(conn)

	defer conn.Close()
	log.Println("Connection accepted")

	for {
		message, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			log.Println("Error to read message: ", err)
			break
		}

		log.Println("Message recived: ", string(message))

		// conn.Write(message)
	}
	// }
}
