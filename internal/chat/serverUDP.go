package chat

import (
	"net"
	"sync"

	"github.com/rs/zerolog/log"
)

type ServerUDP struct {
	// Protocol    Protocol
	Users       map[string]connectionUDP
	userMux     sync.Mutex
	Address     string
	MessagePort int
	CommandPort int
	messageConn *net.UDPConn
	commandConn *net.UDPConn
}

type connectionUDP struct {
	UserName    string
	MessageAddr *net.UDPAddr
	CommandAddr *net.UDPAddr
}

func (s *ServerUDP) StartListenAndServer() error {
	// serve, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Address, s.MessagePort))
	// serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.Address, s.MessagePort))
	serve, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(s.Address),
		Port: s.MessagePort,
	})

	if err != nil {
		log.Fatal().Msgf("Error to create server: %s", err.Error())

		return err
	}

	s.messageConn = serve
	defer serve.Close()

	log.Info().Msgf("Server started at %s:%d", s.Address, s.MessagePort)
	s.Users = make(map[string]connectionUDP)

	// serverCommands, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Address, s.CommandPort))
	// serverCommandsAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.Address, s.CommandPort))
	serverCommands, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(s.Address),
		Port: s.CommandPort,
	})
	if err != nil {
		log.Fatal().Msgf("Error to create command server: %s", err.Error())

		return err
	}
	log.Info().Msgf("Command server started at %s:%d", s.Address, s.CommandPort)

	s.commandConn = serverCommands
	defer serverCommands.Close()

	go s.handleCommandConnection(serverCommands)
	for {
		raw, addr, err := reciveMessageUDP(s.messageConn)
		if !checkErrorMessage(err) {
			break
		}

		log.Debug().Msgf("Message recived from %s: %v", addr.String(), raw)

		s.handlerMessage((*raw)[0], raw, addr)
	}

	return nil
}

func (s *ServerUDP) handleCommandConnection(conn *net.UDPConn) {
	for {
		raw, addr, err := reciveMessageUDP(conn)
		if !checkErrorMessage(err) {
			break
		}

		log.Debug().Msgf("Message recived from %s: %v", addr.String(), raw)

		s.handlerCommandMessage((*raw)[0], raw, addr)
	}
}

func (s *ServerUDP) handlerCommandMessage(messageType uint8, message *[]byte, addr *net.UDPAddr) {
	switch messageType {
	case MESSAGE_TYPE_JOIN_REQUEST:
		messageJoinRequest, err := UnWrapMessageJoinRequest(message)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			return
		}

		// log.Info().Msgf("User %s request join", messageJoinRequest.UserName)
		log.Info().Msgf("User %s request join with user %s", addr.String(), messageJoinRequest.UserName)

		if _, ok := s.Users[messageJoinRequest.UserName]; ok {
			log.Warn().Msgf("User %s already connected", messageJoinRequest.UserName)

			s.commandConn.WriteToUDP(MessageJoinRequestResponse{
				UserName:  messageJoinRequest.UserName,
				Succeeded: false,
			}.Wrap(), addr)
			return
		}

		s.userMux.Lock()
		// s.Users[messageJoinRequest.UserName] = connection{
		// 	UserName:    messageJoinRequest.UserName,
		// 	MessageConn: conn,
		// }
		s.Users[messageJoinRequest.UserName] = connectionUDP{
			UserName:    messageJoinRequest.UserName,
			CommandAddr: addr,
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

func (s *ServerUDP) sendCommandMessageTo(username string, message MessageInterface) {
	s.userMux.Lock()
	user, ok := s.Users[username]
	if !ok || user.CommandAddr == nil {
		log.Warn().Msgf("User %s not connected in commands", username)
		log.Debug().Msgf("Users: %v", s.Users)
		s.userMux.Unlock()
		return
	}
	s.userMux.Unlock()

	// _, err := io.Copy(user.CommandConn, bytes.NewReader(message.Wrap()))
	_, err := s.commandConn.WriteToUDP(message.Wrap(), user.CommandAddr)
	if err != nil {
		log.Err(err).Msgf("Error to send message to %s", username)
	}
	log.Debug().Msgf("Command message sent to %s: %v", username, message.Wrap())
}

func (s *ServerUDP) handlerMessage(messageType uint8, message *[]byte, addr *net.UDPAddr) {
	switch messageType {
	case MESSAGE_TYPE_JOIN:
		messageJoin, err := UnWrapMessageJoin(message)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			return
		}

		// log.Info().Msgf("User %s request join", messageJoin.UserName)
		log.Info().Msgf("User %s join with user %s", addr.String(), messageJoin.UserName)

		if user, ok := s.Users[messageJoin.UserName]; ok && user.MessageAddr != nil {
			log.Warn().Msgf("User %s already connected", messageJoin.UserName)

			s.messageConn.Write(MessageJoinResponse{
				UserName:  messageJoin.UserName,
				Succeeded: false,
			}.Wrap())
			return
		}

		s.userMux.Lock()
		user, ok := s.Users[messageJoin.UserName]
		if !ok {
			s.Users[messageJoin.UserName] = connectionUDP{
				UserName:    messageJoin.UserName,
				MessageAddr: addr,
			}
		} else {
			user.MessageAddr = addr
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
		for username := range s.Users {
			users = append(users, username)
		}
		s.userMux.Unlock()

		log.Debug().Msgf("Users: %v", users)

		s.sendMessageTo(messageListUser.origin, MessageListUserResponse{
			Users: users,
		})

	case MESSAGE_TYPE_FILE:
		messageFile, err := UnWrapMessageFile(message)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			return
		}

		log.Info().Msgf("%s send to %s: %s with %d bytes", messageFile.Origin, messageFile.Target, messageFile.Filename, messageFile.Filesize)

		if messageFile.Target == MESSAGE_TARGET_ALL {
			s.sendMessageToAll(messageFile)
			return
		}

		s.sendMessageTo(messageFile.Target, messageFile)

	default:
		log.Warn().Msgf("Message type %d not implemented", messageType)
	}
}

func (s *ServerUDP) sendMessageTo(username string, message MessageInterface) {
	s.userMux.Lock()
	user, ok := s.Users[username]
	if !ok || user.MessageAddr == nil {
		log.Warn().Msgf("User %s not connected", username)
		log.Debug().Msgf("Users: %v", s.Users)
		s.userMux.Unlock()
		return
	}
	s.userMux.Unlock()

	// _, err := io.Copy(user.MessageConn, bytes.NewReader(message.Wrap()))
	_, err := s.messageConn.WriteToUDP(message.Wrap(), user.MessageAddr)
	if err != nil {
		log.Err(err).Msgf("Error to send message to %s", username)
	}

	log.Debug().Msgf("Message sent to %s: %v", username, message.Wrap())
}

func (s *ServerUDP) sendMessageToAll(message MessageInterface) {
	s.userMux.Lock()
	for name, user := range s.Users {
		if user.MessageAddr == nil {
			log.Warn().Msgf("User %s not connected", name)
			continue
		}
		user := user
		name := name
		go func() {
			// _, err := io.Copy(user.MessageConn, bytes.NewReader(message.Wrap()))
			_, err := s.messageConn.WriteToUDP(message.Wrap(), user.MessageAddr)
			if err != nil {
				log.Err(err).Msgf("Error to send message to %s", name)
			}
			log.Debug().Msgf("Message sent to %s", name)
		}()
	}
	s.userMux.Unlock()
}
