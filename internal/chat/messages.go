package chat

import (
	"bufio"
	"bytes"
	"errors"

	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidMessageType = errors.New("invalid message type")
)

type MessageInterface interface {
	// UnWrap(bytes *[]byte) Message
	Wrap() []byte
	GetOrigin() string
	GetTarget() string
	GetType() uint8
}

type MessageText struct {
	MessageLen uint32
	Origin     string
	Target     string
	Text       string
}

func (m MessageText) Wrap() []byte {
	bytes := make([]byte, 0, 1+len(m.Text))

	length := uint32(len(m.Text))
	length += uint32(len(m.Origin))
	length += uint32(len(m.Target))
	length += 3

	log.Debug().Msgf("Message length: %d", length)

	bytes = append(bytes, MESSAGE_TYPE_TEXT)

	bytes = append(bytes, byte(length>>24))
	bytes = append(bytes, byte(length>>16))
	bytes = append(bytes, byte(length>>8))
	bytes = append(bytes, byte(length))

	bytes = append(bytes, []byte(m.Origin)...)
	bytes = append(bytes, '\n')
	bytes = append(bytes, []byte(m.Target)...)
	bytes = append(bytes, '\n')
	bytes = append(bytes, []byte(m.Text)...)
	bytes = append(bytes, '\n')

	return bytes
}

func UnWrapMessageText(rawMessage *[]byte) (*MessageText, error) {
	var msg MessageText

	messageType := (*rawMessage)[0]
	if messageType != MESSAGE_TYPE_TEXT {
		return nil, ErrInvalidMessageType
	}

	reader := bufio.NewReader(bytes.NewReader((*rawMessage)[1:]))

	// msg.MessageLen = uint32((*rawMessage)[1])<<24 | uint32((*rawMessage)[2])<<16 | uint32((*rawMessage)[3])<<8 | uint32((*rawMessage)[4])
	// msg.Origin = bytes.ReadBytes('\n')

	origin, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	msg.Origin = string(origin[:len(origin)-1])

	target, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	msg.Target = string(target[:len(target)-1])

	text, err := reader.ReadBytes('\n')
	if err != nil {

		return nil, err
	}
	msg.Text = string(text[:len(text)-1])

	// bytes.NewReader((*rawMessage)[5:]).ReadBytes('\n')

	// msg.Origin, err :=

	return &msg, nil
}

type MessageFile struct {
	Origin   string
	Target   string
	Filename string
	Filesize uint32
	Filedata []byte
}
