package chat

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type MessageFile struct {
	MessageLen uint64
	Origin     string
	Target     string
	Filename   string
	Filesize   uint64
	Filedata   []byte
}

func (m MessageFile) Wrap() []byte {

	length := uint64(len(m.Origin))
	length += uint64(len(m.Target))
	length += uint64(len(m.Filename))
	length += uint64(len(m.Filedata))
	length += 12 // 3 bytes for \n delimiters, 1 byte for message type, 8 bytes for message length
	length += 8  // 8 bytes for file size

	bytes := make([]byte, 0, int(length))

	bytes = append(bytes, MESSAGE_TYPE_FILE)
	bytes = binary.BigEndian.AppendUint64(bytes, length)

	bytes = append(bytes, []byte(m.Origin)...)
	bytes = append(bytes, '\n')
	bytes = append(bytes, []byte(m.Target)...)
	bytes = append(bytes, '\n')
	bytes = append(bytes, []byte(m.Filename)...)
	bytes = append(bytes, '\n')
	if m.Filesize == 0 {
		m.Filesize = uint64(len(m.Filedata))
	}
	bytes = binary.BigEndian.AppendUint64(bytes, m.Filesize)
	bytes = append(bytes, m.Filedata...)

	return bytes
}

func (m MessageFile) GetType() uint8 {
	return MESSAGE_TYPE_FILE
}

func UnWrapMessageFile(rawMessage *[]byte) (*MessageFile, error) {
	if len(*rawMessage) < 9 {
		return nil, ErrInvalidMessageType
	}

	if (*rawMessage)[0] != MESSAGE_TYPE_FILE {
		return nil, ErrInvalidMessageType
	}

	var msg MessageFile

	msg.MessageLen = binary.BigEndian.Uint64((*rawMessage)[1:9])

	reader := bufio.NewReader(bytes.NewReader((*rawMessage)[9:]))
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

	filename, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	msg.Filename = string(filename[:len(filename)-1])

	index := 1 + 8 + len(origin) + len(target) + len(filename) // 1 byte for message type, 8 bytes for message length

	// log.Debug().Msgf("Message length: %d", index)

	msg.Filesize = binary.BigEndian.Uint64((*rawMessage)[index:])
	index += 8

	// log.Debug().Msgf("Message length: %d", msg.Filesize)

	msg.Filedata = (*rawMessage)[index:]

	// log.Debug().Msgf("Message data: %v", msg.Filedata)

	return &msg, nil
}
