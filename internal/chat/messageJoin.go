package chat

import "encoding/binary"

type MessageJoin struct {
	MessageLen uint64
	UserName   string
}

func (m MessageJoin) Wrap() []byte {

	length := uint64(len(m.UserName))
	length += 10 // 1 byte for message type, 8 bytes for message length, 1 byte for \n delimiter

	bytes := make([]byte, 0, length)

	bytes = append(bytes, byte(MESSAGE_TYPE_JOIN))
	bytes = binary.BigEndian.AppendUint64(bytes, length)

	bytes = append(bytes, []byte(m.UserName)...)
	bytes = append(bytes, '\n')

	return bytes
}

func (m MessageJoin) GetType() uint8 {
	return MESSAGE_TYPE_JOIN
}

func UnWrapMessageJoin(rawMessage *[]byte) (*MessageJoinRequest, error) {
	var msg MessageJoinRequest

	messageType := (*rawMessage)[0]
	if uint8(messageType) != MESSAGE_TYPE_JOIN {
		return nil, ErrInvalidMessageType
	}

	msg.MessageLen = binary.BigEndian.Uint64((*rawMessage)[1:9])
	msg.UserName = string((*rawMessage)[9 : len(*rawMessage)-1])
	return &msg, nil
}
