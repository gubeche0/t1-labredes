package chat

import "encoding/binary"

type MessageResponseJoin struct {
	MessageLen uint64
	UserName   string
	Succeeded  bool
}

func (m MessageResponseJoin) Wrap() []byte {
	length := uint64(len(m.UserName))
	length += 11 // 1 byte for message type, 8 bytes for message length, 1 byte for \n delimiter, 1 byte for succeeded

	bytes := make([]byte, 0, length)

	bytes = append(bytes, byte(MESSAGE_TYPE_RESPONSE_JOIN))
	bytes = binary.BigEndian.AppendUint64(bytes, length)

	bytes = append(bytes, []byte(m.UserName)...)
	bytes = append(bytes, '\n')

	if m.Succeeded {
		bytes = append(bytes, '1')
	} else {
		bytes = append(bytes, '0')
	}

	return bytes
}

func (m MessageResponseJoin) GetType() uint8 {
	return MESSAGE_TYPE_RESPONSE_JOIN
}

func UnWrapMessageResponseJoin(rawMessage *[]byte) (*MessageResponseJoin, error) {
	var msg MessageResponseJoin

	messageType := (*rawMessage)[0]
	if uint8(messageType) != MESSAGE_TYPE_RESPONSE_JOIN {
		return nil, ErrInvalidMessageType
	}

	msg.MessageLen = binary.BigEndian.Uint64((*rawMessage)[1:9])
	msg.UserName = string((*rawMessage)[9 : len(*rawMessage)-2])
	msg.Succeeded = (*rawMessage)[len(*rawMessage)-1] == '1'
	return &msg, nil
}
