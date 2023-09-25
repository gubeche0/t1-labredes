package chat

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type ClientChat struct {
	conn        net.Conn
	User        string
	Address     string
	MessagePort int
	CommandPort int
}

func (c *ClientChat) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Address, c.MessagePort))
	if err != nil {
		log.Fatal().Msgf("Error to connect: %s", err.Error())
	}
	defer conn.Close()
	c.conn = conn

	go c.listenMessage()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")

		text, _ := reader.ReadString('\n')

		msg := MessageText{
			Origin: c.User,
			Target: "all",
			Text:   text,
		}

		_, err := conn.Write(msg.Wrap())
		if err != nil {
			log.Fatal().Err(err).Msg("Error to send message")
		}

		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")

			return errors.New("TCP client exiting")
		}
	}
}

func (c ClientChat) listenMessage() {
	for {
		message, err := bufio.NewReader(c.conn).ReadBytes('\n')
		if err == io.EOF {
			log.Fatal().Msg("Connection closed")
		} else if err != nil {
			fmt.Println(err)
			return
		}
		if len(message) == 0 {
			continue
		}
		fmt.Print("Server: " + string(message))
	}
}

func (c ClientChat) HandleControllerMessage() {

}

func (c ClientChat) SendMessage(message MessageInterface) {

}
