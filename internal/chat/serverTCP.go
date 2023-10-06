package chat

import (
	"fmt"
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

	log.Info().Msgf("Server started at %s:%d", s.Address, s.MessagePort)
	s.Users = make(map[string]connection)

	serverCommands, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Address, s.CommandPort))
	if err != nil {
		log.Fatal().Msgf("Error to create command server: %s", err.Error())

		return err
	}
	log.Info().Msgf("Command server started at %s:%d", s.Address, s.CommandPort)

	defer serverCommands.Close()

	go func() {
		for {
			conn, err := serverCommands.Accept()
			if err != nil {
				log.Fatal().Msgf("Error to accept connection: %s", err.Error())
			}

			go s.handleCommandConnection(conn)
		}
	}()

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
		raw, err := reciveMessage(conn)
		if !checkErrorMessage(err) {
			break
		}

		log.Debug().Msgf("Message recived: %v", raw)

		s.handlerMessage((*raw)[0], raw, conn)
	}

	s.userMux.Lock()
	for _, user := range s.Users {
		if user.MessageConn == conn {
			if user.CommandConn != nil {
				user.CommandConn.Close()
			}
			delete(s.Users, user.UserName)
			log.Info().Msgf("User %s disconnected", user.UserName)
			break
		}
	}
	s.userMux.Unlock()
}

func (s *ServerTCP) handlerMessage(messageType uint8, message *[]byte, conn net.Conn) {
	switch messageType {
	case MESSAGE_TYPE_JOIN:
		messageJoin, err := UnWrapMessageJoin(message)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			return
		}

		// log.Info().Msgf("User %s request join", messageJoin.UserName)
		log.Info().Msgf("User %s join with user %s", conn.RemoteAddr().String(), messageJoin.UserName)

		if user, ok := s.Users[messageJoin.UserName]; ok && user.MessageConn != nil {
			log.Warn().Msgf("User %s already connected", messageJoin.UserName)

			conn.Write(MessageJoinResponse{
				UserName:  messageJoin.UserName,
				Succeeded: false,
			}.Wrap())
			return
		}

		s.userMux.Lock()
		user, ok := s.Users[messageJoin.UserName]
		if !ok {
			s.Users[messageJoin.UserName] = connection{
				UserName:    messageJoin.UserName,
				MessageConn: conn,
			}
		} else {
			user.MessageConn = conn
			s.Users[messageJoin.UserName] = user
		}
		s.userMux.Unlock()

		s.sendMessageTo(messageJoin.UserName, MessageJoinRequestResponse{
			UserName:  messageJoin.UserName,
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

	case MESSAGE_TYPE_LIST_USERS:
		messageListUser, err := UnWrapMessageListUser(message)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			return
		}

		log.Info().Msgf("User %s request list users", messageListUser.origin)

		users := make([]string, 0)
		s.userMux.Lock()
		for _, user := range s.Users {
			users = append(users, user.UserName)
		}
		s.userMux.Unlock()

		log.Debug().Msgf("Users: %v", users)

		s.sendMessageTo(messageListUser.origin, MessageListUserResponse{
			Users: users,
		})

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
		s.userMux.Unlock()
		return
	}
	s.userMux.Unlock()

	conn := user.MessageConn
	_, err := conn.Write(message.Wrap())
	if err != nil {
		log.Err(err).Msgf("Error to send message to %s", username)
	}

	log.Debug().Msgf("Message sent to %s: %v", username, message.Wrap())
}

func (s *ServerTCP) sendMessageToAll(message MessageInterface) {
	s.userMux.Lock()
	for _, user := range s.Users {
		if user.MessageConn == nil {
			log.Warn().Msgf("User %s not connected", user.UserName)
			continue
		}
		user := user
		go func() {
			_, err := user.MessageConn.Write(message.Wrap())
			if err != nil {
				log.Err(err).Msgf("Error to send message to %s", user.UserName)
			}
			log.Debug().Msgf("Message sent to %s", user.UserName)
		}()
	}
	s.userMux.Unlock()
}

func (s *ServerTCP) handleCommandConnection(conn net.Conn) {
	defer conn.Close()
	log.Info().Msg("Command connection accepted")

	for {
		raw, err := reciveMessage(conn)
		if !checkErrorMessage(err) {
			break
		}

		log.Debug().Msgf("Command recived: %v", raw)

		s.handlerCommandMessage((*raw)[0], raw, conn)
	}

	s.userMux.Lock()
	for _, user := range s.Users {
		if user.CommandConn == conn {
			user.CommandConn = nil
			log.Info().Msgf("User %s disconnected from commands", user.UserName)
			break
		}
	}
	s.userMux.Unlock()
}

func (s *ServerTCP) handlerCommandMessage(messageType uint8, message *[]byte, conn net.Conn) {
	switch messageType {
	case MESSAGE_TYPE_JOIN_REQUEST:
		messageJoinRequest, err := UnWrapMessageJoinRequest(message)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			return
		}

		// log.Info().Msgf("User %s request join", messageJoinRequest.UserName)
		log.Info().Msgf("User %s request join with user %s", conn.RemoteAddr().String(), messageJoinRequest.UserName)

		if _, ok := s.Users[messageJoinRequest.UserName]; ok {
			log.Warn().Msgf("User %s already connected", messageJoinRequest.UserName)

			conn.Write(MessageJoinRequestResponse{
				UserName:  messageJoinRequest.UserName,
				Succeeded: false,
			}.Wrap())
			return
		}

		s.userMux.Lock()
		// s.Users[messageJoinRequest.UserName] = connection{
		// 	UserName:    messageJoinRequest.UserName,
		// 	MessageConn: conn,
		// }
		s.Users[messageJoinRequest.UserName] = connection{
			UserName:    messageJoinRequest.UserName,
			CommandConn: conn,
		}
		s.userMux.Unlock()

		s.sendCommandMessageTo(messageJoinRequest.UserName, MessageJoinRequestResponse{
			UserName:  messageJoinRequest.UserName,
			Succeeded: true,
		})
	default:
		log.Warn().Msgf("Command message type %d not implemented", messageType)
	}
}

func (s *ServerTCP) sendCommandMessageTo(username string, message MessageInterface) {
	s.userMux.Lock()
	user, ok := s.Users[username]
	if !ok || user.CommandConn == nil {
		log.Warn().Msgf("User %s not connected in commands", username)
		log.Debug().Msgf("Users: %v", s.Users)
		s.userMux.Unlock()
		return
	}
	s.userMux.Unlock()

	conn := user.CommandConn
	_, err := conn.Write(message.Wrap())
	if err != nil {
		log.Err(err).Msgf("Error to send message to %s", username)
	}
}
