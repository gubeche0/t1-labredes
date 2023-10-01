package chat

import "encoding/binary"

type MessageRequestJoin struct {
	MessageLen uint32
	UserName   string
}

func (m MessageRequestJoin) Wrap() []byte {

	length := uint32(len(m.UserName))
	length += 6 // 1 byte for message type, 4 bytes for message length, 1 byte for \n delimiter

	bytes := make([]byte, 0, length)

	bytes = append(bytes, byte(MESSAGE_TYPE_REQUEST_JOIN))
	bytes = binary.BigEndian.AppendUint32(bytes, length)

	bytes = append(bytes, []byte(m.UserName)...)
	bytes = append(bytes, '\n')

	return bytes
}

func UnWrapMessageRequestJoin(rawMessage *[]byte) (*MessageRequestJoin, error) {
	var msg MessageRequestJoin

	messageType := (*rawMessage)[0]
	if uint8(messageType) != MESSAGE_TYPE_REQUEST_JOIN {
		return nil, ErrInvalidMessageType
	}

	msg.MessageLen = binary.BigEndian.Uint32((*rawMessage)[1:5])
	msg.UserName = string((*rawMessage)[5 : len(*rawMessage)-1])
	return &msg, nil
}
