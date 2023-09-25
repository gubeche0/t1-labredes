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
	Users       map[string]net.Conn
	Address     string
	MessagePort int
	CommandPort int
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

		go s.HandleConnection(conn)
	}

}

func (s ServerTCP) HandleConnection(conn net.Conn) {
	defer conn.Close()
	log.Info().Msg("Connection accepted")

	for {
		buf := new(bytes.Buffer)

		var messageType uint8
		var messageSize uint32
		err := binary.Read(conn, binary.LittleEndian, &messageType)
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
		checkErrorMessage(err)

		// log.Debug().Msgf("Message recived: %v", (buf.String()))
		log.Debug().Msgf("Message recived: %v", (buf.Bytes()))

		// raw := append([]byte{messageType}, buf.Bytes()...)
		raw := []byte{messageType}
		sizeRaw := make([]byte, 4)
		binary.LittleEndian.PutUint32(sizeRaw, messageSize)
		raw = append(raw, sizeRaw...)
		raw = append(raw, buf.Bytes()...)

		messageText, err := UnWrapMessageText(&raw)
		if err != nil {
			log.Err(err).Msg("Error to unwrap message")
			continue
		}

		log.Info().Msgf("%s send to %s: %s", messageText.Origin, messageText.Target, messageText.Text)

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Err(err).Msg("Error to write message")
			break
		}
	}
}

func checkErrorMessage(err error) bool {
	if err == io.EOF {
		log.Info().Msg("Connection closed")
		return false
	} else if err != nil {
		log.Err(err).Msg("Error to read message")
		return false
	}

	return true
}
