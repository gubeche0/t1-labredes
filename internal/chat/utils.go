package chat

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"

	"github.com/rs/zerolog/log"
)

func checkErrorMessage(err error) bool {
	if err == io.EOF {
		log.Debug().Msg("Connection closed")
		return false
	} else if err != nil {
		log.Err(err).Msg("Error to read message")
		return false
	}

	return true
}

func checkErrorMessageClient(err error) bool {
	if err == io.EOF {
		log.Fatal().Msg("Connection closed")
		return false
	} else if err != nil {
		log.Fatal().Err(err).Msg("Error to read message")
		return false
	}

	return true
}

func reciveMessage(conn net.Conn) (*[]byte, error) {
	buf := new(bytes.Buffer)

	var messageType uint8
	var messageSize uint64
	err := binary.Read(conn, binary.BigEndian, &messageType)
	if err != nil {
		return nil, err
	}
	err = binary.Read(conn, binary.BigEndian, &messageSize)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Message type recived: %v", messageType)
	log.Debug().Msgf("Message size recived: %v", messageSize)

	// message, err := bufio.NewReader(conn).ReadBytes('\n')
	_, err = io.CopyN(buf, conn, int64(messageSize-9)) // 9 bytes for message type (1) and message size (8)
	if err != nil {
		return nil, err
	}

	raw := make([]byte, 9, messageSize)
	raw[0] = messageType
	binary.BigEndian.PutUint64(raw[1:], messageSize)
	raw = append(raw, buf.Bytes()...)

	return &raw, nil
}
