package chat

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
)

var (
	ErrInvalidMessageType = errors.New("invalid message type")
)

type MessageInterface interface {
	// UnWrap(bytes *[]byte) Message
	Wrap() []byte
	// GetOrigin() string
	// GetTarget() string
	GetType() uint8
}

type MessageText struct {
	MessageLen uint32
	Origin     string
	Target     string
	Text       string
}

func (m MessageText) Wrap() []byte {
	length := uint32(len(m.Text))
	length += uint32(len(m.Origin))
	length += uint32(len(m.Target))
	length += 8 // 3 bytes for \n delimiters, 1 byte for message type, 4 bytes for message length
	// log.Debug().Msgf("Message length: %d", length)

	bytes := make([]byte, 0, int(length))

	bytes = append(bytes, MESSAGE_TYPE_TEXT)
	bytes = binary.BigEndian.AppendUint32(bytes, length)

	bytes = append(bytes, []byte(m.Origin)...)
	bytes = append(bytes, '\n')
	bytes = append(bytes, []byte(m.Target)...)
	bytes = append(bytes, '\n')
	bytes = append(bytes, []byte(m.Text)...)
	bytes = append(bytes, '\n')

	return bytes
}

func (m MessageText) GetOrigin() string {
	return m.Origin
}

func (m MessageText) GetTarget() string {
	return m.Target
}

func (m MessageText) GetType() uint8 {
	return MESSAGE_TYPE_TEXT
}

func UnWrapMessageText(rawMessage *[]byte) (*MessageText, error) {
	var msg MessageText

	messageType := (*rawMessage)[0]
	if uint8(messageType) != MESSAGE_TYPE_TEXT {
		return nil, ErrInvalidMessageType
	}

	msg.MessageLen = binary.BigEndian.Uint32((*rawMessage)[1:5])

	reader := bufio.NewReader(bytes.NewReader((*rawMessage)[5:]))

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

	return &msg, nil
}

type MessageFile struct {
	Origin   string
	Target   string
	Filename string
	Filesize uint32
	Filedata []byte
}
