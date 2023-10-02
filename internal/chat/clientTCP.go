package chat

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

		c.handlerInput(text)
	}
}

func (c ClientChat) listenMessage() {
	for {
		buf := new(bytes.Buffer)

		var messageType uint8
		var messageSize uint64
		err := binary.Read(c.conn, binary.BigEndian, &messageType)
		checkErrorMessageClient(err)

		err = binary.Read(c.conn, binary.BigEndian, &messageSize)
		checkErrorMessageClient(err)

		_, err = io.CopyN(buf, c.conn, int64(messageSize-9)) // 9 bytes for message type(1) and message size(8)
		checkErrorMessageClient(err)

		raw := make([]byte, 9, messageSize)
		raw[0] = messageType
		binary.BigEndian.PutUint64(raw[1:], messageSize)
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
	log.Debug().Msgf("Sending message: %v", message.Wrap())
	_, err := c.conn.Write(message.Wrap())
	if err != nil {
		log.Fatal().Err(err).Msg("Error to send message")
	}
}

func (c ClientChat) handlerInput(text string) {
	if len(text) == 0 {
		return
	} else if len(text) == 1 && text[0] == '\n' {
		return
	}

	if strings.HasPrefix(text, "/") {
		c.handleCommand(text)
		return
	}

	msg := MessageText{
		Origin: c.User,
		Target: MESSAGE_TARGET_ALL,
		Text:   text,
	}

	c.SendMessage(msg)
}

func (c ClientChat) handleCommand(text string) {
	command := strings.TrimSpace(text[1:])
	command = strings.ToLower(command)
	commandArgs := strings.Split(command, " ")

	switch commandArgs[0] {
	// case "listusers":
	// c.handleCommandList()
	// case "sendprivate":
	// c.handleSendPrivate(commandArgs)
	// case "sendfile":
	// c.handleCommandSendFile(commandArgs)
	case "exit":
		log.Info().Msg("TCP client exiting...")
		os.Exit(0)
	case "help":
		c.handleCommandHelp()
	default:
		log.Warn().Msgf("Command %s not implemented", commandArgs[0])
	}
}

// func (c ClientChat) handleCommandList() {}
// func (c ClientChat) handleSendPrivate(commandArgs []string) {}
// func (c ClientChat) handleCommandSendFile(commandArgs []string) {}

func (c ClientChat) handleCommandHelp() {
	fmt.Println("Help commands:")
	fmt.Println("  /listUsers")
	fmt.Println("  /sendPrivate <user> <message>")
	fmt.Println("  /sendFile <user> <file>")
	fmt.Println("  /exit")
	fmt.Println("  /help")
}
