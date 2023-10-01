package chat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/rs/zerolog/log"
)

type ServerTCP struct {
	// Protocol    Protocol
	Users       map[string]connection
	Address     string
	MessagePort int
	CommandPort int
}

type connection struct {
	UserName    string
	MessageConn net.Conn
	CommandConn net.Conn
}

func (s *ServerTCP) StartListenAndServer() error {
	serve, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Address, s.MessagePort))
	if err != nil {
		log.Fatal().Msgf("Error to create server: %s", err.Error())

		return err
	}

	for {
		conn, err := serve.Accept()
		if err != nil {
			log.Fatal().Msgf("Error to accept connection: %s", err.Error())
		}

		go s.handleConnection(conn)
	}
}

func (s *ServerTCP) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Info().Msg("Connection accepted")
	// s.Users[""] = connection{
	// 	MessageConn: conn,
	// }

	for {
		buf := new(bytes.Buffer)

		var messageType uint8
		var messageSize uint32
		err := binary.Read(conn, binary.BigEndian, &messageType)
		if !checkErrorMessage(err) {
			break
		}
		err = binary.Read(conn, binary.BigEndian, &messageSize)
		if !checkErrorMessage(err) {
			break
		}

		log.Debug().Msgf("Message type recived: %v", messageType)
		log.Debug().Msgf("Message size recived: %v", messageSize)

		// message, err := bufio.NewReader(conn).ReadBytes('\n')
		_, err = io.CopyN(buf, conn, int64(messageSize))
		if !checkErrorMessage(err) {
			break
		}

		raw := make([]byte, 5+messageSize)
		raw[0] = messageType
		binary.BigEndian.PutUint32(raw[1:], messageSize)
		raw = append(raw, buf.Bytes()...)

		log.Debug().Msgf("Message recived: %v", raw)

		s.handlerMessage(messageType, &raw)

	}
}

func (s ServerTCP) handlerMessage(messageType uint8, message *[]byte) {
	switch messageType {
	case MESSAGE_TYPE_TEXT:
		messageText, err := UnWrapMessageText(message)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			return
		}

		log.Info().Msgf("%s send to %s: %s %d", messageText.Origin, messageText.Target, messageText.Text, messageText.MessageLen)

		s.sendMessageTo(messageText.Origin, messageText)

	}
}

func (s ServerTCP) sendMessageTo(username string, message MessageInterface) {
	// conn := s.Users[username].
	user, ok := s.Users[username]
	if !ok || user.MessageConn == nil {
		log.Warn().Msgf("User %s not connected", username)
		return
	}

	conn := user.MessageConn
	_, err := conn.Write(message.Wrap())
	if err != nil {
		log.Err(err).Msgf("Error to send message to %s", username)
	}
}
