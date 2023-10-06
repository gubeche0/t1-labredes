package chat

import "encoding/binary"

type MessageListUserResponse struct {
	MessageLen uint64
	Users      []string
}

func (m MessageListUserResponse) Wrap() []byte {
	length := uint64(0)
	for _, user := range m.Users {
		length += uint64(len(user))
		length += 1 // 1 byte for \n delimiter
	}
	length += 1 // 1 byte for message type
	length += 8 // 8 bytes for message length

	bytes := make([]byte, 0, length)

	bytes = append(bytes, byte(MESSAGE_TYPE_LIST_USERS_RESPONSE))
	bytes = binary.BigEndian.AppendUint64(bytes, length)

	for _, user := range m.Users {
		bytes = append(bytes, []byte(user)...)
		bytes = append(bytes, '\n')
	}

	return bytes
}

func (m MessageListUserResponse) GetType() uint8 {
	return MESSAGE_TYPE_LIST_USERS_RESPONSE
}

func UnWrapMessageListUserResponse(rawMessage *[]byte) (*MessageListUserResponse, error) {
	var msg MessageListUserResponse

	messageType := (*rawMessage)[0]
	if uint8(messageType) != MESSAGE_TYPE_LIST_USERS_RESPONSE {
		return nil, ErrInvalidMessageType
	}

	msg.MessageLen = binary.BigEndian.Uint64((*rawMessage)[1:9])
	msg.Users = make([]string, 0)

	currentUser := make([]byte, 0)
	for i := 9; i < len(*rawMessage); i++ {
		if (*rawMessage)[i] == '\n' {
			msg.Users = append(msg.Users, string(currentUser))
			currentUser = make([]byte, 0)
		} else {
			currentUser = append(currentUser, (*rawMessage)[i])
		}
	}

	return &msg, nil
}
