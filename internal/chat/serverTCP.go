package chat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
)

type ServerTCP struct {
	// Protocol    Protocol
	Users       map[string]connection
	userMux     sync.Mutex
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

	defer serve.Close()

	s.Users = make(map[string]connection)

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
		_, err = io.CopyN(buf, conn, int64(messageSize-5)) // 5 bytes for message type and message size
		if !checkErrorMessage(err) {
			break
		}

		raw := make([]byte, 5, messageSize)
		raw[0] = messageType
		binary.BigEndian.PutUint32(raw[1:], messageSize)
		raw = append(raw, buf.Bytes()...)

		log.Debug().Msgf("Message recived: %v", raw)

		s.handlerMessage(messageType, &raw, conn)
	}

	s.userMux.Lock()
	for _, user := range s.Users {
		if user.MessageConn == conn {
			delete(s.Users, user.UserName)
			log.Info().Msgf("User %s disconnected", user.UserName)
			break
		}
	}
	s.userMux.Unlock()
}

func (s *ServerTCP) handlerMessage(messageType uint8, message *[]byte, conn net.Conn) {
	switch messageType {
	case MESSAGE_TYPE_REQUEST_JOIN:
		messageRequestJoin, err := UnWrapMessageRequestJoin(message)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			return
		}

		log.Info().Msgf("User %s request join", messageRequestJoin.UserName)

		if _, ok := s.Users[messageRequestJoin.UserName]; ok {
			log.Warn().Msgf("User %s already connected", messageRequestJoin.UserName)

			conn.Write(MessageResponseJoin{
				UserName:  messageRequestJoin.UserName,
				Succeeded: false,
			}.Wrap())
			return
		}

		s.userMux.Lock()
		s.Users[messageRequestJoin.UserName] = connection{
			UserName:    messageRequestJoin.UserName,
			MessageConn: conn,
		}
		s.userMux.Unlock()

		s.sendMessageTo(messageRequestJoin.UserName, MessageResponseJoin{
			UserName:  messageRequestJoin.UserName,
			Succeeded: true,
		})

	case MESSAGE_TYPE_TEXT:
		messageText, err := UnWrapMessageText(message)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			return
		}

		log.Info().Msgf("%s send to %s: %s", messageText.Origin, messageText.Target, messageText.Text)

		if messageText.Target == MESSAGE_TARGET_ALL {
			s.sendMessageToAll(messageText)
			return
		}

		s.sendMessageTo(messageText.Target, messageText)

	default:
		log.Warn().Msgf("Message type %d not implemented", messageType)
	}
}

func (s *ServerTCP) sendMessageTo(username string, message MessageInterface) {
	s.userMux.Lock()
	user, ok := s.Users[username]
	if !ok || user.MessageConn == nil {
		log.Warn().Msgf("User %s not connected", username)
		log.Debug().Msgf("Users: %v", s.Users)
		return
	}
	s.userMux.Unlock()

	conn := user.MessageConn
	_, err := conn.Write(message.Wrap())
	if err != nil {
		log.Err(err).Msgf("Error to send message to %s", username)
	}
}

func (s *ServerTCP) sendMessageToAll(message MessageInterface) {
	s.userMux.Lock()
	for _, user := range s.Users {
		if user.MessageConn == nil {
			continue
		}
		user := user
		go func() {
			_, err := user.MessageConn.Write(message.Wrap())
			if err != nil {
				log.Err(err).Msgf("Error to send message to %s", user.UserName)
			}
		}()
	}
	s.userMux.Unlock()
}
