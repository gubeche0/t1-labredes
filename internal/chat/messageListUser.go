package chat

import "encoding/binary"

type MessageListUser struct {
	MessageLen uint64
	origin     string
}

func (m MessageListUser) Wrap() []byte {
	length := uint64(len(m.origin))
	length += 10 // 1 byte for message type, 8 bytes for message length, 1 byte for \n delimiter

	bytes := make([]byte, 0, length)

	bytes = append(bytes, byte(MESSAGE_TYPE_LIST_USERS))
	bytes = binary.BigEndian.AppendUint64(bytes, length)

	bytes = append(bytes, []byte(m.origin)...)
	bytes = append(bytes, '\n')

	return bytes
}

func (m MessageListUser) GetType() uint8 {
	return MESSAGE_TYPE_LIST_USERS
}

func UnWrapMessageListUser(rawMessage *[]byte) (*MessageListUser, error) {
	var msg MessageListUser

	messageType := (*rawMessage)[0]
	if uint8(messageType) != MESSAGE_TYPE_LIST_USERS {
		return nil, ErrInvalidMessageType
	}

	msg.MessageLen = binary.BigEndian.Uint64((*rawMessage)[1:9])
	msg.origin = string((*rawMessage)[9 : len(*rawMessage)-1])
	return &msg, nil
}
