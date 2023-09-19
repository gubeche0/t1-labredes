package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"

	"github.com/gubeche0/raw-socket-t1-labredes/internal/chat"
)

var (
	MessagePort = flag.Int("message-port", 9000, "Port to send message")
	CommandPort = flag.Int("command-port", 9001, "Port to send command")
	User        = flag.String("user", "", "User to connect")

	Destination = flag.String("destination", "localhost", "Destination to send message")
)

func main() {
	flag.Parse()
	if *User == "" {
		// Generate random user
		*User = fmt.Sprintf("Anonymous_%d", rand.Intn(1000))
	}

	client := chat.ClientChat{
		User: *User,
		Conn: TCPConn{
			Address:     *Destination,
			MessagePort: *MessagePort,
			CommandPort: *CommandPort,
		},
	}

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

type TCPConn struct {
	Address     string
	MessagePort int
	CommandPort int

	readMessageChannel            *chan chat.MessageInterface
	readControllerMessageChannel  *chan chat.ControllerMessage
	writeMessageChannel           *chan chat.MessageInterface
	writeControllerMessageChannel *chan chat.ControllerMessage
}

func (c TCPConn) Conection() error {

	c.handleControllerMessage()

	readMessageChannel := make(chan chat.MessageInterface)
	c.readMessageChannel = &readMessageChannel
	return nil
}

func (c TCPConn) GetReadMessageChannel() <-chan chat.MessageInterface {
	return *c.readMessageChannel
}

func (c TCPConn) GetReadControllerMessageChannel() <-chan chat.ControllerMessage {
	return *c.readControllerMessageChannel
}

func (c TCPConn) GetWriteMessageChannel() chan<- chat.MessageInterface {
	return *c.writeMessageChannel
}

func (c TCPConn) GetWriteControllerMessageChannel() chan<- chat.ControllerMessage {
	return *c.writeControllerMessageChannel
}

func (c TCPConn) Close() error {

	close(*c.readMessageChannel)
	close(*c.readControllerMessageChannel)
	close(*c.writeMessageChannel)
	close(*c.writeControllerMessageChannel)

	return nil
}

func (c *TCPConn) handleControllerMessage() {
	readControllerMessageChannel := make(chan chat.ControllerMessage)
	c.readControllerMessageChannel = &readControllerMessageChannel

	writeControllerMessageChannel := make(chan chat.ControllerMessage)
	c.writeControllerMessageChannel = &writeControllerMessageChannel

	conexao, erro1 := net.Dial("tcp", "127.0.0.1:8081")
	if erro1 != nil {
		fmt.Println(erro1)
		os.Exit(3)
	}

	go func() {
		message := <-*c.writeControllerMessageChannel
		conexao.Write(message.Wrap())

		// conexao.Read(make([]byte, 1))
		// conexao.Read(make([]byte, 1))
		mensagem, err3 := bufio.NewReader(conexao).ReadBytes('\n')
		// bufio.NewScanner(conexao).Scan()
		// mensagem := bufio.NewScanner(conexao).Bytes()
		if err3 != nil {
			fmt.Println(err3)
			os.Exit(3)
		}
		fmt.Println("Mensagem do servidor: ", mensagem)
		msg, err := chat.UnWrapControllerMessageRaw(&mensagem)
		if err != nil {
			log.Print("Error to unwrap message")
		}

		*c.readControllerMessageChannel <- msg
	}()
}

// Conection() error
// GetReadMessageChannel() <-chan MessageInterface
// GetReadControllerMessageChannel() <-chan ControllerMessage
// GetWriteMessageChannel() chan<- MessageInterface
// GetWriteControllerMessageChannel() chan<- ControllerMessage
// Close() error
