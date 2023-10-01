package chat

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

	requestJoin := MessageRequestJoin{
		UserName: c.User,
	}

	log.Debug().Msgf("Sending message: %v", requestJoin.Wrap())

	_, err = conn.Write(requestJoin.Wrap())
	if err != nil {
		log.Fatal().Err(err).Msg("Error to send message")
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")

		text, _ := reader.ReadString('\n')
		if len(text) == 0 {
			continue
		} else if len(text) == 1 && text[0] == '\n' {
			continue
		}

		msg := MessageText{
			Origin: c.User,
			Target: "all",
			Text:   text,
		}

		log.Debug().Msgf("Sending message: %v", msg.Wrap())
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
		buf := new(bytes.Buffer)

		var messageType uint8
		var messageSize uint32
		err := binary.Read(c.conn, binary.BigEndian, &messageType)
		if !checkErrorMessageClient(err) {
			break
		}
		err = binary.Read(c.conn, binary.BigEndian, &messageSize)
		if !checkErrorMessageClient(err) {
			break
		}

		_, err = io.CopyN(buf, c.conn, int64(messageSize-5)) // 5 bytes for message type and message size
		if !checkErrorMessageClient(err) {
			break
		}

		raw := make([]byte, 5, messageSize)
		raw[0] = messageType
		binary.BigEndian.PutUint32(raw[1:], messageSize)
		raw = append(raw, buf.Bytes()...)

		log.Debug().Msgf("Message recived: %v", raw)

		c.handlerMessage(messageType, &raw)
	}
}

func (c ClientChat) HandleControllerMessage() {

}

func (c ClientChat) handlerMessage(messageType uint8, message *[]byte) {
	switch messageType {
	case MESSAGE_TYPE_RESPONSE_JOIN:
		msg, err := UnWrapMessageResponseJoin(message)
		if err != nil {
			log.Warn().Err(err).Msg("Error to unwrap message")
			return
		}
		if !msg.Succeeded {
			log.Fatal().Msg("Error to join")
			return
		}

		log.Info().Msgf("Connected to server with user: %s", msg.UserName)

	case MESSAGE_TYPE_TEXT:
		msg, err := UnWrapMessageText(message)
		if err != nil {
			log.Warn().Err(err).Msg("Error to unwrap message")
			return
		}
		if msg.Origin == c.User {
			return
		}

		fmt.Printf("%s: %s\n", msg.Origin, msg.Text)

	default:
		log.Warn().Msgf("Message type %d not implemented", messageType)
	}

}

func (c ClientChat) SendMessage(message MessageInterface) {

}

func checkErrorMessageClient(err error) bool {
	if err == io.EOF {
		log.Fatal().Msg("Connection closed")
		return false
	} else if err != nil {
		log.Err(err).Msg("Error to read message")
		return false
	}

	return true
}
