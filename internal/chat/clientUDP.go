package chat

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type ClientUDP struct {
	conn        *net.UDPConn
	CommandConn *net.UDPConn
	User        string
	Address     string
	MessagePort int
	CommandPort int
}

func (c *ClientUDP) Connect() error {
	commandConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(c.Address),
		Port: c.CommandPort,
	})

	if err != nil {
		log.Fatal().Msgf("Error to connect: %s", err.Error())
	}
	defer commandConn.Close()
	c.CommandConn = commandConn

	c.requestUserName()

	go c.listenCommandMessage()

	// conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", c.Address, c.MessagePort))
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(c.Address),
		Port: c.MessagePort,
	})
	if err != nil {
		log.Fatal().Msgf("Error to connect: %s", err.Error())
	}
	defer conn.Close()
	c.conn = conn

	go c.listenMessage()

	requestJoin := MessageJoin{
		UserName: c.User,
	}

	log.Debug().Msgf("Sending message: %v", requestJoin.Wrap())

	reader := bytes.NewReader(requestJoin.Wrap())
	_, err = io.Copy(conn, reader)
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

func (c ClientUDP) requestUserName() {
	requestJoin := MessageJoinRequest{
		UserName: c.User,
	}
	log.Debug().Msgf("Request join with user: %s", c.User)

	_, err := c.CommandConn.Write(requestJoin.Wrap())
	if err != nil {
		log.Fatal().Err(err).Msg("Error to send message")
	}

	// log.Debug().Msgf("Waiting response to join with user: %s", c.User)
	raw, _, err := reciveMessageUDP(c.CommandConn)
	if err != nil {
		log.Fatal().Err(err).Msg("Error to recive message")
	}
	// log.Debug().Msgf("Response recived: %v", raw)

	msg, err := UnWrapMessageJoinRequestResponse(raw)
	if err != nil {
		log.Fatal().Err(err).Msg("Error to unwrap message")
	}

	if !msg.Succeeded {
		log.Fatal().Msgf("Error to request join with user: %s", c.User)
	}
}

func (c ClientUDP) listenMessage() {
	for {
		raw, _, err := reciveMessageUDP(c.conn)
		checkErrorMessageClient(err)

		log.Debug().Msgf("Message recived: %v", raw)

		c.handlerMessage((*raw)[0], raw)
	}
}

func (c ClientUDP) listenCommandMessage() {
	for {
		raw, _, err := reciveMessageUDP(c.CommandConn)
		if !checkErrorMessage(err) {
			break
		}

		log.Debug().Msgf("Command message recived: %v", raw)

		c.HandleCommandMessage((*raw)[0], raw)
	}
}

func (c ClientUDP) HandleCommandMessage(messageType uint8, message *[]byte) {
	switch messageType {
	default:
		log.Warn().Msgf("Command message type %d not implemented", messageType)
	}
}

func (c ClientUDP) handlerMessage(messageType uint8, message *[]byte) {
	switch messageType {
	case MESSAGE_TYPE_JOIN_REQUEST_RESPONSE:
		msg, err := UnWrapMessageJoinRequestResponse(message)
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

	case MESSAGE_TYPE_LIST_USERS_RESPONSE:
		msg, err := UnWrapMessageListUserResponse(message)
		if err != nil {
			log.Warn().Err(err).Msg("Error to unwrap message")
			return
		}

		log.Debug().Msgf("Users: %v", msg.Users)
		fmt.Println("Users:")
		for _, user := range msg.Users {
			fmt.Println("  ", user)
		}

	case MESSAGE_TYPE_FILE:
		msg, err := UnWrapMessageFile(message)
		if err != nil {
			log.Warn().Err(err).Msg("Error to unwrap message")
			return
		}

		if msg.Origin == c.User {
			return
		}

		dir := fmt.Sprintf("./dl-%s", c.User)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				log.Warn().Err(err).Msg("Error to create directory")
				return
			}
		}

		file, err := os.Create(fmt.Sprintf("%s/%s", dir, msg.Filename))
		if err != nil {
			log.Warn().Err(err).Msg("Error to create file")
			return
		}

		_, err = file.Write(msg.Filedata)
		if err != nil {
			log.Warn().Err(err).Msg("Error to write file")
			return
		}

		log.Info().Msgf("File %s received from %s", msg.Filename, msg.Origin)

	default:
		log.Warn().Msgf("Message type %d not implemented", messageType)
	}

}

func (c ClientUDP) SendMessage(message MessageInterface) {
	log.Debug().Msgf("Sending message: %v", message.Wrap())

	reader := bytes.NewReader(message.Wrap())
	_, err := io.Copy(c.conn, reader)
	if err != nil {
		log.Fatal().Err(err).Msg("Error to send message")
	}
}

func (c ClientUDP) handlerInput(text string) {
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

func (c *ClientUDP) handleCommand(text string) {
	command := strings.TrimSpace(text[1:])
	commandArgs := strings.Split(command, " ")

	command = strings.ToLower(commandArgs[0])

	switch command {
	case "sendprivate":
		c.handleSendPrivate(commandArgs)
	case "sendfile":
		c.handleCommandSendFile(commandArgs)
	case "sendfileto":
		c.handleCommandSendFile(commandArgs)
	case "exit":
		log.Info().Msg("TCP client exiting...")
		os.Exit(0)
	// case "close":
	// 	c.CommandConn.Close()
	case "listusers":
		c.SendMessage(MessageListUser{
			origin: c.User,
		})
	case "help":
		c.handleCommandHelp()
	case "clear":
		fmt.Print("\033[H\033[2J")
	default:
		log.Warn().Msgf("Command %s not implemented", commandArgs[0])
	}
}

// func (c ClientChat) handleCommandList() {}

func (c ClientUDP) handleSendPrivate(commandArgs []string) {
	if len(commandArgs) < 3 {
		log.Warn().Msg("Invalid arguments")
		return
	}

	destination := commandArgs[1]
	text := strings.Join(commandArgs[2:], " ")

	c.SendMessage(MessageText{
		Origin: c.User,
		Target: destination,
		Text:   text,
	})
}

func (c ClientUDP) handleCommandSendFile(commandArgs []string) {
	if commandArgs[0] == "sendfile" && len(commandArgs) < 2 {
		log.Warn().Msg("Invalid arguments")
		return
	} else if commandArgs[0] == "sendfileto" && len(commandArgs) < 3 {
		log.Warn().Msg("Invalid arguments")
		return
	}

	target := MESSAGE_TARGET_ALL
	upload := commandArgs[1]
	if commandArgs[0] == "sendfileto" {
		target = commandArgs[1]
		upload = commandArgs[2]

	}

	file, err := os.Open(upload)
	if err != nil {
		log.Warn().Err(err).Msg("Error to open file")
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Warn().Err(err).Msg("Error to get file info")
		return
	}

	fileData := make([]byte, fileInfo.Size())
	_, err = file.Read(fileData)
	if err != nil {
		log.Warn().Err(err).Msg("Error to read file")
		return
	}

	filename := fileInfo.Name()
	filename = strings.ReplaceAll(filename, "\n", "")
	filename = strings.ReplaceAll(filename, "\r", "")
	filename = strings.ReplaceAll(filename, "\t", "")

	msg := MessageFile{
		Origin:   c.User,
		Target:   target,
		Filename: filename,
		Filesize: uint64(fileInfo.Size()),
		Filedata: fileData,
	}

	// log.Debug().Msgf("Sending file: %v", msg)

	c.SendMessage(msg)
}

func (c ClientUDP) handleCommandHelp() {
	fmt.Println("Help commands:")
	fmt.Println("  /listUsers")
	fmt.Println("  /sendPrivate <user> <message>")
	fmt.Println("  /sendFile <file>")
	fmt.Println("  /sendFileTo <user> <file>")
	fmt.Println("  /exit")
	fmt.Println("  /help")
	fmt.Println("  /clear")
}
